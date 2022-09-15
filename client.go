// This is somewhat unique to this type of receiver, the logic to create and dispatch requests
// is defined here.

package githubstatsreceiver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// common errors for the api call. note that the errcache is not a 4xx but a 202,
// it means that there was technically no issue with the request but that the api has not
// completed the process to return the metrics yet and the requester should try again soon
var (
    errCache           = errors.New("status 202, process not cached")
    errUnauthorized    = errors.New("status 403, unauthorized")
    errUnauthenticated = errors.New("status 401, unauthenticated")
)

// default client which will
type defaultGithubMetricsClient struct {
    client     *http.Client
    url        *url.URL
    headers    []header
    logger     *zap.Logger
}

type header struct {
    headerType  string
    value       string
}

type githubMetricsClient interface {
    getRepoChanges(ctx context.Context, c Config) (*commitStats, error)
    getCommitStats(ctx context.Context, c Config) (*commitActivity, error)
    logInfo(msg string)
}

var _ githubMetricsClient = (*defaultGithubMetricsClient)(nil)

// build your githubreceiverclinet! handles each part of creating a Client which will (eventually) be controlled by a scraper
func newDefaultClient(settings component.TelemetrySettings, c Config, h component.Host) (*defaultGithubMetricsClient, error) {
    client, err := c.HTTPClientSettings.ToClient(h, settings)
    if err != nil {
        return nil, err 
    }

    url, err := url.Parse(c.Endpoint)
    if err != nil {
        return nil, err
    }

    // if you're curious about why these are what they are be sure to check out the github metrics documentation
    // at https://docs.github.com/en/rest/metrics/statistics. Sorry if following that link makes you sad for any reason...
    var headers []header

    headers = append(headers,
        header{"Authorization", fmt.Sprintf("Bearer %s", c.APIKey)},
        header{"Accept", "application/vnd.github+json"},
        )

    return &defaultGithubMetricsClient{
        client:   client,
        url:      url,
        headers:  headers,
        logger:   settings.Logger,
    }, nil
}

func (client defaultGithubMetricsClient) logInfo(msg string) {
    client.logger.Info(fmt.Sprintf("%s", msg))
    return
}
// generic method to build and submit a request
func (client defaultGithubMetricsClient) makeRequest(ctx context.Context, p string) ([]byte, error) {
    // build your endpoint
    endpoint, err := client.url.Parse(p)
    if err != nil {
        return nil, err
    }

    // build your request
    req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
    if err != nil {
        return nil, err
    }

    // add your headers
    for i := 0; i < len(client.headers); i++ {
        req.Header.Add(client.headers[i].headerType, client.headers[i].value)
    }

    r, err := client.client.Do(req)
    if err != nil {
        return nil, err
    }

    // Straightforward check your request! we only return values of 200 everything else gets logged and returns an error
    // that we defined at the top. This is really similar to how the Elastisearch receiver handles requests and is a pretty
    // common go pattern when interacting with API's
    if r.StatusCode == 200 {
        return io.ReadAll(r.Body)
    }

    body, err := io.ReadAll(r.Body)
    
    client.logger.Debug(
        "Request to github api failed",
        zap.String("endpoint", endpoint.String()),
        zap.Int("status", r.StatusCode),
        zap.ByteString("body", body),
        zap.NamedError("body_read_error", err),
        )
    switch r.StatusCode {
    case 403:
        return nil, errUnauthorized
    case 401:
        return nil, errUnauthenticated
    case 202:
        // the caching error is not really an error, its the api asking for you to try again in a second
        time.Sleep(1 * time.Second)
        return client.makeRequest(ctx, p)
    default:
        return nil, fmt.Errorf("non 200 status returned: %d. endpoint: %s", r.StatusCode, endpoint.String())
    }
}

// these wrap the generic request function for each endpoint we want to hit.
func (client defaultGithubMetricsClient) getRepoChanges(ctx context.Context, c Config) (*commitStats, error) {
    p := fmt.Sprintf("/repos/%s/%s/stats/commit_activity", c.GitUsername, c.RepoName) 
    body, err := client.makeRequest(ctx, p)
    if err != nil {
        return nil, err
    }
    
    comstats := []commitStats{}

    err = json.Unmarshal(body, &comstats)
    return &comstats[len(comstats)-1], err
}

func (client defaultGithubMetricsClient) getCommitStats(ctx context.Context, c Config) (*commitActivity, error) {
    p := fmt.Sprintf("/repos/%s/%s/stats/code_frequency", c.GitUsername, c.RepoName)
    body, err := client.makeRequest(ctx, p)
    if err != nil {
        return nil, err
    }

    client.logger.Info("hello from getCommitStats")
    comAct, err := newCommitActivity(body)

    return comAct, err
}


