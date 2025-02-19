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

package k8sobserver

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Extensions[typeStr] = factory
	cfg, err := configtest.LoadConfigAndValidate(path.Join(".", "testdata", "config.yaml"), factories)

	require.Nil(t, err)
	require.NotNil(t, cfg)

	require.Len(t, cfg.Extensions, 2)

	ext0 := cfg.Extensions[config.NewComponentID(typeStr)]
	assert.EqualValues(t, factory.CreateDefaultConfig(), ext0)

	ext1 := cfg.Extensions[config.NewComponentIDWithName(typeStr, "1")]
	assert.EqualValues(t,
		&Config{
			ExtensionSettings: config.NewExtensionSettings(config.NewComponentIDWithName(typeStr, "1")),
			Node:              "node-1",
			APIConfig:         k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeKubeConfig},
		},
		ext1)
}

func TestValidate(t *testing.T) {
	cfg := &Config{
		ExtensionSettings: config.NewExtensionSettings(config.NewComponentIDWithName(typeStr, "1")),
		Node:              "node-1",
		APIConfig:         k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeKubeConfig},
	}

	err := cfg.Validate()
	require.Nil(t, err)

	cfg.APIConfig.AuthType = "invalid"
	err = cfg.Validate()
	require.NotNil(t, err)
}
