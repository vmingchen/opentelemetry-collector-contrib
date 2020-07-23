// Copyright 2019, OpenTelemetry Authors
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

package signalfxexporter

import (
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/translation"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/splunk"
)

const (
	// The value of "type" key in configuration.
	typeStr = "signalfx"

	defaultHTTPTimeout = time.Second * 5
)

// Factory is the factory for SignalFx exporter.
type Factory struct {
}

// Type gets the type of the Exporter config created by this factory.
func (f *Factory) Type() configmodels.Type {
	return configmodels.Type(typeStr)
}

// CreateDefaultConfig creates the default configuration for exporter.
func (f *Factory) CreateDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: configmodels.Type(typeStr),
			NameVal: typeStr,
		},
		Timeout: defaultHTTPTimeout,
		AccessTokenPassthroughConfig: splunk.AccessTokenPassthroughConfig{
			AccessTokenPassthrough: true,
		},
		SendCompatibleMetrics: false,
		TranslationRules:      nil,
	}
}

// CreateTraceExporter creates a trace exporter based on this config.
func (f *Factory) CreateTraceExporter(
	logger *zap.Logger,
	config configmodels.Exporter,
) (component.TraceExporterOld, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}

// CreateMetricsExporter creates a metrics exporter based on this config.
func (f *Factory) CreateMetricsExporter(
	logger *zap.Logger,
	config configmodels.Exporter,
) (exp component.MetricsExporterOld, err error) {

	expCfg := config.(*Config)
	if expCfg.SendCompatibleMetrics && expCfg.TranslationRules == nil {
		expCfg.TranslationRules, err = loadDefaultTranslationRules()
		if err != nil {
			return nil, err
		}
	}

	exp, err = New(expCfg, logger)

	if err != nil {
		return nil, err
	}

	return exp, nil
}

func loadDefaultTranslationRules() ([]translation.Rule, error) {
	config := Config{}

	v := otelconfig.NewViper()
	v.SetConfigType("yaml")
	v.ReadConfig(strings.NewReader(translation.DefaultTranslationRulesYaml))
	err := v.UnmarshalExact(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to load default translation rules: %v", err)
	}

	return config.TranslationRules, nil
}
