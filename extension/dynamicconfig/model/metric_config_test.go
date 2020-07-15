// Copyright The OpenTelemetry Authors
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

package model

import (
	"bytes"
	"testing"
)

func TestMetricConfigProto(t *testing.T) {
	config := MetricConfig{
		Schedules: []*Schedule{{}, {}},
	}

	configProto := config.Proto()

	if len(configProto.Schedules) != 2 {
		t.Errorf("improper conversion to proto")
	}
}

func TestMetricConfigHash(t *testing.T) {
	configA := MetricConfig{
		Schedules: []*Schedule{
			{Period: "MIN_1"},
			{Period: "MIN_5"},
		},
	}

	configB := MetricConfig{
		Schedules: []*Schedule{
			{Period: "MIN_5"},
			{Period: "MIN_1"},
		},
	}

	configC := MetricConfig{
		Schedules: []*Schedule{
			{Period: "MIN_1"},
		},
	}

	if !bytes.Equal(configA.Hash(), configB.Hash()) {
		t.Errorf("identical configs with different hashes")
	}

	if bytes.Equal(configA.Hash(), configC.Hash()) {
		t.Errorf("different configs with identical hashes")
	}
}

func TestMetricConfigHashEmpty(t *testing.T) {
	config := &MetricConfig{}
	hash := config.Hash()
	if !bytes.Equal(hash, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) {
		t.Errorf("expected all zeros, got: %v", hash)
	}
}
