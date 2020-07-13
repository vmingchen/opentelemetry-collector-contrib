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
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

// This file implements factory for receiver_creator. A receiver_creator can create other receivers at runtime.

const (
	typeStr = "receiver_creator"
)

// Factory is the factory for receiver_creator.
type Factory struct {
}

var _ component.ReceiverFactoryOld = (*Factory)(nil)

// Type gets the type of the Receiver config created by this factory.
func (f *Factory) Type() configmodels.Type {
	return configmodels.Type(typeStr)
}

// CustomUnmarshaler returns custom unmarshaler for receiver_creator config.
func (f *Factory) CustomUnmarshaler() component.CustomUnmarshaler {
	return func(sourceViperSection *viper.Viper, intoCfg interface{}) error {
		if sourceViperSection == nil {
			// Nothing to do if there is no config given.
			return nil
		}
		c := intoCfg.(*Config)

		if err := sourceViperSection.Unmarshal(&c); err != nil {
			return err
		}

		receiversCfg := viperSub(sourceViperSection, receiversConfigKey)

		for subreceiverKey := range receiversCfg.AllSettings() {
			cfgSection := viperSub(receiversCfg, subreceiverKey).GetStringMap(configKey)
			subreceiver, err := newReceiverTemplate(subreceiverKey, cfgSection)
			if err != nil {
				return err
			}

			// Unmarshals receiver_creator configuration like rule.
			if err = receiversCfg.UnmarshalKey(subreceiverKey, &subreceiver); err != nil {
				return fmt.Errorf("failed to deserialize sub-receiver %q: %s", subreceiverKey, err)
			}

			subreceiver.rule, err = newRule(subreceiver.Rule)
			if err != nil {
				return fmt.Errorf("subreceiver %q rule is invalid: %v", subreceiverKey, err)
			}

			c.receiverTemplates[subreceiverKey] = subreceiver
		}

		return nil
	}
}

// CreateDefaultConfig creates the default configuration for receiver_creator.
func (f *Factory) CreateDefaultConfig() configmodels.Receiver {
	return &Config{
		ReceiverSettings: configmodels.ReceiverSettings{
			TypeVal: configmodels.Type(typeStr),
			NameVal: typeStr,
		},
		receiverTemplates: map[string]receiverTemplate{},
	}
}

// CreateTraceReceiver creates a trace receiver based on provided config.
func (f *Factory) CreateTraceReceiver(context.Context, *zap.Logger, configmodels.Receiver,
	consumer.TraceConsumerOld) (component.TraceReceiver, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}

// CreateMetricsReceiver creates a metrics receiver based on provided config.
func (f *Factory) CreateMetricsReceiver(
	ctx context.Context,
	logger *zap.Logger,
	cfg configmodels.Receiver,
	consumer consumer.MetricsConsumerOld,
) (component.MetricsReceiver, error) {
	return newReceiverCreator(logger, consumer, cfg.(*Config))
}
