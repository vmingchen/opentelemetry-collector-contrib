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
)

func TestAddConfigBlock(t *testing.T) {
	configBlocks := []*ConfigBlock{
		{
			MetricConfig: &MetricConfig{
				Schedules: []*Schedule{
					{Period: "SEC_1"},
				},
			},
		},
		{
			MetricConfig: &MetricConfig{
				Schedules: []*Schedule{
					{Period: "SEC_5"}, {Period: "DAY_1"},
				},
			},
		},
		{
			MetricConfig: &MetricConfig{
				Schedules: []*Schedule{},
			},
		},
	}

	var totalBlock ConfigBlock
	for _, block := range configBlocks {
		totalBlock.Add(block)
	}

	scheds := totalBlock.MetricConfig.Schedules
	if len(scheds) != 3 {
		t.Errorf("expected 3 schedules, found: %v", len(scheds))
	}

	if scheds[0].Period != "SEC_1" || scheds[1].Period != "SEC_5" || scheds[2].Period != "DAY_1" {
		t.Errorf("expected periods SEC_1, SEC_5, DAY_1, found: %v", scheds)
	}
}

func TestHash(t *testing.T) {
	configA := ConfigBlock{
		MetricConfig: &MetricConfig{
			Schedules: []*Schedule{
				{Period: "SEC_1"},
			},
		},
		TraceConfig: &TraceConfig{},
	}

	configB := ConfigBlock{
		TraceConfig: &TraceConfig{},
		MetricConfig: &MetricConfig{
			Schedules: []*Schedule{
				{Period: "SEC_1"},
			},
		},
	}

	configC := ConfigBlock{
		MetricConfig: &MetricConfig{
			Schedules: []*Schedule{
				{Period: "SEC_5"},
			},
		},
		TraceConfig: &TraceConfig{},
	}

	if !bytes.Equal(configA.Hash(), configB.Hash()) {
		t.Errorf("identical configs with different hashes")
	}

	if bytes.Equal(configA.Hash(), configC.Hash()) {
		t.Errorf("different configs with identical hashes")
	}
}
