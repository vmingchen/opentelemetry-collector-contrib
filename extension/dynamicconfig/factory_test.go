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
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcheck"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/testutils"
)

func TestFactory_Type(t *testing.T) {
	factory := Factory{}
	require.Equal(t, configmodels.Type(typeStr), factory.Type())
}

func TestFactory_CreateDefaultConfig(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	assert.Equal(t, &Config{
		ExtensionSettings: configmodels.ExtensionSettings{
			NameVal: typeStr,
			TypeVal: typeStr,
		},
		Endpoint:        "0.0.0.0:55700",
		Target:          "",
		LocalConfigFile: "dynamic-config-local-schedules.yaml",
		WaitTime:        30,
	}, cfg)

	assert.NoError(t, configcheck.ValidateConfig(cfg))
	ext, err := factory.CreateExtension(context.Background(), component.ExtensionCreateParams{Logger: zap.NewNop()}, cfg)
	require.NoError(t, err)
	require.NotNil(t, ext)

	// Restore instance tracking from factory, for other tests.
	atomic.StoreInt32(&instanceState, instanceNotCreated)
}

func TestFactory_CreateExtension(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.Endpoint = testutils.GetAvailableLocalAddress(t)

	ext, err := factory.CreateExtension(context.Background(), component.ExtensionCreateParams{Logger: zap.NewNop()}, cfg)
	require.NoError(t, err)
	require.NotNil(t, ext)

	// Restore instance tracking from factory, for other tests.
	atomic.StoreInt32(&instanceState, instanceNotCreated)
}

func TestFactory_CreateExtensionOnlyOnce(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.Endpoint = testutils.GetAvailableLocalAddress(t)

	ext, err := factory.CreateExtension(context.Background(), component.ExtensionCreateParams{Logger: zap.NewNop()}, cfg)
	require.NoError(t, err)
	require.NotNil(t, ext)

	ext1, err := factory.CreateExtension(context.Background(), component.ExtensionCreateParams{Logger: zap.NewNop()}, cfg)
	require.Error(t, err)
	require.Nil(t, ext1)

	// Restore instance tracking from factory, for other tests.
	atomic.StoreInt32(&instanceState, instanceNotCreated)
}
