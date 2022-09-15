package githubstatsreceiver

import (
	"testing"
	//"time"
    "fmt"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
  f := NewFactory()
  assert.Equal(t, "githubmetrics", fmt.Sprintf("%v", f.Type()))

  conf := f.CreateDefaultConfig()
  assert.NotNil(t, conf)

  cfg := conf.(*Config)
  //defaultDur, _ := time.ParseDuration("3600s")
  assert.Equal(t, "https://api.github.com", cfg.Endpoint)
  //assert.Equal(t, defaultDur, cfg.CollectionInterval)
}
