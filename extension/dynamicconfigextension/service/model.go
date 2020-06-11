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

package service

import (
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

// ConfigBackend defines a general backend that the service can read
// configuration data from.
type ConfigBackend interface {
	GetFingerprint() []byte
	IsSameFingerprint(fingerprint []byte) bool
	BuildConfigResponse() *pb.ConfigResponse
}

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

type Schedule struct {
	InclusionPatterns []Pattern
	ExclusionPatterns []Pattern
	Period            CollectionPeriod
	Metadata          []byte
}

func (schedule *Schedule) Proto() *pb.ConfigResponse_MetricConfig_Schedule {
	incSlice := make([]*pb.ConfigResponse_MetricConfig_Schedule_Pattern, len(schedule.InclusionPatterns))
	excSlice := make([]*pb.ConfigResponse_MetricConfig_Schedule_Pattern, len(schedule.ExclusionPatterns))

	for i, incPat := range schedule.InclusionPatterns {
		incSlice[i] = incPat.Proto()
	}

	for i, excPat := range schedule.ExclusionPatterns {
		excSlice[i] = excPat.Proto()
	}

	proto := &pb.ConfigResponse_MetricConfig_Schedule{
		InclusionPatterns: incSlice,
		ExclusionPatterns: excSlice,
		Period:            schedule.Period.Proto(),
		Metadata:          schedule.Metadata,
	}

	return proto
}

type Pattern struct {
	Equals     string
	StartsWith string
}

func (p *Pattern) Proto() *pb.ConfigResponse_MetricConfig_Schedule_Pattern {
	if len(p.Equals) > 0 {
		return &pb.ConfigResponse_MetricConfig_Schedule_Pattern{
			Match: &pb.ConfigResponse_MetricConfig_Schedule_Pattern_Equals{
				Equals: p.Equals,
			},
		}
	} else {
		return &pb.ConfigResponse_MetricConfig_Schedule_Pattern{
			Match: &pb.ConfigResponse_MetricConfig_Schedule_Pattern_StartsWith{
				StartsWith: p.StartsWith,
			},
		}
	}
}

type CollectionPeriod string

func (period CollectionPeriod) Proto() pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod {
	interval := pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod_value[string(period)]
	return pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod(interval)
}
