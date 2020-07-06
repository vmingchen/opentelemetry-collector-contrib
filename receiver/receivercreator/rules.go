// Copyright 2020, OpenTelemetry Authors
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

package receivercreator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer"
)

// rule wraps expr rule for later evaluation.
type rule struct {
	program *vm.Program
}

// ruleRe is used to verify the rule starts type check.
var ruleRe = regexp.MustCompile(`^type\.(pod|port)`)

type endpointEnv map[string]interface{}

// endpointToEnv converts an endpoint into a map suitable for expr evaluation.
func endpointToEnv(endpoint observer.Endpoint) (endpointEnv, error) {
	ruleTypes := map[string]interface{}{
		"port": false,
		"pod":  false,
	}

	switch o := endpoint.Details.(type) {
	case observer.Pod:
		ruleTypes["pod"] = true
		return map[string]interface{}{
			"type":        ruleTypes,
			"endpoint":    endpoint.Target,
			"name":        o.Name,
			"labels":      o.Labels,
			"annotations": o.Annotations,
		}, nil
	case observer.Port:
		ruleTypes["port"] = true
		return map[string]interface{}{
			"type":     ruleTypes,
			"endpoint": endpoint.Target,
			"name":     o.Name,
			"port":     o.Port,
			"pod": map[string]interface{}{
				"name":   o.Pod.Name,
				"labels": o.Pod.Labels,
			},
			"protocol": o.Protocol,
		}, nil
	default:
		return nil, fmt.Errorf("unknown endpoint details type %T", endpoint.Details)
	}
}

// newRule creates a new rule instance.
func newRule(ruleStr string) (rule, error) {
	if ruleStr == "" {
		return rule{}, errors.New("rule cannot be empty")
	}
	if !ruleRe.MatchString(ruleStr) {
		// TODO: Try validating against bytecode instead.
		return rule{}, errors.New("rule must specify type")
	}

	// TODO: Maybe use https://godoc.org/github.com/antonmedv/expr#Env in type checking
	// depending on type == specified.
	v, err := expr.Compile(ruleStr)
	if err != nil {
		return rule{}, err
	}
	return rule{v}, nil
}

// eval the rule against the given endpoint.
func (r *rule) eval(env endpointEnv) (bool, error) {
	res, err := expr.Run(r.program, env)
	if err != nil {
		return false, err
	}
	if ret, ok := res.(bool); ok {
		return ret, nil
	}
	return false, errors.New("rule did not return a boolean")
}
