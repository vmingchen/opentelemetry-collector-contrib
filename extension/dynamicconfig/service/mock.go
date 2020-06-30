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
	"net"
	"testing"

	res "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
	"google.golang.org/grpc"
)

var mockFingerprint = []byte("There once was a cat named Gretchen")
var mockResponse = &pb.ConfigResponse{
	Fingerprint: mockFingerprint,
}

func alterFingerprint(newFingerprint []byte) {
	mockFingerprint = newFingerprint
	mockResponse.Fingerprint = mockFingerprint
}

type mockBackend struct{}

func (mock *mockBackend) GetFingerprint(_ *res.Resource) ([]byte, error) {
	return []byte(mockFingerprint), nil
}

func (mock *mockBackend) BuildConfigResponse(_ *res.Resource) (*pb.ConfigResponse, error) {
	return mockResponse, nil
}

func (mock *mockBackend) Close() error {
	return nil
}

func withMockConfig() Option {
	return func(builder *serviceBuilder) {
		builder.backend = &mockBackend{}
	}
}

// startMockServer is a test utility to start a quick-n-dirty gRPC server.
func startMockServer(t *testing.T, configService *ConfigService,
	quit <-chan struct{}, done chan<- struct{}) string {

	listen, err := net.Listen("tcp", ":0")
	address := listen.Addr()

	if listen == nil || err != nil {
		t.Fatalf("fail to listen: %v", err)
	}

	server := grpc.NewServer()
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
