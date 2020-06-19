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

package dynamicconfig

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configmodels"
)

func TestLoadConfig(t *testing.T) {
	factories, err := config.ExampleComponents()
	assert.NoError(t, err)

	factory := &Factory{}
	factories.Extensions[typeStr] = factory
	cfg, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	require.Nil(t, err)
	require.NotNil(t, cfg)

	ext0 := cfg.Extensions["dynamicconfig"]
	assert.Equal(t, factory.CreateDefaultConfig(), ext0)

	ext1 := cfg.Extensions["dynamicconfig/1"]
	assert.Equal(t,
		&Config{
			ExtensionSettings: configmodels.ExtensionSettings{
				TypeVal: "dynamicconfig",
				NameVal: "dynamicconfig/1",
			},
			Endpoint:        "localhost:12345",
			LocalConfigFile: "wut-the-wut.yaml",
		},
		ext1,
	)

	assert.Equal(t, 1, len(cfg.Service.Extensions))
	assert.Equal(t, "dynamicconfig/1", cfg.Service.Extensions[0])
}
