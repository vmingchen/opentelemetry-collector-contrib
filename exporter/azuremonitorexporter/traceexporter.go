// Copyright OpenTelemetry Authors
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

package azuremonitorexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.uber.org/zap"
)

type traceExporter struct {
	config           *Config
	transportChannel transportChannel
	logger           *zap.Logger
}

type traceVisitor struct {
	processed int
	err       error
	exporter  *traceExporter
}

// Called for each tuple of Resource, InstrumentationLibrary, and Span
func (v *traceVisitor) visit(
	resource pdata.Resource,
	instrumentationLibrary pdata.InstrumentationLibrary, span pdata.Span) (ok bool) {

	envelope, err := spanToEnvelope(resource, instrumentationLibrary, span, v.exporter.logger)
	if err != nil {
		// record the error and short-circuit
		v.err = err
		return false
	}

	// apply the instrumentation key to the envelope
	envelope.IKey = v.exporter.config.InstrumentationKey

	// This is a fire and forget operation
	v.exporter.transportChannel.Send(envelope)
	v.processed++

	return true
}

func (exporter *traceExporter) onTraceData(context context.Context, traceData pdata.Traces) (droppedSpans int, err error) {
	spanCount := traceData.SpanCount()
	if spanCount == 0 {
		return 0, nil
	}

	visitor := &traceVisitor{exporter: exporter}
	Accept(traceData, visitor)
	return (spanCount - visitor.processed), visitor.err
}

// Returns a new instance of the trace exporter
func newTraceExporter(config *Config, transportChannel transportChannel, logger *zap.Logger) (component.TraceExporter, error) {

	exporter := &traceExporter{
		config:           config,
		transportChannel: transportChannel,
		logger:           logger,
	}

	return exporterhelper.NewTraceExporter(config, exporter.onTraceData)
}
