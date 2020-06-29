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
	"fmt"

	"google.golang.org/grpc"

	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

type UpdateStrategy uint8

const (
	Default UpdateStrategy = iota
	OnGetFingerprint
)

type RemoteConfigBackend struct {
	target		   string
	conn           *grpc.ClientConn
	client         pb.DynamicConfigClient
	updateStrategy UpdateStrategy
	chs            *responseMonitorChan
}

type responseMonitorChan struct {
	getResp    chan *pb.ConfigResponse
	updateResp chan *pb.ConfigResponse
	quit       chan struct{}
}

func monitorResponse(chs *responseMonitorChan) {
	var resp *pb.ConfigResponse

	for {
		select {
		case chs.getResp <- resp:
		case resp = <-chs.updateResp:
		case <-chs.quit:
			return
		}
	}
}

func NewRemoteConfigBackend(target string) (*RemoteConfigBackend, error) {
	backend := &RemoteConfigBackend{
		target: 		target,
		conn:           nil,
		client:         nil,
		updateStrategy: Default,
		chs: &responseMonitorChan{
			getResp:    make(chan *pb.ConfigResponse),
			updateResp: make(chan *pb.ConfigResponse),
			quit:       make(chan struct{}),
		},
	}

	if err := backend.initConn(); err != nil {
		return nil, err
	}

	go monitorResponse(backend.chs)
	return backend, nil
}

func (backend *RemoteConfigBackend) initConn() error {
	conn, err := grpc.Dial(
		backend.target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("remote config backend fail to connect: %w", err)
	}

	backend.conn = conn
	backend.client = pb.NewDynamicConfigClient(conn)
	return nil
}


func (backend *RemoteConfigBackend) GetUpdateStrategy() UpdateStrategy {
	return backend.updateStrategy
}

func (backend *RemoteConfigBackend) SetUpdateStrategy(strategy UpdateStrategy) {
	if strategy == Default || strategy == OnGetFingerprint {
		backend.updateStrategy = strategy
	}
}

func (backend *RemoteConfigBackend) GetFingerprint(resource *res.Resource) ([]byte, error) {
	if err := backend.syncRemote(resource); err != nil {
		return nil, fmt.Errorf("fail to get fingerprint: %w", err)
	}

	resp := <-backend.chs.getResp
	return resp.Fingerprint, nil
}

func (backend *RemoteConfigBackend) BuildConfigResponse(resource *res.Resource) (*pb.ConfigResponse, error) {
	if backend.updateStrategy == Default {
		if err := backend.syncRemote(resource); err != nil {
			return nil, fmt.Errorf("fail to build config resp: %w", err)
		}
	}

	resp := <-backend.chs.getResp
	return resp, nil
}

func (backend *RemoteConfigBackend) syncRemote(resource *res.Resource) error {
	var lastKnownFingerprint []byte
	if lastResp := <-backend.chs.getResp; lastResp != nil {
		lastKnownFingerprint = lastResp.Fingerprint
	}

	req := &pb.ConfigRequest{
		Resource:             resource,
		LastKnownFingerprint: lastKnownFingerprint,
	}

	resp, err := backend.client.GetConfig(context.Background(), req)
	if err != nil {
		return err
	}

	backend.chs.updateResp <- resp
	return nil
}

func (backend *RemoteConfigBackend) Close() error {
	backend.chs.quit <- struct{}{}

	// TODO: gRPC connection seems to take inordinately long to close
	if err := backend.conn.Close(); err != nil {
		return fmt.Errorf("remote config backend fail to close connection: %w", err)
	}

	return nil
}
