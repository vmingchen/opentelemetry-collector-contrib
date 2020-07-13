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

package simpleprometheusreceiver

import (
	"time"

	"go.opentelemetry.io/collector/config/configmodels"
)

// Config defines configuration for simple prometheus receiver.
type Config struct {
	configmodels.ReceiverSettings `mapstructure:",squash"`
	httpConfig                    `mapstructure:",squash"`
	// TODO: Use one of the configs from core.
	// The target endpoint.
	Endpoint string `mapstructure:"endpoint"`
	// CollectionInterval is the interval at which metrics should be collected
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
	// MetricsPath the path to the metrics endpoint.
	MetricsPath string `mapstructure:"metrics_path"`
	// Whether or not to use pod service account to authenticate.
	UseServiceAccount bool `mapstructure:"use_service_account"`
}

// TODO: Move to a common package for use by other receivers and also pull
// in other utilities from
// https://github.com/signalfx/signalfx-agent/blob/master/pkg/core/common/httpclient/http.go.
type httpConfig struct {
	// Whether not TLS is enabled
	TLSEnabled bool      `mapstructure:"tls_enabled"`
	TLSConfig  tlsConfig `mapstructure:"tls_config"`
}

// tlsConfig holds common TLS config options
type tlsConfig struct {
	// Path to the CA cert that has signed the TLS cert.
	CAFile string `mapstructure:"ca_file"`
	// Path to the client TLS cert to use for TLS required connections.
	CertFile string `mapstructure:"cert_file"`
	// Path to the client TLS key to use for TLS required connections.
	KeyFile string `mapstructure:"key_file"`
	// Whether or not to verify the exporter's TLS cert.
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}
