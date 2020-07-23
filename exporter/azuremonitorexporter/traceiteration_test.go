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
	"testing"

	mock "github.com/stretchr/testify/mock"
	"go.opentelemetry.io/collector/consumer/pdata"
)

type mockVisitor struct {
	mock.Mock
}

func (v *mockVisitor) visit(resource pdata.Resource, instrumentationLibrary pdata.InstrumentationLibrary, span pdata.Span) (ok bool) {
	args := v.Called(resource, instrumentationLibrary, span)
	return args.Bool(0)
}

// Tests the iteration logic over a pdata.Traces type when no ResourceSpans are provided
func TestTraceDataIterationNoResourceSpans(t *testing.T) {
	traces := pdata.NewTraces()

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic over a pdata.Traces type when a ResourceSpans is nil
func TestTraceDataIterationResourceSpansIsNil(t *testing.T) {
	traces := pdata.NewTraces()
	resourceSpans := pdata.NewResourceSpans()
	traces.ResourceSpans().Append(&resourceSpans)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic over a pdata.Traces type when a Resource is nil
func TestTraceDataIterationResourceIsNil(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic over a pdata.Traces type when InstrumentationLibrarySpans is nil
func TestTraceDataIterationInstrumentationLibrarySpansIsNil(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rs := traces.ResourceSpans().At(0)
	r := rs.Resource()
	r.InitEmpty()
	instrumentationLibrarySpans := pdata.NewInstrumentationLibrarySpans()
	rs.InstrumentationLibrarySpans().Append(&instrumentationLibrarySpans)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic over a pdata.Traces type when there are no Spans
func TestTraceDataIterationNoSpans(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rs := traces.ResourceSpans().At(0)
	r := rs.Resource()
	r.InitEmpty()
	instrumentationLibrarySpans := pdata.NewInstrumentationLibrarySpans()
	instrumentationLibrarySpans.InitEmpty()
	rs.InstrumentationLibrarySpans().Append(&instrumentationLibrarySpans)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic over a pdata.Traces type when the Span is nil
func TestTraceDataIterationSpanIsNil(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rs := traces.ResourceSpans().At(0)
	r := rs.Resource()
	r.InitEmpty()
	rs.InstrumentationLibrarySpans().Resize(1)
	ilss := rs.InstrumentationLibrarySpans().At(0)
	span := pdata.NewSpan()
	ilss.Spans().Append(&span)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 0)
}

// Tests the iteration logic if the visitor returns true
func TestTraceDataIterationNoShortCircuit(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rs := traces.ResourceSpans().At(0)
	r := rs.Resource()
	r.InitEmpty()
	rs.InstrumentationLibrarySpans().Resize(1)
	ilss := rs.InstrumentationLibrarySpans().At(0)
	ilss.Spans().Resize(2)

	visitor := getMockVisitor(true)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 2)
}

// Tests the iteration logic short circuit if the visitor returns false
func TestTraceDataIterationShortCircuit(t *testing.T) {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rs := traces.ResourceSpans().At(0)
	r := rs.Resource()
	r.InitEmpty()
	rs.InstrumentationLibrarySpans().Resize(1)
	ilss := rs.InstrumentationLibrarySpans().At(0)
	ilss.Spans().Resize(2)

	visitor := getMockVisitor(false)

	Accept(traces, visitor)

	visitor.AssertNumberOfCalls(t, "visit", 1)
}

func getMockVisitor(returns bool) *mockVisitor {
	visitor := new(mockVisitor)
	visitor.On("visit", mock.Anything, mock.Anything, mock.Anything).Return(returns)
	return visitor
}
