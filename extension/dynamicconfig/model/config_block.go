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

type ConfigBlock struct {
	Resource     []string
	MetricConfig *MetricConfig
	TraceConfig  *TraceConfig
}

func (block *ConfigBlock) Add(other *ConfigBlock) {
	if block.MetricConfig == nil {
		block.MetricConfig = &MetricConfig{}
	}

	block.MetricConfig.Schedules = append(
		block.MetricConfig.Schedules,
		other.MetricConfig.Schedules...)
}

func (block *ConfigBlock) Hash() []byte {
	hashes := [][]byte{
		block.MetricConfig.Hash(),
		block.TraceConfig.Hash(),
	}

	return combineHash(hashes)
}
