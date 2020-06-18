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
	"fmt"
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
	viper        *viper.Viper
	MetricConfig *model.MetricConfig
	fingerprint  []byte
	waitTime     int32

	sync.Mutex
}

func NewLocalConfigBackend(configFile string) (*LocalConfigBackend, error) {
	backend := &LocalConfigBackend{
		viper:    config.NewViper(),
		waitTime: 30, // TODO: need more refined strategy for setting this
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
		backend.updateConfig()
	})

	return backend, nil
}

func (backend *LocalConfigBackend) updateConfig() error {
	var metricConfig model.MetricConfig
	if err := backend.viper.UnmarshalExact(&metricConfig); err != nil {
		return fmt.Errorf("local backend failed to decode config: %w", err)
	}

	backend.Lock()
	defer backend.Unlock()

	backend.MetricConfig = &metricConfig
	backend.fingerprint = hashConfig(&metricConfig)

	return nil
}

func hashConfig(obj *model.MetricConfig) []byte {
	return obj.Hash()
}

func (backend *LocalConfigBackend) GetFingerprint() []byte {
	backend.Lock()
	defer backend.Unlock()

	fingerprint := make([]byte, len(backend.fingerprint))
	copy(fingerprint, backend.fingerprint)
	return fingerprint
}

func (backend *LocalConfigBackend) IsSameFingerprint(fingerprint []byte) bool {
	backend.Lock()
	defer backend.Unlock()

	if len(fingerprint) == 0 {
		return false
	}

	return bytes.Equal(backend.fingerprint, fingerprint)
}

func (backend *LocalConfigBackend) BuildConfigResponse() *pb.ConfigResponse {
	backend.Lock()
	defer backend.Unlock()

	return &pb.ConfigResponse{
		Fingerprint:          backend.fingerprint,
		MetricConfig:         backend.MetricConfig.Proto(),
		SuggestedWaitTimeSec: backend.waitTime,
	}
}
