package promdb

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestCollector_Collect(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	metricProvider := &mockMetricProvider{
		DBStats: sql.DBStats{
			MaxOpenConnections: rand.Intn(1000),
			OpenConnections:    rand.Intn(1000),
			InUse:              rand.Intn(1000),
			Idle:               rand.Intn(1000),
			WaitCount:          rand.Int63n(1000),
			WaitDuration:       time.Duration(rand.Int63n(1000)),
			MaxIdleClosed:      rand.Int63n(1000),
			MaxIdleTimeClosed:  rand.Int63n(1000),
			MaxLifetimeClosed:  rand.Int63n(1000),
		},
	}
	collector := NewCollector("my_db", metricProvider)

	t.Run("Collect sends all registered metrics", func(t *testing.T) {
		// GIVEN
		metricsReceiver := make(chan prometheus.Metric, math.MaxInt32)

		// WHEN
		collector.Collect(metricsReceiver)
		close(metricsReceiver)

		// THEN
		var collectedMetrics []prometheus.Metric
		for m := range metricsReceiver {
			collectedMetrics = append(collectedMetrics, m)
		}
		assert.Len(t, collectedMetrics, 9)
	})

	t.Run("Collect sends all registered metrics description", func(t *testing.T) {
		// GIVEN
		descs := make(chan *prometheus.Desc, math.MaxInt32)

		// WHEN
		collector.Describe(descs)
		close(descs)

		// THEN
		var descriptions []prometheus.Metric
		for m := range descs {
			descriptions = append(descriptions, m)
		}
		assert.Len(t, descriptions, 9)
	})
}

type mockMetricProvider struct {
	sql.DBStats
}

func (m *mockMetricProvider) Stats() sql.DBStats {
	return m.DBStats
}
