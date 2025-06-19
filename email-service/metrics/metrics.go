package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EmailsSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_emails_sent_total",
		Help: "Total number of emails sent",
	}, []string{"sender"})

	EmailErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_email_errors_total",
		Help: "Total number of email sending errors",
	}, []string{"sender"})

	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "sender", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
	}, []string{"method", "sender"})
)
