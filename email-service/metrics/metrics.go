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
)
