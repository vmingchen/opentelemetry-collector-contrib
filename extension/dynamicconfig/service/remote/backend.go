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

package remote

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

type UpdateStrategy uint8

const (
	// updates cached response on GetFingerprint() and BuildConfigResponse()
	Default UpdateStrategy = iota

	// updates cached response only on GetFingerprint()
	OnGetFingerprint
)

type Backend struct {
	remoteConfigAddress string
	conn                *grpc.ClientConn
	client              pb.DynamicConfigClient
	updateStrategy      UpdateStrategy

	mu   sync.Mutex
	resp *pb.ConfigResponse
}

func NewBackend(remoteConfigAddress string) (*Backend, error) {
	backend := &Backend{
		remoteConfigAddress: remoteConfigAddress,
		conn:                nil,
		client:              nil,
		updateStrategy:      Default,
	}

	if err := backend.initConn(); err != nil {
		return nil, err
	}

	return backend, nil
}

func (backend *Backend) initConn() error {
	conn, err := grpc.Dial(
		backend.remoteConfigAddress,
		grpc.WithInsecure(), // TODO: consider security implications
	)
	if err != nil {
		return fmt.Errorf("remote config backend fail to connect: %w", err)
	}

	backend.conn = conn
	backend.client = pb.NewDynamicConfigClient(conn)
	return nil
}

func (backend *Backend) GetUpdateStrategy() UpdateStrategy {
	return backend.updateStrategy
}

func (backend *Backend) SetUpdateStrategy(strategy UpdateStrategy) {
	if strategy == Default || strategy == OnGetFingerprint {
		backend.updateStrategy = strategy
	}
}

func (backend *Backend) GetFingerprint(resource *res.Resource) ([]byte, error) {
	if err := backend.syncRemote(resource); err != nil {
		return nil, fmt.Errorf("fail to get fingerprint: %w", err)
	}

	backend.mu.Lock()
	resp := backend.resp
	backend.mu.Unlock()

	return resp.Fingerprint, nil
}

func (backend *Backend) BuildConfigResponse(resource *res.Resource) (*pb.ConfigResponse, error) {
	if backend.updateStrategy == Default {
		if err := backend.syncRemote(resource); err != nil {
			return nil, fmt.Errorf("fail to build config resp: %w", err)
		}
	}

	backend.mu.Lock()
	resp := backend.resp
	backend.mu.Unlock()

	return resp, nil
}

func (backend *Backend) syncRemote(resource *res.Resource) error {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	var lastKnownFingerprint []byte
	if backend.resp != nil {
		lastKnownFingerprint = backend.resp.Fingerprint
	}

	req := &pb.ConfigRequest{
		Resource:             resource,
		LastKnownFingerprint: lastKnownFingerprint,
	}

	resp, err := backend.client.GetConfig(context.Background(), req)
	if err != nil {
		return err
	}

	backend.resp = resp
	return nil
}

func (backend *Backend) Close() error {
	if err := backend.conn.Close(); err != nil {
		return fmt.Errorf("remote config backend fail to close connection: %w", err)
	}

	return nil
}
