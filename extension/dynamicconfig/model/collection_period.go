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
//
// Contains common models for the dynamic config service. The corresponding
// Proto() methods convert the model representation to a usable struct for
// protobuf marshalling.

package model

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strconv"
)

type CollectionPeriod string

// TODO: check for library to parse time duration
// TODO: consider how to open up the list of recommended periods for extension
func (period CollectionPeriod) Proto() (int32, error) {
	switch period {
	case "SEC_1":
		return 1, nil
	case "SEC_5":
		return 5, nil
	case "SEC_10":
		return 10, nil
	case "SEC_30":
		return 30, nil
	case "MIN_1":
		return 60, nil
	case "MIN_5":
		return 300, nil
	case "MIN_10":
		return 600, nil
	case "MIN_30":
		return 1800, nil
	case "HR_1":
		return 3600, nil
	case "HR_2":
		return 7200, nil
	case "HR_4":
		return 14400, nil
	case "HR_12":
		return 43200, nil
	case "DAY_1":
		return 86400, nil
	case "DAY_7":
		return 604800, nil
	default:
		value, err := strconv.ParseInt(string(period), 10, 32)
		if err != nil {
			return 0, fmt.Errorf("fail to parse period: %v", err)
		}

		if value < 0 {
			return 0, fmt.Errorf("cannot process negative period: %v", value)
		}

		return int32(value), nil
	}
}

func (period CollectionPeriod) Hash() []byte {
	hasher := fnv.New64a()
	periodSec, _ := period.Proto()

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(periodSec))

	hasher.Write(bs)
	return hasher.Sum(nil)
}
