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

package dynamicconfigextension

import (
	"go.opentelemetry.io/collector/config/configmodels"
)

// Config has the configuration for the extension enabling the dynamic
// configuration service
type Config struct {
	configmodels.ExtensionSettings `mapstructure:",squash"`

	// Endpoint is the address and port used to communicate the config updates
	// The default value is localhost:55700.
	Endpoint string `mapstructure:"endpoint"`

	// LocalConfigFile is the local record of configuration updates, applied
	// when a third-party config service backend is not used.
	LocalConfigFile string `mapstructure:"local_config_file"`
}
