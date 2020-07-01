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

package remote_test

import (
	"bytes"
	"net"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service/mock"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/service/remote"
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
	"google.golang.org/grpc"
)

// startServer is a test utility to start a quick-n-dirty gRPC server using the
// mock backend.
func StartServer(t *testing.T, quit <-chan struct{}, done chan<- struct{}) string {
	listen, err := net.Listen("tcp", ":0")
	address := listen.Addr()

	if listen == nil || err != nil {
		t.Fatalf("fail to listen: %v", err)
	}

	server := grpc.NewServer()
	configService, _ := service.NewConfigService(service.WithMockBackend())
	pb.RegisterDynamicConfigServer(server, configService)

	go func() {
		done <- struct{}{}
		if err := server.Serve(listen); err != nil {
			t.Errorf("fail to serve: %v", err)
		}
	}()

	go func() {
		<-quit
		configService.Stop()
		server.Stop()

		done <- struct{}{}
	}()

	return address.String()
}

func SetUpServer(t *testing.T) (*remote.Backend, chan struct{}, chan struct{}) {
	quit := make(chan struct{})
	done := make(chan struct{})

	// making mock third-party
	address := StartServer(t, quit, done)
	<-done

	// making remote backend
	backend, err := remote.NewBackend(address)
	if err != nil {
		t.Fatalf("fail to init remote config backend")
	}

	return backend, quit, done
}

func TearDownServer(t *testing.T, backend *remote.Backend, quit chan struct{}, done chan struct{}) {
	quit <- struct{}{}
	if err := backend.Close(); err != nil {
		t.Errorf("fail to close backend: %v", err)
	}

	<-done
}

func TestNewBackend(t *testing.T) {
	if _, err := remote.NewBackend(":0"); err != nil {
		t.Errorf("fail to initialize a new remote backend: %v", err)
	}
}

// TODO: set nondefault update strategy
func TestUpdateStrategy(t *testing.T) {
	backend, err := remote.NewBackend("")
	if err != nil {
		t.Fatalf("fail to init remote config backend: %v", err)
	}

	if strategy := backend.GetUpdateStrategy(); strategy != remote.Default {
		t.Errorf("expected strategy Default, got %v", strategy)
	}

	backend.SetUpdateStrategy(remote.OnGetFingerprint)

	if strategy := backend.GetUpdateStrategy(); strategy != remote.OnGetFingerprint {
		t.Errorf("expected strategy OnGetFingerprint, got %v", strategy)
	}

	const NonsenseStrategy remote.UpdateStrategy = 255
	backend.SetUpdateStrategy(NonsenseStrategy)

	if strategy := backend.GetUpdateStrategy(); strategy != remote.OnGetFingerprint {
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

	if !bytes.Equal(fingerprint, mock.GlobalFingerprint) {
		t.Errorf("expected fingerprint %v, got %v", mock.GlobalFingerprint, fingerprint)
	}
}

func TestBuildConfigResponseRemote(t *testing.T) {
	backend, quit, done := SetUpServer(t)
	defer TearDownServer(t, backend, quit, done)

	resp := buildResp(t, backend)
	if !bytes.Equal(resp.Fingerprint, mock.GlobalResponse.Fingerprint) {
		t.Errorf("expected resp %v, got %v", mock.GlobalResponse, resp)
	}

	newFingerprint := []byte("actually, I believe Gretchen was a cow")
	mock.AlterFingerprint(newFingerprint)

	backend.SetUpdateStrategy(remote.OnGetFingerprint)

	resp = buildResp(t, backend)
	if bytes.Equal(resp.Fingerprint, mock.GlobalResponse.Fingerprint) {
		t.Errorf("expected resp and mock fingerprints to be different, both: %v", resp)
	}

	backend.GetFingerprint(nil)

	resp = buildResp(t, backend)
	if !bytes.Equal(resp.Fingerprint, mock.GlobalResponse.Fingerprint) {
		t.Errorf("expected resp %v, got %v", mock.GlobalResponse, resp)
	}
}

func buildResp(t *testing.T, backend *remote.Backend) *pb.ConfigResponse {
	resp, err := backend.BuildConfigResponse(nil)
	if err != nil {
		t.Errorf("fail to build config response: %v", err)
	}

	return resp
}
