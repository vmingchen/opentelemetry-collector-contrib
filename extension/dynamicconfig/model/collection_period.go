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
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
	"hash/fnv"
)

type CollectionPeriod string

func (period CollectionPeriod) Proto() pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod {
	interval := pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod_value[string(period)]
	return pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod(interval)
}

func (period CollectionPeriod) Hash() []byte {
	hasher := fnv.New64a()
	hasher.Write([]byte(period.Proto().String()))
	return hasher.Sum(nil)

}
