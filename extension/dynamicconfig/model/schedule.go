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
	"hash/fnv"

	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/experimental/metricconfigservice"
)

type Schedule struct {
	InclusionPatterns []Pattern
	ExclusionPatterns []Pattern
	PeriodSec         CollectionPeriod
}

func (schedule *Schedule) Proto() (*pb.MetricConfigResponse_Schedule, error) {
	incSlice := make([]*pb.MetricConfigResponse_Schedule_Pattern, len(schedule.InclusionPatterns))
	excSlice := make([]*pb.MetricConfigResponse_Schedule_Pattern, len(schedule.ExclusionPatterns))

	var err error
	for i, incPat := range schedule.InclusionPatterns {
		incSlice[i], err = incPat.Proto()
		if err != nil {
			return nil, err
		}
	}

	for i, excPat := range schedule.ExclusionPatterns {
		excSlice[i], err = excPat.Proto()
		if err != nil {
			return nil, err
		}
	}

	periodProto, err := schedule.PeriodSec.Proto()
	if err != nil {
		return nil, err
	}

	proto := &pb.MetricConfigResponse_Schedule{
		InclusionPatterns: incSlice,
		ExclusionPatterns: excSlice,
		PeriodSec:         periodProto,
	}

	return proto, nil
}

func (schedule *Schedule) Hash() []byte {
	incHashes := make([][]byte, len(schedule.InclusionPatterns))
	excHashes := make([][]byte, len(schedule.ExclusionPatterns))

	for i, incPat := range schedule.InclusionPatterns {
		incHashes[i] = incPat.Hash()
	}

	for i, excPat := range schedule.ExclusionPatterns {
		excHashes[i] = excPat.Hash()
	}

	hashes := [][]byte{
		[]byte("InclusionPatterns"),
		combineHash(incHashes),
		[]byte("ExclusionPatterns"),
		combineHash(excHashes),
		schedule.PeriodSec.Hash(),
	}

	hasher := fnv.New64a()
	for _, val := range hashes {
		hasher.Write(val)
	}

	return hasher.Sum(nil)
}
