# githubstatsreceiver
Scrape information on github repositories in your organization and emit them as metrics.

## Metrics

These are the metrics available for this scraper.

| Name | Description | Unit | Type | Attributes |
| ---- | ----------- | ---- | ---- | ---------- |
| **github.codechanges.additions** | Code additions to repo | 1 | Gauge(Int) | <ul> </ul> |
| **github.codechanges.deletions** | Code deletions to repo | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.friday** | Number of commits on Friday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.monday** | Number of commits on Monday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.saturday** | Number of commits on Saturday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.sunday** | Number of commits on Sunday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.thursday** | Number of commits on Thursday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.total.weekly** | Number of total commits this week (beginning Sunday) | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.tuesday** | Number of commits on Tuesday | 1 | Gauge(Int) | <ul> </ul> |
| **github.commits.wednesday** | Number of commits on Wednesday | 1 | Gauge(Int) | <ul> </ul> |

**Highlighted metrics** are emitted by default. Other metrics are optional and not emitted by default.
Any metric can be enabled or disabled with the following scraper configuration:

```yaml
metrics:
  <metric_name>:
    enabled: <true|false>
```

### Resource attributes

| Name | Description | Type |
| ---- | ----------- | ---- |
| github.repo.name | The name of the repo being scraped | String |
| github.repo.username | Repository owners username | Empty |

### Metric attributes

| Name | Description | Values |
| ---- | ----------- | ------ |

## Example Config
Like other receivers, githubstatsreceiver is enabled by adding it to the pipeline. The example config provided below shows the most common fields to configure.
```
receivers:
  githubstatsreceiver:
    api_key: <api key>
    repo_name:
    - <repo 1>
    - <repo 2>
    git_username: <username associated with repositories>
    metrics:
      metric_1:
        enabled: <true|false>
      metric_2:
        enabled: <true|false>
```
