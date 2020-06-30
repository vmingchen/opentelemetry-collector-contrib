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

	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

// ConfigBackend defines a general backend that the service can read
// configuration data from.
type ConfigBackend interface {
	GetFingerprint(*res.Resource) ([]byte, error)
	BuildConfigResponse(*res.Resource) (*pb.ConfigResponse, error)
	Close() error
}

// ConfigService implements the server side of the gRPC service for config
// updates.
type ConfigService struct {
	pb.UnimplementedDynamicConfigServer // for forward compatability
	backend                             ConfigBackend
}

func NewConfigService(opts ...Option) (*ConfigService, error) {
	builder := &serviceBuilder{}
	for _, opt := range opts {
		opt(builder)
	}

	backend, err := builder.build()
	if err != nil {
		return nil, err
	}

	return &ConfigService{backend: backend}, nil
}

type serviceBuilder struct {
	remoteConfigAddress string
	filepath            string
	updateStrategy      UpdateStrategy
	waitTime            int32

	// overrides build() to use this given backend.
	// NOTE: intended for testing only!
	backend ConfigBackend
}

// TODO: implement LocalConfigBackend as fall-back
func (builder *serviceBuilder) build() (ConfigBackend, error) {
	if builder.backend != nil {
		return builder.backend, nil
	}

	if builder.remoteConfigAddress != "" {
		backend, err := NewRemoteConfigBackend(builder.remoteConfigAddress)
		if err != nil {
			return nil, err
		}

		if builder.updateStrategy != 0 {
			backend.SetUpdateStrategy(builder.updateStrategy)
		}

		return backend, nil
	}

	if builder.filepath != "" {
		backend, err := NewLocalConfigBackend(builder.filepath)
		if err != nil {
			return nil, err
		}

		if builder.waitTime > 0 {
			backend.waitTime = builder.waitTime
		}

		return backend, nil

	}

	return nil, errors.New("missing backend specification")
}

type Option func(*serviceBuilder)

func WithRemoteConfig(remoteConfigAddress string) Option {
	return func(builder *serviceBuilder) {
		builder.remoteConfigAddress = remoteConfigAddress
	}
}

func WithUpdateStrategy(strategy UpdateStrategy) Option {
	return func(builder *serviceBuilder) {
		builder.updateStrategy = strategy
	}
}

func WithLocalConfig(filepath string) Option {
	return func(builder *serviceBuilder) {
		builder.filepath = filepath
	}
}

func WithWaitTime(time int32) Option {
	return func(builder *serviceBuilder) {
		builder.waitTime = time
	}
}

// TODO: Match req.Resource to appropriate configs
func (service *ConfigService) GetConfig(ctx context.Context, req *pb.ConfigRequest) (*pb.ConfigResponse, error) {
	var resp *pb.ConfigResponse
	backendFingerprint, err := service.backend.GetFingerprint(req.Resource)
	if err != nil {
		return nil, fmt.Errorf("fail to read fingerprint from backend: %w", err)
	}

	if bytes.Equal(backendFingerprint, req.LastKnownFingerprint) {
		resp = &pb.ConfigResponse{Fingerprint: backendFingerprint}
	} else {
		resp, err = service.backend.BuildConfigResponse(req.Resource)
		if err != nil {
			return nil, fmt.Errorf("backend fail to build config response: %w", err)
		}
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
