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
	"context"
	"errors"

	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

// ConfigService implements the server side of the gRPC service for config
// updates.
type ConfigService struct {
	pb.UnimplementedDynamicConfigServer // for forward compatability
	backend                             ConfigBackend
}

func NewConfigService(opts ...Option) (*ConfigService, error) {
	service := &ConfigService{}
	for _, opt := range opts {
		if err := opt(service); err != nil {
			return nil, err
		}
	}

	if service.backend == nil {
		return nil, errors.New("config service is missing backend")
	}

	return service, nil
}

type Option func(*ConfigService) error

func WithLocalConfig(filepath string) Option {
	return func(service *ConfigService) error {
		backend, err := NewLocalConfigBackend(filepath)
		service.backend = backend
		return err
	}
}

func (service *ConfigService) GetConfig(ctx context.Context, req *pb.ConfigRequest) (*pb.ConfigResponse, error) {
	var resp *pb.ConfigResponse
	if service.backend.IsSameFingerprint(req.LastKnownFingerprint) {
		resp = &pb.ConfigResponse{Fingerprint: service.backend.GetFingerprint()}
	} else {
		resp = service.backend.BuildConfigResponse()
	}

	return resp, nil
}
