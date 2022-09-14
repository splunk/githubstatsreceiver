package githubstatsreceiver

import (
	"testing"
	"time"
    "path"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestValidConfig(t *testing.T) {
    factories, err := componenttest.NopFactories()
    assert.NoError(t, err)

    factories.Receivers[typeStr] = NewFactory()

    conf, err := servicetest.LoadConfigAndValidate(path.Join("testdata", "test_config.yaml"), factories)
    assert.NoError(t, err)

    cfg := conf.Receivers[config.NewComponentID(typeStr)].(*Config)
    assert.Equal(t, "abc123", cfg.APIKey)
    assert.Equal(t, "truck", cfg.RepoName)
    assert.Equal(t, "buck", cfg.GitUsername)
    duration, _ := time.ParseDuration("3600s")
    assert.Equal(t, duration, cfg.CollectionInterval)
}

