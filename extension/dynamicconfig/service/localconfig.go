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
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/config"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/dynamicconfig/model"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

// LocalConfigBackend is a ConfigBackend that uses a local file to determine
// what schedules to change. The file is read live, so changes to it will
// reflect immediately in the configs.
type LocalConfigBackend struct {
	viper *viper.Viper

	mu           sync.Mutex
	metricConfig *model.MetricConfig
	fingerprint  []byte

	waitTime int32
	updateCh chan struct{} // syncs updates; meant for testing
}

func NewLocalConfigBackend(configFile string) (*LocalConfigBackend, error) {
	backend := &LocalConfigBackend{
		viper:    config.NewViper(),
		waitTime: 30,
	}
	backend.viper.SetConfigFile(configFile)

	if err := backend.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("local backend failed to read config: %w", err)
	}

	if err := backend.updateConfig(); err != nil {
		return nil, err
	}

	backend.viper.WatchConfig()
	backend.viper.OnConfigChange(func(e fsnotify.Event) {
		if err := backend.updateConfig(); err != nil {
			log.Printf("failed to update configs: %v", err)
		}

		if backend.updateCh != nil {
			backend.updateCh <- struct{}{}
		}
	})

	return backend, nil
}

func (backend *LocalConfigBackend) updateConfig() error {
	var config model.MetricConfig
	if err := backend.viper.UnmarshalExact(&config); err != nil {
		return fmt.Errorf("local backend failed to decode config: %w", err)
	}

	backend.mu.Lock()
	defer backend.mu.Unlock()

	backend.metricConfig = &config
	backend.fingerprint = hashConfig(&config)

	return nil
}

func hashConfig(obj *model.MetricConfig) []byte {
	return obj.Hash()
}

func (backend *LocalConfigBackend) GetFingerprint() []byte {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	fingerprint := make([]byte, len(backend.fingerprint))
	copy(fingerprint, backend.fingerprint)
	return fingerprint
}

func (backend *LocalConfigBackend) BuildConfigResponse() *pb.ConfigResponse {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	return &pb.ConfigResponse{
		Fingerprint:          backend.fingerprint,
		MetricConfig:         backend.metricConfig.Proto(),
		SuggestedWaitTimeSec: backend.waitTime,
	}
}
