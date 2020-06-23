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

	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

func TestPatternProto(t *testing.T) {
	pattern := Pattern{
		StartsWith: "/my/metric",
	}

	p := pattern.Proto()

	if p.Match.(*pb.ConfigResponse_MetricConfig_Schedule_Pattern_StartsWith).StartsWith != "/my/metric" {
		t.Errorf("improper conversion to proto")
	}
}

func TestPatternDuplicateProto(t *testing.T) {
	pattern := Pattern{
		Equals:     "/use/this/rule",
		StartsWith: "/not/this/one",
	}

	p := pattern.Proto()

	if p.Match.(*pb.ConfigResponse_MetricConfig_Schedule_Pattern_Equals).Equals != "/use/this/rule" {
		t.Errorf("improper conversion to proto")
	}
}

func TestPatternHash(t *testing.T) {
	configA := Pattern{
		Equals:     "/use/this/rule",
		StartsWith: "/not/this/one",
	}

	configB := Pattern{
		Equals: "/use/this/rule",
	}

	configC := Pattern{
		StartsWith: "/use/this/rule",
	}

	if !bytes.Equal(configA.Hash(), configB.Hash()) {
		t.Errorf("identical configs with different hashes")
	}

	if bytes.Equal(configA.Hash(), configC.Hash()) {
		t.Errorf("different configs with identical hashes")
	}
}
