// Copyright 2019, OpenTelemetry Authors
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

package sapmreceiver

// This file implements factory for SAPM receiver.

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// The value of "type" key in configuration.
	typeStr = "sapm"

	// Default endpoints to bind to.
	defaultEndpoint = ":7276"
)

// Factory is the factory for SAPM receiver.
type Factory struct {
}

// Type gets the type of the Receiver config created by this factory.
func (f *Factory) Type() configmodels.Type {
	return configmodels.Type(typeStr)
}

// CustomUnmarshaler returns nil because we don't need custom unmarshaling for this config.
func (f *Factory) CustomUnmarshaler() component.CustomUnmarshaler {
	return nil
}

// CreateDefaultConfig creates the default configuration for SAPM receiver.
func (f *Factory) CreateDefaultConfig() configmodels.Receiver {
	return &Config{
		ReceiverSettings: configmodels.ReceiverSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		Endpoint: defaultEndpoint,
	}
}

// extract the port number from string in "address:port" format. If the
// port number cannot be extracted returns an error.
// TODO make this a utility function
func extractPortFromEndpoint(endpoint string) (int, error) {
	_, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return 0, fmt.Errorf("endpoint is not formatted correctly: %s", err.Error())
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("endpoint port is not a number: %s", err.Error())
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port number must be between 1 and 65535")
	}
	return int(port), nil
}

// verify that the configured port is not 0
func (rCfg *Config) validate() error {
	_, err := extractPortFromEndpoint(rCfg.Endpoint)
	if err != nil {
		return err
	}
	return nil
}

// CreateTraceReceiver creates a trace receiver based on provided config.
func (f *Factory) CreateTraceReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	nextConsumer consumer.TraceConsumer,
) (component.TraceReceiver, error) {
	// assert config is SAPM config
	rCfg := cfg.(*Config)

	err := rCfg.validate()
	if err != nil {
		return nil, err
	}

	// Create the receiver.
	return New(ctx, params, rCfg, nextConsumer)
}

// CreateMetricsReceiver creates a metrics receiver based on provided config.
func (f *Factory) CreateMetricsReceiver(
	_ context.Context,
	_ component.ReceiverCreateParams,
	_ configmodels.Receiver,
	_ consumer.MetricsConsumer,
) (component.MetricsReceiver, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}
