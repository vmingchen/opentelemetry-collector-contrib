# Redis Receiver

### Overview

The Redis receiver is designed to retrieve Redis INFO data from a single Redis
instance, build metrics from that data, and send them to the next consumer at a
configurable interval.

Status: beta

### Details

The Redis INFO command returns information and statistics about a Redis
server (see [https://redis.io/commands/info](https://redis.io/commands/info) for
details). The Redis receiver extracts values from the result and converts them to open
telemetry metrics. Details about the metrics produced by the Redis receiver
can be found by browsing [metric_functions.go](metric_functions.go).

For example, one of the fields returned by the Redis INFO command is
`used_cpu_sys` which indicates the system CPU consumed by the Redis server,
expressed in seconds, since the start of the Redis instance.

The Redis receiver turns this data into a gauge...

```go
func usedCPUSys() *redisMetric {
	return &redisMetric{
		key:    "used_cpu_sys",
		name:   "redis/cpu/time",
		units:  "s",
		mdType: metricspb.MetricDescriptor_GAUGE_DOUBLE,
		labels: map[string]string{"state": "sys"},
	}
}
```

with a metric name of "redis/cpu/time" and a units value of "s" (seconds).

# Configuration

Note: this receiver is in beta and configuration fields are subject to change.

Example configuration:

```yaml
receivers:
  redis:
    endpoint: "localhost:6379"
    service_name: "my-test-redis"
    collection_interval: 10s
    password: $REDIS_PASSWORD
```

### endpoint

The hostname and port of the Redis instance, separated by a colon.

_Required._

### collection_interval (default: 10s)

This receiver runs on an interval. Each time it runs, it queries Redis, creates
metrics, and sends them to the next consumer. The `collection_interval`
configuration option tells this receiver the duration between runs.

This value must be a string readable by Golang's `ParseDuration` function:
e.g. "1h30m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

_Required._

### service_name

The logical name of the Redis server. This value will be added as a
`service_name` Resource label and may end up as a dimension on exported
metrics, depending on the exporter.

_Required._

### password

The password used to access the Redis instance; must match the password
specified in the `requirepass` server configuration option.

Note: as with all Open Telemetry configuration values, a reference to an
environment variable is supported. For example, to pick up the value of
an environment variable `REDIS_PASSWORD`, you could use a configuration like
the following:

```yaml
receivers:
  redis:
    endpoint: "localhost:6379"
    service_name: "my-test-redis"
    collection_interval: 10s
    password: $REDIS_PASSWORD
```

_Optional._
