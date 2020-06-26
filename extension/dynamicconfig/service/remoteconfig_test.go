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
	"testing"

	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

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

}
