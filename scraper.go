// This file contains all of the logic for setting up/initializing the scraper which will pull
// from the REST API. The scraper object will wrap the client and handle the request dispatching,
// converting the json response to a metric, and emitting this metric.

package githubstatsreceiver

import (
	"context"
	"time"

	"github.com/splunk/githubstatsreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scrapererror"
)

type githubMetricsScraper struct {
    client   githubMetricsClient
    settings component.TelemetrySettings
    conf     *Config
    mb       *metadata.MetricsBuilder
}

func newGithubMetricsScraper(settings component.ReceiverCreateSettings, conf *Config) (*githubMetricsScraper) {
    return &githubMetricsScraper{
        settings: settings.TelemetrySettings,
        conf: conf,
        mb: metadata.NewMetricsBuilder(conf.Metrics, settings.BuildInfo),
    }
}

// helper function for initiating scraper. This is particularly useful when used with the 
// scraperhelper API.
func (s *githubMetricsScraper) start(_ context.Context, h component.Host) (err error) {
    s.client, err = newDefaultClient(s.settings, *s.conf, h)
    s.client.logInfo("client created")
    return
}

// this is pretty cool: the actual scrape function is controlled by the scraperhelper in the factory 
// and it aggregates all of the different scrapemethods and their respective errors/emitted metrics.
func (s *githubMetricsScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
    errs := &scrapererror.ScrapeErrors{}

    now := pcommon.NewTimestampFromTime(time.Now())

    s.scrapeCommitStats(ctx, now, errs)
    s.scrapeRepoChanges(ctx, now, errs)

    return s.mb.Emit(), errs.Combine()
}

// wrapper functions for each different endpoint. Also contains the operations to insert the values into the metrics
func (s *githubMetricsScraper) scrapeRepoChanges(ctx context.Context, t pcommon.Timestamp, errs *scrapererror.ScrapeErrors) {
    repoChanges, err := s.client.getRepoChanges(ctx, *s.conf)
    if err != nil {
        errs.Add(err)
        return
    }

    s.mb.RecordCommitsSundayDataPoint(t, int64(repoChanges.Days[0].(float64)))
    s.mb.RecordCommitsMondayDataPoint(t, int64(repoChanges.Days[1].(float64)))
    s.mb.RecordCommitsTuesdayDataPoint(t, int64(repoChanges.Days[2].(float64)))
    s.mb.RecordCommitsWednesdayDataPoint(t, int64(repoChanges.Days[3].(float64)))
    s.mb.RecordCommitsThursdayDataPoint(t, int64(repoChanges.Days[4].(float64)))
    s.mb.RecordCommitsFridayDataPoint(t, int64(repoChanges.Days[5].(float64)))
    s.mb.RecordCommitsSaturdayDataPoint(t, int64(repoChanges.Days[6].(float64)))

    s.mb.RecordCommitsTotalWeeklyDataPoint(t, repoChanges.TotalCommits)
    s.mb.RecordCommitsTotalYtdDataPoint(t, repoChanges.TotalCommits)

    return
}

func (s *githubMetricsScraper) scrapeCommitStats(ctx context.Context, t pcommon.Timestamp, errs *scrapererror.ScrapeErrors) {
    commitStats, err := s.client.getCommitStats(ctx, *s.conf)
    if err != nil {
        errs.Add(err)
        return
    }

    s.mb.RecordCodechangesAdditionsDataPoint(t, int64(commitStats.Insertions))
    s.mb.RecordCodechangesDeletionsDataPoint(t, int64(commitStats.Deletions))

    return
}
