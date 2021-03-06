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
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strconv"
)

// A CollectionPeriod represents the period with which metrics should be
// collected. For optimization purposes, it is strongly recommended to specify
// it as one of the following strings:
//    * "SEC_1":
//    * "SEC_5":
//    * "SEC_10":
//    * "SEC_30":
//    * "MIN_1":
//    * "MIN_5":
//    * "MIN_10":
//    * "MIN_30":
//    * "HR_1":
//    * "HR_2":
//    * "HR_4":
//    * "HR_12":
//    * "DAY_1":
//    * "DAY_7":
//
// However if you have a compelling reason to use a period not included above,
// you may also specify the period value in seconds, as a quoted integer value.
// For example, CollectionPeriod = "60" is equivalent to CollectionPeriod = "MIN_1"
type CollectionPeriod string

var PeriodToInt32 = map[CollectionPeriod]int32{
	"SEC_1":  1,
	"SEC_5":  5,
	"SEC_10": 10,
	"SEC_30": 30,
	"MIN_1":  60,
	"MIN_5":  300,
	"MIN_10": 600,
	"MIN_30": 1800,
	"HR_1":   3600,
	"HR_2":   7200,
	"HR_4":   14400,
	"HR_12":  43200,
	"DAY_1":  86400,
	"DAY_7":  604800,
}

// Proto converts the CollectionPeriod into an int32, for use in the protobuf
// message. It will return an error if the CollectionPeriod has been initialized
// to an unusable value (e.g. arbitrary strings, negative values)
func (period CollectionPeriod) Proto() (int32, error) {
	time, ok := PeriodToInt32[period]
	if !ok {
		return parsePeriod(string(period))
	}

	return time, nil
}

func parsePeriod(periodStr string) (int32, error) {
	value, err := strconv.ParseInt(periodStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("fail to parse period: %v", err)
	}

	if value < 0 {
		return 0, fmt.Errorf("cannot process negative period: %v", value)
	}

	return int32(value), nil
}

// Hash generates an FNVa 64 bit hash of the int32 (little endian) value of
// the CollectionPeriod.
func (period CollectionPeriod) Hash() []byte {
	hasher := fnv.New64a()
	periodSec, _ := period.Proto()

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(periodSec))

	hasher.Write(bs)
	return hasher.Sum(nil)
}
