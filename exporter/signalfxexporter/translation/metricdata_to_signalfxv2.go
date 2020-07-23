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

package translation

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	sfxpb "github.com/signalfx/com_signalfx_metrics_protobuf/model"
	"go.opentelemetry.io/collector/consumer/consumerdata"
	"go.uber.org/zap"
)

// Some fields on SignalFx protobuf are pointers, in order to reduce
// allocations create the most used ones.
var (
	// SignalFx metric types used in the conversions.
	sfxMetricTypeGauge             = sfxpb.MetricType_GAUGE
	sfxMetricTypeCumulativeCounter = sfxpb.MetricType_CUMULATIVE_COUNTER

	// Array used to map OpenCensus metric descriptor to SignaFx metric type.
	ocDescriptorToMetricType = [8]*sfxpb.MetricType{
		nil, // index metricspb.MetricDescriptor_UNSPECIFIED = 0

		&sfxMetricTypeGauge, // index metricspb.MetricDescriptor_GAUGE_INT64 = 1
		&sfxMetricTypeGauge, // index metricspb.MetricDescriptor_GAUGE_DOUBLE = 2

		// For distribution total count, sum, and individual buckets are all
		// cumulative counters.
		&sfxMetricTypeCumulativeCounter, // index metricspb.MetricDescriptor_GAUGE_DISTRIBUTION = 3

		&sfxMetricTypeCumulativeCounter, // index metricspb.MetricDescriptor_CUMULATIVE_INT64 = 4
		&sfxMetricTypeCumulativeCounter, // index metricspb.MetricDescriptor_CUMULATIVE_DOUBLE = 5

		// For distribution total count, sum, and individual buckets are all
		// cumulative counters.
		&sfxMetricTypeCumulativeCounter, // index metricspb.MetricDescriptor_CUMULATIVE_DISTRIBUTION = 6

		// For summary total count and sum are cumulative counters, however,
		// quantiles are gauges.
		&sfxMetricTypeGauge, // index metricspb.MetricDescriptor_SUMMARY = 7
	}

	// Pre-defined SignalFx Metric types according to their usage. These are
	// handy because the SignalFx protobuf needs pointers not constants.
	// The mapping of these comes from SignalFx docs regarding Prometheus metrics.
	// https://docs.signalfx.com/en/latest/integrations/agent/monitors/prometheus-exporter.html#overview
	bucketMetricType     = &sfxMetricTypeCumulativeCounter
	quantileMetricType   = &sfxMetricTypeGauge
	sumMetricType        = &sfxMetricTypeCumulativeCounter
	totalCountMetricType = &sfxMetricTypeCumulativeCounter

	// Some standard dimension keys.
	// upper bound dimension key for histogram buckets.
	upperBoundDimensionKey = "upper_bound"
	// quantile dimension key for summary quantiles.
	quantileDimensionKey = "quantile"

	// infinity bound dimension value is used on all histograms.
	infinityBoundSFxDimValue = float64ToDimValue(math.Inf(1))
)

func MetricDataToSignalFxV2(
	logger *zap.Logger,
	metricTranslator *MetricTranslator,
	md consumerdata.MetricsData,
) (sfxDataPoints []*sfxpb.DataPoint, numDroppedTimeSeries int) {
	sfxDataPoints, numDroppedTimeSeries = metricDataToSfxDataPoints(logger, md)
	if metricTranslator != nil {
		sfxDataPoints = metricTranslator.TranslateDataPoints(logger, sfxDataPoints)
	}
	sanitizeDataPointDimensions(sfxDataPoints)
	return
}

func metricDataToSfxDataPoints(
	logger *zap.Logger,
	md consumerdata.MetricsData,
) (sfxDataPoints []*sfxpb.DataPoint, numDroppedTimeSeries int) {

	// The final number of data points is likely larger than len(md.Metrics)
	// but at least that is expected.
	sfxDataPoints = make([]*sfxpb.DataPoint, 0, len(md.Metrics))

	var err error

	// Labels from Node and Resource.
	// TODO: Options to add lib, service name, etc as dimensions?
	//  Q.: what about resource type?
	nodeAttribs := md.Node.GetAttributes()
	resourceAttribs := md.Resource.GetLabels()
	numExtraDimensions := len(nodeAttribs) + len(resourceAttribs)
	var extraDimensions []*sfxpb.Dimension
	if numExtraDimensions > 0 {
		extraDimensions = make([]*sfxpb.Dimension, 0, numExtraDimensions)
		extraDimensions = appendAttributesToDimensions(extraDimensions, nodeAttribs)
		extraDimensions = appendAttributesToDimensions(extraDimensions, resourceAttribs)
	}

	for _, metric := range md.Metrics {
		if metric == nil || metric.MetricDescriptor == nil {
			logger.Warn("Received nil metrics data or nil descriptor for metrics")
			numDroppedTimeSeries += len(metric.GetTimeseries())
			continue
		}

		// Build the fixed parts for this metrics from the descriptor.
		descriptor := metric.MetricDescriptor
		metricName := descriptor.Name

		metricType := fromOCMetricDescriptorToMetricType(descriptor.Type)
		numLabels := len(descriptor.LabelKeys)

		for _, series := range metric.Timeseries {
			dimensions := make([]*sfxpb.Dimension, numLabels+numExtraDimensions)
			copy(dimensions, extraDimensions)
			for i := 0; i < numLabels; i++ {
				dimension := &sfxpb.Dimension{
					Key:   descriptor.LabelKeys[i].Key,
					Value: series.LabelValues[i].Value,
				}
				dimensions[numExtraDimensions+i] = dimension
			}

			for _, dp := range series.Points {

				var msec int64
				if dp.Timestamp != nil {
					msec = dp.Timestamp.Seconds*1e3 + int64(dp.Timestamp.Nanos)/1e6
				}

				sfxDataPoint := &sfxpb.DataPoint{
					// Source field is not set by code seen at
					// github.com/signalfx/golib/sfxclient/httpsink.go
					Metric:     metricName,
					MetricType: metricType,
					Timestamp:  msec,
					Dimensions: dimensions,
				}

				switch pv := dp.Value.(type) {
				case *metricspb.Point_Int64Value:
					sfxDataPoint.Value = sfxpb.Datum{IntValue: &pv.Int64Value}
					sfxDataPoints = append(sfxDataPoints, sfxDataPoint)

				case *metricspb.Point_DoubleValue:
					sfxDataPoint.Value = sfxpb.Datum{DoubleValue: &pv.DoubleValue}
					sfxDataPoints = append(sfxDataPoints, sfxDataPoint)

				case *metricspb.Point_DistributionValue:
					sfxDataPoints, err = appendDistributionValues(
						sfxDataPoints,
						sfxDataPoint,
						pv.DistributionValue)
					if err != nil {
						numDroppedTimeSeries++
						logger.Warn(
							"Timeseries for distribution metric dropped",
							zap.Error(err),
							zap.String("metric", sfxDataPoint.Metric))
					}
				case *metricspb.Point_SummaryValue:
					sfxDataPoints, err = appendSummaryValues(
						sfxDataPoints,
						sfxDataPoint,
						pv.SummaryValue)
					if err != nil {
						numDroppedTimeSeries++
						logger.Warn(
							"Timeseries for summary metric dropped",
							zap.Error(err),
							zap.String("metric", sfxDataPoint.Metric))
					}
				default:
					numDroppedTimeSeries++
					logger.Warn(
						"Timeseries dropped to unexpected metric type",
						zap.String("metric", sfxDataPoint.Metric))
				}

			}
		}
	}

	return sfxDataPoints, numDroppedTimeSeries
}

func appendAttributesToDimensions(
	dimensions []*sfxpb.Dimension,
	attribs map[string]string,
) []*sfxpb.Dimension {

	for k, v := range attribs {
		dim := &sfxpb.Dimension{
			Key:   k,
			Value: v,
		}
		dimensions = append(dimensions, dim)
	}
	return dimensions
}

func fromOCMetricDescriptorToMetricType(ocType metricspb.MetricDescriptor_Type) *sfxpb.MetricType {
	if ocType > metricspb.MetricDescriptor_UNSPECIFIED &&
		ocType <= metricspb.MetricDescriptor_SUMMARY {
		return ocDescriptorToMetricType[int(ocType)]
	}

	return nil
}

func appendDistributionValues(
	sfxDataPoints []*sfxpb.DataPoint,
	sfxBaseDataPoint *sfxpb.DataPoint,
	distributionValue *metricspb.DistributionValue,
) ([]*sfxpb.DataPoint, error) {

	// Translating distribution values per symmetrical recommendations to Prometheus:
	// https://docs.signalfx.com/en/latest/integrations/agent/monitors/prometheus-exporter.html#overview

	// 1. The total count gets converted to a cumulative counter called
	// <basename>_count.
	// 2. The total sum gets converted to a cumulative counter called <basename>.
	sfxDataPoints = appendTotalAndSum(
		sfxDataPoints,
		sfxBaseDataPoint,
		&distributionValue.Count,
		&distributionValue.Sum)

	// 3. Each histogram bucket is converted to a cumulative counter called
	// <basename>_bucket and will include a dimension called upper_bound that
	// specifies the maximum value in that bucket. This metric specifies the
	// number of events with a value that is less than or equal to the upper
	// bound.
	metricName := sfxBaseDataPoint.Metric + "_bucket"
	explicitBuckets := distributionValue.BucketOptions.GetExplicit()
	if explicitBuckets == nil {
		return sfxDataPoints, fmt.Errorf(
			"unknown bucket options type for metric %q",
			sfxBaseDataPoint.Metric)
	}
	bounds := explicitBuckets.Bounds
	sfxBounds := make([]string, len(bounds)+1)
	for i := 0; i < len(bounds); i++ {
		sfxBounds[i] = float64ToDimValue(bounds[i])
	}
	sfxBounds[len(sfxBounds)-1] = infinityBoundSFxDimValue

	for i, bucket := range distributionValue.Buckets {

		// Adding the "upper_bound" dimension.
		bucketDimensions := make([]*sfxpb.Dimension, len(sfxBaseDataPoint.Dimensions)+1)
		copy(bucketDimensions, sfxBaseDataPoint.Dimensions)

		bucketDimensions[len(bucketDimensions)-1] = &sfxpb.Dimension{
			Key:   upperBoundDimensionKey,
			Value: sfxBounds[i],
		}

		bucketDP := *sfxBaseDataPoint
		bucketDP.Dimensions = bucketDimensions
		bucketDP.Metric = metricName
		bucketDP.MetricType = bucketMetricType
		count := bucket.Count
		bucketDP.Value = sfxpb.Datum{IntValue: &count}

		sfxDataPoints = append(sfxDataPoints, &bucketDP)
	}

	return sfxDataPoints, nil
}

func appendSummaryValues(
	sfxDataPoints []*sfxpb.DataPoint,
	sfxBaseDataPoint *sfxpb.DataPoint,
	summaryValue *metricspb.SummaryValue,
) ([]*sfxpb.DataPoint, error) {

	// Translating summary values per symmetrical recommendations to Prometheus:
	// https://docs.signalfx.com/en/latest/integrations/agent/monitors/prometheus-exporter.html#overview

	// 1. The total count gets converted to a cumulative counter called
	// <basename>_count.
	// 2. The total sum gets converted to a cumulative counter called <basename>
	count := summaryValue.GetCount().GetValue()
	sum := summaryValue.GetSum().GetValue()
	sfxDataPoints = appendTotalAndSum(
		sfxDataPoints,
		sfxBaseDataPoint,
		&count,
		&sum)

	// 3. Each quantile value is converted to a gauge called <basename>_quantile
	// and will include a dimension called quantile that specifies the quantile.
	percentiles := summaryValue.GetSnapshot().GetPercentileValues()
	if percentiles == nil {
		return sfxDataPoints, fmt.Errorf(
			"unknown percentiles values for summary metric %q",
			sfxBaseDataPoint.Metric)
	}
	metricName := sfxBaseDataPoint.Metric + "_quantile"
	for _, quantile := range percentiles {

		// Adding the "quantile" dimension.
		quantileDimensions := make([]*sfxpb.Dimension, len(sfxBaseDataPoint.Dimensions)+1)
		copy(quantileDimensions, sfxBaseDataPoint.Dimensions)

		// If a dimension "quantile" was already specified: the last one wins.
		quantileDimensions[len(quantileDimensions)-1] = &sfxpb.Dimension{
			Key:   quantileDimensionKey,
			Value: float64ToDimValue(quantile.Percentile),
		}

		quantileDP := *sfxBaseDataPoint
		quantileDP.Dimensions = quantileDimensions
		quantileDP.Metric = metricName
		quantileDP.MetricType = quantileMetricType
		value := quantile.Value
		quantileDP.Value = sfxpb.Datum{DoubleValue: &value}

		sfxDataPoints = append(sfxDataPoints, &quantileDP)
	}

	return sfxDataPoints, nil
}

func appendTotalAndSum(
	sfxDataPoints []*sfxpb.DataPoint,
	sfxBaseDataPoint *sfxpb.DataPoint,
	count *int64,
	sum *float64,
) []*sfxpb.DataPoint {

	sfxDataPoints = append(
		sfxDataPoints,
		buildTotalDataPoint(sfxBaseDataPoint, count),
		buildSumDataPoint(sfxBaseDataPoint, sum))

	return sfxDataPoints
}

func buildTotalDataPoint(
	sfxBaseDataPoint *sfxpb.DataPoint,
	count *int64,
) *sfxpb.DataPoint {

	totalCountDP := *sfxBaseDataPoint
	totalCountName := sfxBaseDataPoint.Metric + "_count"
	totalCountDP.Metric = totalCountName
	totalCountDP.MetricType = totalCountMetricType
	totalCountDP.Value = sfxpb.Datum{IntValue: count}

	return &totalCountDP
}

func buildSumDataPoint(
	sfxBaseDataPoint *sfxpb.DataPoint,
	sum *float64,
) *sfxpb.DataPoint {

	sumDP := *sfxBaseDataPoint
	sumDP.MetricType = sumMetricType
	sumDP.Value = sfxpb.Datum{DoubleValue: sum}

	return &sumDP
}

// sanitizeDataPointLabels replaces all characters unsupported by SignalFx backend
// in metric label keys and with "_"
func sanitizeDataPointDimensions(dps []*sfxpb.DataPoint) {
	for _, dp := range dps {
		for _, d := range dp.Dimensions {
			d.Key = filterKeyChars(d.Key)
		}
	}
}

func filterKeyChars(str string) string {
	filterMap := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}

	return strings.Map(filterMap, str)
}

func float64ToDimValue(f float64) string {
	// Parameters below are the same used by Prometheus
	// see https://github.com/prometheus/common/blob/b5fe7d854c42dc7842e48d1ca58f60feae09d77b/expfmt/text_create.go#L450
	// SignalFx agent uses a different pattern
	// https://github.com/signalfx/signalfx-agent/blob/5779a3de0c9861fa07316fd11b3c4ff38c0d78f0/internal/monitors/prometheusexporter/conversion.go#L77
	// The important issue here is consistency with the exporter, opting for the
	// more common one used by Prometheus.
	return strconv.FormatFloat(f, 'g', -1, 64)
}
