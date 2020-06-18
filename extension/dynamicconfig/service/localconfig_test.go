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
    "io/ioutil"
    "os"
    "testing"
    "time"
)

func TestNewLocalConfig(t *testing.T) {
    if _, err := NewLocalConfigBackend("woot.txt"); err == nil {
        t.Errorf("failed to catch nonexistant config file")
    }

    if _, err := NewLocalConfigBackend("../testdata/schedules_bad.yaml"); err == nil {
        t.Errorf("failed to catch impropoer config file")
    }

    if _, err := NewLocalConfigBackend("../testdata/schedules.yaml"); err != nil {
        t.Errorf("failed to read config file")
    }
}

func TestUpdateConfig(t *testing.T) {
    originalSchedule := `Schedules:
    - Period: MIN_5`
    updatedSchedule := `Schedules:
    - Period: MIN_1`

    tmpfile, err := ioutil.TempFile("", "schedule.*.yaml")
    if err != nil {
        t.Fatalf("cannot open tempfile: %v", err)
    }
    defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.WriteString(originalSchedule); err != nil {
        tmpfile.Close()
        t.Errorf("cannot write schedule: %v", err)
    }

    backend, err := NewLocalConfigBackend(tmpfile.Name())
    if err != nil {
        t.Errorf("fail to create backend: %v", err)
    }
    backend.updateCh = make(chan struct{})

    if backend.MetricConfig.Schedules[0].Period != "MIN_5" {
        t.Errorf("update incorrect: wanted Period=MIN_5, got MetricConfig: %v",
            backend.MetricConfig)
    }

    if _, err := tmpfile.Seek(0, 0); err != nil {
        t.Fatalf("cannot seek to beginning: %v", err)
    }

    if _, err := tmpfile.WriteString(updatedSchedule); err != nil {
        tmpfile.Close()
        t.Errorf("cannot write schedule: %v", err)
    }

    timeout := make(chan struct{}, 1)
    go func() {
        time.Sleep(5 * time.Second)
        timeout <- struct{}{}
    }()

    select {
    case <-backend.updateCh:
        if backend.MetricConfig.Schedules[0].Period != "MIN_1" {
            t.Errorf("update incorrect: wanted Period=MIN_1, got MetricConfig: %v",
                backend.MetricConfig)
        }
    case <-timeout:
        t.Errorf("local config update timed out")
    }

}
