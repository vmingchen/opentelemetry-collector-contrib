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
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testFileBackend(t *testing.T) {
	t.Log("starting file backend test")
	sec1Schedule := `Schedules:
    - Period: SEC_1`
	sec5Schedule := `Schedules:
    - Period: SEC_5`

	schedFile := getSchedulesFile(t)
	writeString(t, schedFile, sec1Schedule)

	t.Log("starting file backend collector")
	backendCmd, stderr := startCollectorWithFileBackend(t)
	defer backendCmd.Process.Kill()

	t.Log("starting sample application")
	appCmd := startSampleApp(t)
	defer appCmd.Process.Kill()

	t.Log("capturing logs for period=SEC_1")
	avgDuration := timeLogs(t, stderr, 10, 10)
	t.Log("avg duration:", avgDuration)
	if !fuzzyEqualDuration(avgDuration, time.Second, 999*time.Millisecond) {
		t.Errorf("expected period=SEC_1, got: %v", avgDuration)
	}

	writeString(t, schedFile, sec5Schedule)

	t.Log("propogating period=SEC_5")
	time.Sleep(6 * time.Second)

	t.Log("capturing logs for period=SEC_5")
	avgDuration = timeLogs(t, stderr, 10, 10)
	t.Log("avg duration:", avgDuration)
	if !fuzzyEqualDuration(avgDuration, 5*time.Second, 999*time.Millisecond) {
		t.Errorf("expected period=SEC_5, got: %v", avgDuration)
	}
}

func getSchedulesFile(t *testing.T) *os.File {
	file, err := os.OpenFile("testdata/schedules.yaml", os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("fail to open schedules.yaml: %v", err)
	}

	return file
}

// TODO: move to test-common package?
func writeString(t *testing.T, file *os.File, text string) {
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("cannot seek to beginning: %v", err)
	}

	if err := file.Truncate(0); err != nil {
		t.Fatalf("cannot truncate: %v", err)
	}

	if _, err := file.WriteString(text); err != nil {
		file.Close()
		t.Fatalf("cannot write schedule: %v", err)
	}
}

func startCollectorWithFileBackend(t *testing.T) (*exec.Cmd, io.ReadCloser) {
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
	return cmd, stderr
}

func startSampleApp(t *testing.T) *exec.Cmd {
	cmd := exec.Command("app/main")

	if err := cmd.Start(); err != nil {
		t.Fatalf("fail to start app: %v", err)
	}

	return cmd
}

func waitForReady(stderr io.ReadCloser, done chan<- struct{}) error {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		nextLine := scanner.Text()

		if strings.Contains(nextLine, "Everything is ready.") {
			done <- struct{}{}
			return nil
		}

		if strings.Contains(nextLine, "Error:") {
			done <- struct{}{}
			return fmt.Errorf("collector fail: %v", nextLine)
		}
	}

	return fmt.Errorf("end of input reached without reading finish")
}

func timeLogs(t *testing.T, stderr io.ReadCloser, numSamples, discard int) time.Duration {
	scanner := bufio.NewScanner(stderr)

	primeLogTimer(scanner, discard)

	var total time.Duration
	var prevTime time.Time
	for i := 0; i < numSamples; i++ {
		scanner.Scan()
		nextLine := scanner.Text()

		if strings.Contains(nextLine, "MetricsExporter") {
			timeStamp := strings.Fields(nextLine)[0]

			logTimeFmt := "2006-01-02T15:04:05.999-0700"
			timeObj, err := time.Parse(logTimeFmt, timeStamp)
			if err != nil {
				t.Errorf("fail to parse time: %v", timeStamp)
			}

			if !prevTime.IsZero() {
				total += timeObj.Sub(prevTime)
			}

			prevTime = timeObj
		}
	}

	return total / time.Duration(numSamples)

}

func primeLogTimer(scanner *bufio.Scanner, discard int) {
	for {
		scanner.Scan()
		if strings.Contains(scanner.Text(), "MetricsExporter") {
			if discard > 0 {
				discard--
			} else {
				return
			}
		}
	}
}

func fuzzyEqualDuration(first, second, tolerance time.Duration) bool {
	difference := float64(first - second)
	return math.Abs(difference) < float64(tolerance)
}
