// Responsible for defining, receiving and validating the configuration of the receiver.

package githubstatsreceiver

import (
	"fmt"

	"github.com/splunk/githubstatsreceiver/internal/metadata"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

// scraperhelper and confighttp are provided by opentelemetry sdk to help extend
// settings in a more standardized way. also saves us the trouble of having to write that
// boilerplate ourselves.
type Config struct {
    scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
    confighttp.HTTPClientSettings           `mapstructure:",squash"`
    Metrics      metadata.MetricsSettings   `mapstructure:"metrics"`
    APIKey       string                     `mapstructure:"api_key"`
    RepoName     []string                   `mapstructure:"repo_name"`
    GitUsername  string                     `mapstructure:"git_username"`
}

func (cfg *Config) Validate() error {
    if (len(cfg.RepoName) == 0) {
        return fmt.Errorf("You must provide at least one valid repository name")
    }

    if (cfg.GitUsername == "") {
        return fmt.Errorf("You must provide a valid github username")
    }

    if (cfg.APIKey == "") {
        return fmt.Errorf("A valid API key is required for the snowflake receiver")
    }

    //if (cfg.CollectionInterval.Minutes() < 60) {
    //    return fmt.Errorf("Interval must be set to at least 1 hour")
    //}
    
    return nil
}
