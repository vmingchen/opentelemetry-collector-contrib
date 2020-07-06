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

// Package env provides a detector that loads resource information from
// the OTEL_RESOURCE environment variable. A list of labels of the form
// `<key1>=<value1>,<key2>=<value2>,...` is accepted. Domain names and
// paths are accepted as label keys.
package env

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/consumer/pdata"
)

const TypeStr = "env"

// Environment variable used by "env" to decode a resource.
const envVar = "OTEL_RESOURCE"

type Detector struct{}

func (d *Detector) Detect(context.Context) (pdata.Resource, error) {
	res := pdata.NewResource()
	res.InitEmpty()

	labels := strings.TrimSpace(os.Getenv(envVar))
	if labels == "" {
		return res, nil
	}

	err := initializeAttributeMap(res.Attributes(), labels)
	if err != nil {
		res.Attributes().InitEmptyWithCapacity(0)
		return res, err
	}

	return res, nil
}

// labelRegex matches any key=value pair including a trailing comma or the end of the
// string. Captures the trimmed key & value parts, and ignores any superfluous spaces.
var labelRegex = regexp.MustCompile(`\s*([[:ascii:]]{1,256}?)\s*=\s*([[:ascii:]]{0,256}?)\s*(?:,|$)`)

func initializeAttributeMap(am pdata.AttributeMap, s string) error {
	matches := labelRegex.FindAllStringSubmatchIndex(s, -1)

	for len(matches) == 0 {
		return fmt.Errorf("invalid resource format: %q", s)
	}

	prevIndex := 0
	for _, match := range matches {
		// if there is any text between matches, raise an error
		if prevIndex != match[0] {
			return fmt.Errorf("invalid resource format, invalid text: %q", s[prevIndex:match[0]])
		}

		key := s[match[2]:match[3]]
		value := s[match[4]:match[5]]

		var err error
		if value, err = url.QueryUnescape(value); err != nil {
			return fmt.Errorf("invalid resource format in attribute: %q, err: %s", s[match[0]:match[1]], err)
		}
		am.InsertString(key, value)

		prevIndex = match[1]
	}

	// if there is any text after the last match, raise an error
	if matches[len(matches)-1][1] != len(s) {
		return fmt.Errorf("invalid resource format, invalid text: %q", s[matches[len(matches)-1][1]:])
	}

	return nil
}
