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
	"fmt"
	"strings"

	com "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
)

type Config struct {
	ConfigBlocks []*ConfigBlock
}

func (config *Config) Match(resource *res.Resource) *ConfigBlock {
	resourceSet, resourceList := embed(resource)
	totalBlock := &ConfigBlock{
		MetricConfig: &MetricConfig{},
		Resource:     resourceList,
	}

	for _, block := range config.ConfigBlocks {
		if doInclude(block, resourceSet) {
			totalBlock.Add(block)
		}
	}

	return totalBlock
}

func embed(resource *res.Resource) (map[string]bool, []string) {
	if resource == nil {
		resource = &res.Resource{}
	}

	resourceSet := make(map[string]bool)
	resourceList := make([]string, len(resource.Attributes))

	for i, attr := range resource.Attributes {
		attrString := attrToString(attr)
		resourceSet[attrString] = true
		resourceList[i] = attrString
	}

	return resourceSet, resourceList
}

func attrToString(attr *com.KeyValue) string {
	rawValue := attr.Value.String()
	value := strings.Split(rawValue, ":")[1]
	attrString := clean(fmt.Sprintf("%s:%s", attr.Key, value))

	return attrString
}

func doInclude(block *ConfigBlock, resourceSet map[string]bool) bool {
	include := true
	for _, label := range block.Resource {
		label = clean(label)
		include = include && resourceSet[label]
	}

	return include
}

func clean(label string) string {
	label = strings.ReplaceAll(label, " ", "")
	label = strings.ReplaceAll(label, `"`, "")

	return label
}
