package promdb

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "db"
	subsystem = "connections"
)

// MetricProvider provides sql.DBStats
type MetricProvider interface {
	Stats() sql.DBStats
}

// collector implements the prometheus.Collector interface.
type collector struct {
	sg      MetricProvider
	metrics []dbMetric
}

type dbMetric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
	Get       func(stat sql.DBStats) float64
}

// NewCollector creates a new prometheus.Collector for database metrics.
func NewCollector(dbName string, sg MetricProvider) prometheus.Collector {
	labels := prometheus.Labels{"db_name": dbName}
	return &collector{
		sg: sg,
		metrics: []dbMetric{
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "max_open"),
					"Maximum number of open connections to the database.",
					nil,
					labels,
				),
				ValueType: prometheus.GaugeValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.MaxOpenConnections)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "open"),
					"The number of established connections both in use and idle.",
					nil,
					labels,
				),
				ValueType: prometheus.GaugeValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.OpenConnections)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "in_use"),
					"The number of connections currently in use.",
					nil,
					labels,
				),
				ValueType: prometheus.GaugeValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.InUse)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "idle"),
					"The number of idle connections.",
					nil,
					labels,
				),
				ValueType: prometheus.GaugeValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.Idle)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "wait_count"),
					"The total number of connections waited for.",
					nil,
					labels,
				),
				ValueType: prometheus.CounterValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.WaitCount)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "wait_duration_seconds"),
					"The total time blocked waiting for a new connection in seconds.",
					nil,
					labels,
				),
				ValueType: prometheus.CounterValue,
				Get: func(stat sql.DBStats) float64 {
					return stat.WaitDuration.Seconds()
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "max_idle_closed"),
					"The total number of connections closed due to SetMaxIdleConns.",
					nil,
					labels,
				),
				ValueType: prometheus.CounterValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.MaxIdleClosed)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "max_idle_time_closed"),
					"The total number of connections closed due to SetConnMaxIdleTime.",
					nil,
					labels,
				),
				ValueType: prometheus.CounterValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.MaxIdleTimeClosed)
				},
			},
			{
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "max_lifetime_closed"),
					"The total number of connections closed due to SetConnMaxLifetime.",
					nil,
					labels,
				),
				ValueType: prometheus.CounterValue,
				Get: func(stat sql.DBStats) float64 {
					return float64(stat.MaxLifetimeClosed)
				},
			},
		},
	}
}

// Describe implements the prometheus.Collector interface.
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.Desc
	}
}

// Collect implements the prometheus.Collector interface.
func (c collector) Collect(ch chan<- prometheus.Metric) {
	stats := c.sg.Stats()
	for _, m := range c.metrics {
		ch <- prometheus.MustNewConstMetric(
			m.Desc,
			m.ValueType,
			m.Get(stats),
		)
	}
}
