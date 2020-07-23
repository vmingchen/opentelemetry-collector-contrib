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

package awsxrayreceiver

import (
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestLoadConfig(t *testing.T) {
	factories, err := config.ExampleComponents()
	assert.Nil(t, err)

	factory := &Factory{}
	factories.Receivers[configmodels.Type(typeStr)] = factory
	cfg, err := config.LoadConfigFile(
		t, path.Join(".", "testdata", "config.yaml"), factories,
	)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, len(cfg.Receivers), 3)

	// ensure default configurations are generated when users provide
	// nothing.
	r0 := cfg.Receivers[typeStr]
	assert.Equal(t, factory.CreateDefaultConfig(), r0)

	// ensure the UDP endpoint can be properly overwritten
	r1 := cfg.Receivers[typeStr+"/udp_endpoint"].(*Config)
	assert.Equal(t,
		&Config{
			ReceiverSettings: configmodels.ReceiverSettings{
				TypeVal: configmodels.Type(typeStr),
				NameVal: typeStr + "/udp_endpoint",
			},
			TCPAddr: confignet.TCPAddr{
				Endpoint: "localhost:5678",
			},
			ProxyServer: &proxyServer{
				TCPEndpoint:  "0.0.0.0:2000",
				ProxyAddress: "",
				TLSSetting: configtls.TLSClientSetting{
					Insecure:   false,
					ServerName: "",
				},
				Region:      "",
				RoleARN:     "",
				AWSEndpoint: "",
				LocalMode:   aws.Bool(false),
			},
		},
		r1)

	// ensure the fields under proxy_server are properly overwritten
	r2 := cfg.Receivers[typeStr+"/proxy_server"].(*Config)
	assert.Equal(t,
		&Config{
			ReceiverSettings: configmodels.ReceiverSettings{
				TypeVal: configmodels.Type(typeStr),
				NameVal: typeStr + "/proxy_server",
			},
			TCPAddr: confignet.TCPAddr{
				Endpoint: "0.0.0.0:2000",
			},
			ProxyServer: &proxyServer{
				TCPEndpoint:  "localhost:1234",
				ProxyAddress: "https://proxy.proxy.com",
				TLSSetting: configtls.TLSClientSetting{
					Insecure:   true,
					ServerName: "something",
				},
				Region:      "us-west-1",
				RoleARN:     "arn:aws:iam::123456789012:role/awesome_role",
				AWSEndpoint: "https://another.aws.endpoint.com",
				LocalMode:   aws.Bool(true),
			},
		},
		r2)
}
