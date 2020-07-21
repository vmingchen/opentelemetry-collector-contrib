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

package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service/file"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service/mock"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service/remote"
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/experimental/metricconfigservice"
	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
)

// ConfigBackend defines a general backend that the service can read
// configuration data from.
type ConfigBackend interface {
	BuildConfigResponse(*res.Resource) (*pb.MetricConfigResponse, error)
	Close() error
}

// ConfigService implements the server side of the gRPC service for config
// updates.
type ConfigService struct {
	pb.UnimplementedMetricConfigServer // for forward compatability
	backend                             ConfigBackend
}

func NewConfigService(opts ...Option) (*ConfigService, error) {
	builder := &ServiceBuilder{}
	for _, opt := range opts {
		opt(builder)
	}

	backend, err := builder.build()
	if err != nil {
		return nil, err
	}

	return &ConfigService{backend: backend}, nil
}

type ServiceBuilder struct {
	remoteConfigAddress string
	filepath            string
	waitTime            int32

	// overrides build() to use this given backend.
	// NOTE: intended for testing only!
	backend ConfigBackend
}

func (builder *ServiceBuilder) build() (ConfigBackend, error) {
	if builder.backend != nil {
		return builder.backend, nil
	}

	if builder.remoteConfigAddress != "" {
		backend, err := remote.NewBackend(builder.remoteConfigAddress)
		if err != nil {
			return nil, err
		}

		return backend, nil
	}

	if builder.filepath != "" {
		backend, err := file.NewBackend(builder.filepath)
		if err != nil {
			return nil, err
		}

		if builder.waitTime > 0 {
			backend.SetWaitTime(builder.waitTime)
		}

		return backend, nil

	}

	return nil, errors.New("missing backend specification")
}

type Option func(*ServiceBuilder)

func WithRemoteConfig(remoteConfigAddress string) Option {
	return func(builder *ServiceBuilder) {
		builder.remoteConfigAddress = remoteConfigAddress
	}
}

func WithLocalConfig(filepath string) Option {
	return func(builder *ServiceBuilder) {
		builder.filepath = filepath
	}
}

func WithWaitTime(time int32) Option {
	return func(builder *ServiceBuilder) {
		builder.waitTime = time
	}
}

// NOTE: intended for testing only!
func WithMockBackend() Option {
	return func(builder *ServiceBuilder) {
		builder.backend = &mock.Backend{}
	}
}

func (service *ConfigService) GetMetricConfig(ctx context.Context, req *pb.MetricConfigRequest) (*pb.MetricConfigResponse, error) {
	resp, err := service.backend.BuildConfigResponse(req.Resource)
	if err != nil {
		return nil, fmt.Errorf("backend failed to build config response: %w", err)
	}

	if bytes.Equal(resp.Fingerprint, req.LastKnownFingerprint) {
		resp = &pb.MetricConfigResponse{Fingerprint: resp.Fingerprint}
	}

	return resp, nil
}

func (service *ConfigService) Stop() error {
	if service != nil {
		if err := service.backend.Close(); err != nil {
			return fmt.Errorf("fail to stop config service: %w", err)
		}
	}

	return nil
}
