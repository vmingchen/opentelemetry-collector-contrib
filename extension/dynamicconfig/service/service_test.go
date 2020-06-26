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
	"testing"

	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

func TestNewConfigService(t *testing.T) {
	if service, err := NewConfigService(); service != nil || err == nil {
		t.Errorf("no backend specified but service created: %v: %v", service, err)
	}

	if service, err := NewConfigService(withMockConfig()); service == nil || err != nil {
		t.Errorf("backend specified but service not created: %v: %v", service, err)
	}
}

func TestRemoteConfigOption(t *testing.T) {
	t.Errorf("unimplemented")
}

func TestLocalConfigOption(t *testing.T) {
	if service, err := NewConfigService(WithLocalConfig("woot.yaml")); service != nil || err == nil {
		t.Errorf("file does not exist but service created: %v: %v", service, err)
	}

	service, err := NewConfigService(WithLocalConfig("../testdata/schedules.yaml"))
	if service == nil || err != nil {
		t.Errorf("file exists but service not created: %v: %v", service, err)
	}
}

func TestWaitTimeConfigOption(t *testing.T) {
	const testWaitTime = 60

	service, err := NewConfigService(
		WithLocalConfig("../testdata/schedules.yaml"),
		WithWaitTime(testWaitTime),
	)
	if service == nil || err != nil {
		t.Errorf("file exists but service not created: %v: %v", service, err)
	}

	time := service.backend.(*LocalConfigBackend).waitTime
	if time != testWaitTime {
		t.Errorf("wait time of %d requested, found %d", testWaitTime, time)
	}

}

func TestGetConfig(t *testing.T) {
	service, err := NewConfigService(withMockConfig())
	sameFingerprintReq := pb.ConfigRequest{LastKnownFingerprint: mockFingerprint}

	resp, err := service.GetConfig(context.Background(), &sameFingerprintReq)
	if err != nil {
		t.Errorf("failed to get config: %v", err)
	}

	if !bytes.Equal(resp.Fingerprint, mockFingerprint) {
		t.Errorf("expected fingerprint to equal %v: got %v", mockFingerprint, resp.Fingerprint)
	}

	blankReq := pb.ConfigRequest{}
	resp, err = service.GetConfig(context.Background(), &blankReq)
	if err != nil {
		t.Errorf("failed to get config: %v", err)
	}
}
