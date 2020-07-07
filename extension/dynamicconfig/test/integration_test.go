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

// +build integration

package test

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TODO: double check build target works
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Log("warning: not recompiling otelcontribcol: omit -test.short flag to compile new binary")
	} else {
		t.Log("building new collector")
		buildCollector(t)
	}

	t.Log("starting file backend test")
	cmd := startCollectorWithFileBackend(t)
	defer cmd.Process.Kill()

	t.Log("starting sample application")
	quit := make(chan struct{})
	go runSampleApp(quit)

	time.Sleep(10 * time.Second)
	quit <- struct{}{}

}

func buildCollector(t *testing.T) {
	cmd := exec.Command("make", "otelcontribcol")
	cmd.Dir = "../../../" // run in top-level of repo

	if err := cmd.Run(); err != nil {
		t.Fatalf("fail to compile otelcontribcol: %v", err)
	}
}

func startCollectorWithFileBackend(t *testing.T) *exec.Cmd {
	cmdPath := fmt.Sprintf("../../../bin/otelcontribcol_%s_%s",
		runtime.GOOS,
		runtime.GOARCH)
	cmd := exec.Command(cmdPath, "--config", "testdata/file-backend-config.yaml")

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("fail to redirect stderr: %v", err)
	}

	done := make(chan struct{})
	go func(t *testing.T) {
		if err := waitForReady(stderr, done); err != nil {
			t.Fatalf(err.Error())
		}
	}(t)

	if err := cmd.Start(); err != nil {
		t.Fatalf("fail to start otelcontribcol: %v", err)
	}

	<-done
	return cmd
}

func waitForReady(stderr io.ReadCloser, done chan<- struct{}) error {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		nextLine := scanner.Text()
		fmt.Println("[READ]", nextLine)

		if strings.Contains(nextLine, "Everything is ready.") {
			done <- struct{}{}
			// return nil
		}

		if strings.Contains(nextLine, "Error:") {
			done <- struct{}{}
			return fmt.Errorf("collector fail: %v", nextLine)
		}
	}

	return fmt.Errorf("end of input reached without reading finish")
}
