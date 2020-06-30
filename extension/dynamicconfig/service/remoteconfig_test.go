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
	quit := make(chan struct{})
	done := make(chan struct{})

	// making mock third-party
	mockService, _ := NewConfigService(withMockConfig())
	address := startMockServer(t, mockService, quit, done)
	<-done

	// making remote backend
	configService, err := NewConfigService(WithRemoteConfig(address))
	if err != nil {
		t.Fatalf("fail to init remote config service: %v", err)
	}

	backend := configService.backend.(*RemoteConfigBackend)
	return backend, quit, done
}

func TearDownServer(t *testing.T, backend *RemoteConfigBackend, quit chan struct{}, done chan struct{}) {
	quit <- struct{}{}
	if err := backend.Close(); err != nil {
		t.Errorf("fail to close backend: %v", err)
	}

	<-done
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

// TODO: set nondefault update strategy
func TestUpdateStrategy(t *testing.T) {
	backend, err := NewRemoteConfigBackend("")
	if err != nil {
		t.Fatalf("fail to init remote config backend: %v", err)
	}

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

func TestBuildConfigResponseRemote(t *testing.T) {
	backend, quit, done := SetUpServer(t)
	defer TearDownServer(t, backend, quit, done)

	resp := buildResp(t, backend)
	if !bytes.Equal(resp.Fingerprint, mockResponse.Fingerprint) {
		t.Errorf("expected resp %v, got %v", mockResponse, resp)
	}

	newFingerprint := []byte("actually, I believe Gretchen was a cow")
	alterFingerprint(newFingerprint)

	backend.SetUpdateStrategy(OnGetFingerprint)

	resp = buildResp(t, backend)
	if bytes.Equal(resp.Fingerprint, mockResponse.Fingerprint) {
		t.Errorf("expected resp and mock fingerprints to be different, both: %v", resp)
	}

	backend.GetFingerprint(nil)

	resp = buildResp(t, backend)
	if !bytes.Equal(resp.Fingerprint, mockResponse.Fingerprint) {
		t.Errorf("expected resp %v, got %v", mockResponse, resp)
	}
}

func buildResp(t *testing.T, backend *RemoteConfigBackend) *pb.ConfigResponse {
	resp, err := backend.BuildConfigResponse(nil)
	if err != nil {
		t.Errorf("fail to build config response: %v", err)
	}

	return resp
}
