package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsHttPort = 8889 // for promethes to connect

var (
	ClientConnected prometheus.Gauge
	ReqRecvTotal    prometheus.Counter
	RspSendTotal    prometheus.Counter
)

func init() {
	ReqRecvTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_demo_req_recv_total",
	})

	RspSendTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_demo_resp_send_total",
	})

	ClientConnected = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_server_demo_client_connected",
	})

	prometheus.MustRegister(ReqRecvTotal, RspSendTotal, ClientConnected)

	// start the metrics server
	metricsServer := &http.Server{
		Addr: fmt.Sprintf(":%d", metricsHttPort),
	}

	mu := http.NewServeMux()
	mu.Handle("/metrics", promhttp.Handler())
	metricsServer.Handler = mu

	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			fmt.Println("prometheus-exporter http server start failed:", err.Error())
			return
		}
	}()

	fmt.Println("metrics server start at (*:8889)")
}
