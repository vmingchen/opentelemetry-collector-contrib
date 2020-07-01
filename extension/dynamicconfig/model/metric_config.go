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
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

type MetricConfig struct {
	Schedules []Schedule
}

func (config *MetricConfig) Proto() *pb.ConfigResponse_MetricConfig {
	scheduleSlice := make([]*pb.ConfigResponse_MetricConfig_Schedule, len(config.Schedules))
	for i, schedule := range config.Schedules {
		scheduleSlice[i] = schedule.Proto()
	}

	proto := &pb.ConfigResponse_MetricConfig{
		Schedules: scheduleSlice,
	}

	return proto
}

func (config *MetricConfig) Hash() []byte {
	hashes := make([][]byte, len(config.Schedules))
	for i, sched := range config.Schedules {
		hashes[i] = sched.Hash()
	}

	return combineHash(hashes)
}
