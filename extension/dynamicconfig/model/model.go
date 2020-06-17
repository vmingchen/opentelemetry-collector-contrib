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
    "fmt"
    "hash"
    "hash/fnv"
)

var hasher hash.Hash = fnv.New64a()

type Hashable interface {
    Hash() []byte
}

func combineHash(chunks []Hashable) []byte {
    if len(chunks) == 0 {
        return nil
    }

    totalHash := chunks[0].Hash()
    for _, chunk := range chunks[1:] {
        chunkHash := chunk.Hash()
        if len(totalHash) != len(chunkHash) {
            panic(fmt.Sprintf("length mismatch: len(%v) != len(%v)",
                totalHash, chunkHash))
        }

        for i, val := range chunkHash {
            totalHash[i] ^= val
        }
    }

    return totalHash
}
