// test code for client.go

package githubstatsreceiver

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
)

func TestClientCreation(t *testing.T) {
    _, err := newDefaultClient(componenttest.NewNopTelemetrySettings(), Config{
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: "http://api.github.com",
		},
	}, componenttest.NewNopHost())
    assert.Equal(t, err, nil)
}

func TestGetRepoChanges(t *testing.T) {
    payload, _ := os.ReadFile("testdata/commit_activity_test_data.json")

    commstats := []commitStats{}
    err := json.Unmarshal(payload, &commstats)

    assert.Equal(t, err, nil)
    assert.Equal(t, int64(1662249600), commstats[len(commstats)-1].WeekStamp)
    assert.Equal(t, int64(0), int64(commstats[len(commstats)-1].Days[0].(float64)))
    fmt.Printf("%T\n", commstats[len(commstats)-1].Days[0])
}

func TestGetCommitStats(t *testing.T) {
    payload, _ := os.ReadFile("testdata/code_frequency_test_data.json")

    comAct, err := newCommitActivity(payload)

    assert.Equal(t, err, nil)
    assert.Equal(t, int64(1662249600), comAct.WeekStamp)
    assert.Equal(t, int64(-22), comAct.Deletions)
    assert.Equal(t, int64(54), comAct.Insertions)
}
