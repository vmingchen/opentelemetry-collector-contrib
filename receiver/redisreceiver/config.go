// Copyright 2020, OpenTelemetry Authors
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

package redisreceiver

import (
	"time"

	"go.opentelemetry.io/collector/config/configmodels"
)

type config struct {
	configmodels.ReceiverSettings `mapstructure:",squash"`
	// TODO: Use one of the configs from core.
	// The target endpoint.
	Endpoint string `mapstructure:"endpoint"`
	// The duration between Redis metric fetches.
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
	// The logical name of the Redis server. This value will be added as a
	// "service.name" Resource label.
	ServiceName string `mapstructure:"service_name"`

	// TODO allow users to add additional resource key value pairs?

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string `mapstructure:"password"`
}
