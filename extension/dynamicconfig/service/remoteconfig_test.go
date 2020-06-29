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
	"testing"

	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

func SetUpServer(t *testing.T) (*RemoteConfigBackend, chan struct{}, chan struct{}) {
    address := ":50052"

    // making remote backend
    configService, err := NewConfigService(WithRemoteConfig(address))
    if err != nil {
        t.Fatalf("fail to init remote config service: %v", err)
    }

    backend := configService.backend.(*RemoteConfigBackend)
    quit := make(chan struct{})
    done := make(chan struct{})

    // making mock third-party
    mockService, _ := NewConfigService(withMockConfig())
    startMockServer(t, mockService, address, quit, done)
    <-done

    return backend, quit, done
}

func TearDownServer(t *testing.T, backend *RemoteConfigBackend, quit chan struct{}, done chan struct{}) {
    quit <- struct{}{}
    if err := backend.Close(); err != nil {
        t.Errorf("fail to close backend: %v", err)
    }

    <-done
}

func TestResponseMonitor(t *testing.T) {
	chs := &responseMonitorChan{
		getResp:    make(chan *pb.ConfigResponse),
		updateResp: make(chan *pb.ConfigResponse),
		quit:       make(chan struct{}),
	}

	go monitorResponse(chs)

	if resp := <-chs.getResp; resp != nil {
		t.Errorf("expected monitored resp to be nil, got: %v", resp)
	}

	chs.updateResp <- &pb.ConfigResponse{}
	if resp := <-chs.getResp; resp == nil {
		t.Errorf("expected empty resp, got nil")
	}

	chs.quit <- struct{}{}
}

func TestNewRemoteConfigBackend(t *testing.T) {
    backend, quit, done := SetUpServer(t)
    defer TearDownServer(t, backend, quit, done)

    if err := backend.initConn(); err != nil {
        t.Fatalf("failed to connect: %v", err)
    }

    if backend.conn == nil || backend.client == nil {
        t.Errorf("connection structs not properly instantiated")
    }
}

func TestUpdateStrategy(t *testing.T) {
    configService, err := NewConfigService(WithRemoteConfig("0.0.0.0:55800"))
    if err != nil {
        t.Fatalf("fail to init remote config service: %v", err)
    }

    backend := configService.backend.(*RemoteConfigBackend)

    if strategy := backend.GetUpdateStrategy(); strategy != Default {
        t.Errorf("expected strategy Default, got %v", strategy)
    }

    backend.SetUpdateStrategy(OnGetFingerprint)

    if strategy := backend.GetUpdateStrategy(); strategy != OnGetFingerprint {
        t.Errorf("expected strategy OnGetFingerprint, got %v", strategy)
    }

    const NonsenseStrategy UpdateStrategy = 255
    backend.SetUpdateStrategy(NonsenseStrategy)

    if strategy := backend.GetUpdateStrategy(); strategy != OnGetFingerprint {
        t.Errorf("expected strategy OnGetFingerprint, got %v", strategy)
    }
}

func TestGetFingerprintRemote(t *testing.T) {
    backend, quit, done := SetUpServer(t)
    defer TearDownServer(t, backend, quit, done)

    fingerprint, err := backend.GetFingerprint(nil)
    if err != nil {
        t.Errorf("fail to get fingerprint: %v", err)
    }

    if !bytes.Equal(fingerprint, mockFingerprint) {
        t.Errorf("expected fingerprint %v, got %v", mockFingerprint, fingerprint)
    }
}
