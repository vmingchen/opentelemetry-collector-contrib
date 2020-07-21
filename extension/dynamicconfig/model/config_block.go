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
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/experimental/metricconfigservice"
)

type ConfigBlock struct {
	Resource  []string
	Schedules []*Schedule
}

func (block *ConfigBlock) Proto() ([]*pb.MetricConfigResponse_Schedule, error) {
	scheduleSlice := make([]*pb.MetricConfigResponse_Schedule, len(block.Schedules))

	var err error
	for i, schedule := range block.Schedules {
		scheduleSlice[i], err = schedule.Proto()
		if err != nil {
			return nil, err
		}
	}

	return scheduleSlice, nil
}

func (block *ConfigBlock) Hash() []byte {
	if len(block.Schedules) == 0 {
		return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}

	hashes := make([][]byte, len(block.Schedules))
	for i, sched := range block.Schedules {
		hashes[i] = sched.Hash()
	}

	return combineHash(hashes)
}

func (block *ConfigBlock) Add(other *ConfigBlock) {
	block.Schedules = append(
		block.Schedules,
		other.Schedules...)
}
