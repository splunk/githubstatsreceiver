package githubstatsreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestScraper(t *testing.T) {
  s := newGithubMetricsScraper(componenttest.NewNopReceiverCreateSettings(), createDefaultConfig().(*Config))
  err := s.start(context.Background(), componenttest.NewNopHost())
  assert.NoError(t, err)
}
