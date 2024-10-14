package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	pingCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "No of request handled by Ping handler",
		},
	)
	topCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "top_request_count",
			Help: "No of request handled by Top handler",
		})
	userRepoCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_repo_request_count",
			Help: "No of request handled by User repo handler",
		})
	userRepoResultGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_repo_result_count",
			Help: "Number of repos returned by request",
		}, []string{"user"})
)

func registerMetrics() {
	prometheus.MustRegister(pingCounter)
	prometheus.MustRegister(topCounter)
	prometheus.MustRegister(userRepoCounter)
	prometheus.MustRegister(userRepoResultGauge)
}
