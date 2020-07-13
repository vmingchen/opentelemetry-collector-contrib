// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package splunkhecexporter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumerdata"
	"go.opentelemetry.io/collector/testutil/metricstestutil"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	got, err := createExporter(nil, zap.NewNop())
	assert.EqualError(t, err, "nil config")
	assert.Nil(t, got)

	config := &Config{
		Token:    "someToken",
		Endpoint: "https://example.com:8088",
		Timeout:  1 * time.Second,
	}
	got, err = createExporter(config, zap.NewNop())
	assert.NoError(t, err)
	require.NotNil(t, got)
}

func TestConsumeMetricsData(t *testing.T) {
	smallBatch := &consumerdata.MetricsData{
		Node: &commonpb.Node{
			ServiceInfo: &commonpb.ServiceInfo{Name: "test_splunk"},
		},
		Resource: &resourcepb.Resource{Type: "test"},
		Metrics: []*metricspb.Metric{
			metricstestutil.Gauge(
				"test_gauge",
				[]string{"k0", "k1"},
				metricstestutil.Timeseries(
					time.Now(),
					[]string{"v0", "v1"},
					metricstestutil.Double(time.Now(), 123))),
		},
	}
	tests := []struct {
		name                 string
		md                   *consumerdata.MetricsData
		reqTestFunc          func(t *testing.T, r *http.Request)
		httpResponseCode     int
		numDroppedTimeSeries int
		wantErr              bool
	}{
		{
			name: "happy_path",
			md:   smallBatch,
			reqTestFunc: func(t *testing.T, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, "keep-alive", r.Header.Get("Connection"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "OpenTelemetry-Collector Splunk Exporter/v0.0.1", r.Header.Get("User-Agent"))
				assert.Equal(t, "Splunk 1234", r.Header.Get("Authorization"))
				if r.Header.Get("Content-Encoding") == "gzip" {
					t.Fatal("Small batch should not be compressed")
				}
				firstPayload := strings.Split(string(body), "\n\r\n\r")[0]
				var metric splunkMetric
				err = json.Unmarshal([]byte(firstPayload), &metric)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, "test", metric.Source)
				assert.Equal(t, "test_type", metric.SourceType)
				assert.Equal(t, "test_index", metric.Index)

			},
			httpResponseCode: http.StatusAccepted,
		},
		{
			name:                 "response_forbidden",
			md:                   smallBatch,
			reqTestFunc:          nil,
			httpResponseCode:     http.StatusForbidden,
			numDroppedTimeSeries: 1,
			wantErr:              true,
		},
		{
			name:             "large_batch",
			md:               generateLargeBatch(t),
			reqTestFunc:      nil,
			httpResponseCode: http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.reqTestFunc != nil {
					tt.reqTestFunc(t, r)
				}
				w.WriteHeader(tt.httpResponseCode)
			}))
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.NoError(t, err)

			options := &exporterOptions{
				url:   serverURL,
				token: "1234",
			}
			config := &Config{
				Source:     "test",
				SourceType: "test_type",
				Token:      "1234",
				Index:      "test_index",
			}
			sender := buildClient(options, config, zap.NewNop())

			numDroppedTimeSeries, err := sender.pushMetricsData(context.Background(), *tt.md)
			assert.Equal(t, tt.numDroppedTimeSeries, numDroppedTimeSeries)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func generateLargeBatch(t *testing.T) *consumerdata.MetricsData {
	md := &consumerdata.MetricsData{
		Node: &commonpb.Node{
			ServiceInfo: &commonpb.ServiceInfo{Name: "test_splunkhec"},
		},
		Resource: &resourcepb.Resource{Type: "test"},
	}

	ts := time.Now()
	for i := 0; i < 65000; i++ {
		md.Metrics = append(md.Metrics,
			metricstestutil.Gauge(
				"test_"+strconv.Itoa(i),
				[]string{"k0", "k1"},
				metricstestutil.Timeseries(
					time.Now(),
					[]string{"v0", "v1"},
					&metricspb.Point{
						Timestamp: metricstestutil.Timestamp(ts),
						Value:     &metricspb.Point_Int64Value{Int64Value: int64(i)},
					},
				),
			),
		)
	}

	return md
}
