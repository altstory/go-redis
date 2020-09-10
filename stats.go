package redis

import (
	"context"
	"sync"

	"github.com/altstory/go-metrics"
	"github.com/altstory/go-runner"
)

const redisStatsKey = "redis"

var redisMetrics struct {
	Call, Failures *metrics.Metric
}

var metricsOnce sync.Once

func initMetrics() {
	metricsOnce.Do(func() {
		redisMetrics.Call = metrics.Define(&metrics.Def{
			Category: "redis_call",
			Method:   metrics.Sum,
		})
		redisMetrics.Failures = metrics.Define(&metrics.Def{
			Category: "redis_failure",
			Method:   metrics.Sum,
		})
	})
}

func statsForCall(ctx context.Context, err error) {
	runner.StatsFromContext(ctx).Add(redisStatsKey, 1)
	redisMetrics.Call.Add(1)

	if err != nil {
		redisMetrics.Failures.Add(1)
	}
}
