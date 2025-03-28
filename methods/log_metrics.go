package methods

import (
	"github.com/prometheus/client_golang/prometheus"
)

func GetdatetargetlogMetric(RequestLogFilter *prometheus.GaugeVec, queryurl, queryname string, lognum float64) error {
	if lognum != 0 {
		RequestLogFilter.WithLabelValues(queryname, queryurl).Set(lognum)
	} else {
		RequestLogFilter.WithLabelValues(queryname, queryurl).Set(0)
	}
	return nil
}
