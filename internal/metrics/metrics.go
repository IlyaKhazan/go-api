package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_hits_total",
		Help: "Total number of cache hits",
	})

	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_misses_total",
		Help: "Total number of cache misses",
	})

	CacheSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "flights_cache_current_size",
		Help: "Current number of items in the cache",
	})

	InsertTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_insert_total",
		Help: "Total inserts into the cache",
	})

	UpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_update_total",
		Help: "Total updates in the cache",
	})

	DeleteTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_delete_total",
		Help: "Total deletes from the cache",
	})

	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	CacheExpired = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flights_cache_expired_total",
		Help: "Total number of expired cache entries",
	})

	CacheLastCleanup = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "flights_cache_last_cleanup_time_seconds",
		Help: "Unix timestamp of last cache cleanup",
	})
)

func Init() {
	prometheus.MustRegister(CacheHits, CacheMisses, CacheSize, InsertTotal, UpdateTotal, DeleteTotal, HTTPRequests, HTTPDuration, CacheExpired, CacheLastCleanup)
}
