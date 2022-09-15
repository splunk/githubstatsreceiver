// The factory file is responsible for providing the required ReceiverFactory object that every
// receiver must return to the opentelemetry receiever.

package githubstatsreceiver

import (
	"context"
	"time"

	"github.com/splunk/githubstatsreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

// define default values as constants here
const (
    typeStr         = "githubmetrics"
    // actual default
    //defaultInterval = 60 * time.Minute
    defaultInterval = 5 * time.Second
    defaultEndpoint = "https://api.github.com"
    defaultTimeout  = 10 * time.Second
)

// create your default config. This should cover all necessary configuration pieces that are not
// verified in the config.go files
func createDefaultConfig() config.Receiver {
    return &Config{
        HTTPClientSettings: confighttp.HTTPClientSettings{
            Endpoint: defaultEndpoint, 
            Timeout: defaultTimeout,
        },
        ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
            ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
            CollectionInterval: defaultInterval,
        },
        Metrics: metadata.DefaultMetricsSettings(),
    }
}

func NewFactory() component.ReceiverFactory {
    return component.NewReceiverFactory(
        typeStr,
        createDefaultConfig,
        component.WithMetricsReceiver(createMetricsReceiver, component.StabilityLevelAlpha))
}

// pull all the pieces together into a factory that the collector can call!
// using the scraperhelper interface greatly simplifies how much code we need to write
// to actually control the metrics scraping.
func createMetricsReceiver(_ context.Context,
    params component.ReceiverCreateSettings,
    baseCfg config.Receiver,
    consumer consumer.Metrics) (component.MetricsReceiver, error) {
    // try and cast baseCfg into a Config object
    cfg := baseCfg.(*Config)
    gitScraper := newGithubMetricsScraper(params, cfg)
    scraper, err := scraperhelper.NewScraper(typeStr, gitScraper.scrape, scraperhelper.WithStart(gitScraper.start))
    if err != nil {
        return nil, err
    }

    return scraperhelper.NewScraperControllerReceiver(
        &cfg.ScraperControllerSettings,
        params,
        consumer,
        scraperhelper.AddScraper(scraper),
        )
}
