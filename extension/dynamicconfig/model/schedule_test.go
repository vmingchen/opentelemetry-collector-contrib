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

package model

import (
	"bytes"
	"testing"
)

func TestScheduleProto(t *testing.T) {
	schedule := Schedule{
		InclusionPatterns: []Pattern{Pattern{}, Pattern{}},
		ExclusionPatterns: []Pattern{Pattern{}, Pattern{}},
		Period:            "MIN_5",
	}

	p, err := schedule.Proto()
	if err != nil || len(p.InclusionPatterns) != 2 ||
		len(p.ExclusionPatterns) != 2 ||
		p.PeriodSec != 300 {
		t.Errorf("improper conversion to proto")
	}
}

func TestScheduleHash(t *testing.T) {
	configA := Schedule{
		InclusionPatterns: []Pattern{
			Pattern{Equals: "woot"},
			Pattern{StartsWith: "yay"},
		},
	}

	configB := Schedule{
		InclusionPatterns: []Pattern{
			Pattern{StartsWith: "yay"},
			Pattern{Equals: "woot"},
		},
	}

	configC := Schedule{
		ExclusionPatterns: []Pattern{
			Pattern{Equals: "woot"},
			Pattern{StartsWith: "yay"},
		},
	}

	if !bytes.Equal(configA.Hash(), configB.Hash()) {
		t.Errorf("identical configs with different hashes")
	}

	if bytes.Equal(configA.Hash(), configC.Hash()) {
		t.Errorf("different configs with identical hashes")
	}
}
