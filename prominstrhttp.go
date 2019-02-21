package prominstrhttp

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HandlerFunc instruments http.HandlerFunc with total requests, duration and request & response sizes.
func HandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return InstrumentHandler(h).ServeHTTP
}

// Handler instruments http.Handler with total requests, duration and request & response sizes.
func Handler(h http.Handler) http.Handler {
	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests_total",
		Help: "Total number of in flight HTTP requests.",
	})
	if err := prometheus.Register(inFlightGauge); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			inFlightGauge = are.ExistingCollector.(prometheus.Gauge)
		} else {
			panic(err)
		}
	}

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"code", "method"},
	)
	if err := prometheus.Register(counter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			counter = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			panic(err)
		}
	}

	reqDur := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "http_request_duration_seconds request duration histogram.",
		Buckets: prometheus.DefBuckets,
	}, []string{"code", "method"})

	if err := prometheus.Register(reqDur); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			reqDur = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			panic(err)
		}
	}

	reqSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request sizes histogram in bytes.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"code", "method"},
	)

	if err := prometheus.Register(reqSize); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			reqSize = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			panic(err)
		}
	}

	respSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response sizes histogram in bytes.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"code", "method"},
	)

	if err := prometheus.Register(respSize); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			respSize = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			panic(err)
		}
	}

	chain :=
		promhttp.InstrumentHandlerResponseSize(respSize,
			promhttp.InstrumentHandlerRequestSize(reqSize,
				promhttp.InstrumentHandlerInFlight(inFlightGauge,
					promhttp.InstrumentHandlerCounter(counter,
						promhttp.InstrumentHandlerDuration(reqDur, h),
					),
				),
			),
		)

	return chain
}
