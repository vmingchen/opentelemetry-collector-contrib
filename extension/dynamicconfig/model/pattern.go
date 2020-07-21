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
//
// Contains common models for the dynamic config service. The corresponding
// Proto() methods convert the model representation to a usable struct for
// protobuf marshalling.

package model

import (
	"fmt"
	"hash/fnv"

	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/experimental/metricconfigservice"
)

type Pattern struct {
	Equals     string
	StartsWith string
}

func (p *Pattern) Proto() (*pb.MetricConfigResponse_Schedule_Pattern, error) {
	if len(p.Equals) > 0 {
		if len(p.StartsWith) > 0 {
			return nil, fmt.Errorf("only specify StartsWith or Equals, not both")
		}

		return &pb.MetricConfigResponse_Schedule_Pattern{
			Match: &pb.MetricConfigResponse_Schedule_Pattern_Equals{
				Equals: p.Equals,
			},
		}, nil
	} else {
		return &pb.MetricConfigResponse_Schedule_Pattern{
			Match: &pb.MetricConfigResponse_Schedule_Pattern_StartsWith{
				StartsWith: p.StartsWith,
			},
		}, nil
	}
}

func (p *Pattern) Hash() []byte {
	hasher := fnv.New64a()

	if len(p.Equals) > 0 {
		hasher.Write([]byte("Equals"))
		hasher.Write([]byte(p.Equals))
	} else {
		hasher.Write([]byte("StartsWith"))
		hasher.Write([]byte(p.StartsWith))
	}

	return hasher.Sum(nil)
}
