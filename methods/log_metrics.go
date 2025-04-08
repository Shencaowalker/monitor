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

func GetdatetargetlogMetricnew(RequestLogFilter *prometheus.GaugeVec, queryurl, queryname, traceidwithskywalking, traceidwiithaloudata, responsecode, starttime, apiurl string, cost float64) error {
	if cost != 0 {
		RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwithskywalking, traceidwiithaloudata, responsecode, starttime, apiurl).Set(cost)
	} else {
		RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwithskywalking, traceidwiithaloudata, responsecode, starttime, apiurl).Set(0)
	}
	return nil
}

func GeterrorlogMetricnew(RequestLogFilter *prometheus.GaugeVec, queryurl, queryname, traceidwiithaloudata, responsecode, starttime, apiurl string, cost float64) error {
	if cost != 0 {
		RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwiithaloudata, responsecode, starttime, apiurl).Set(cost)
	} else {
		RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwiithaloudata, responsecode, starttime, apiurl).Set(0)
	}
	return nil
}

func GetdatetargetlogMetricerror(RequestLogFilter *prometheus.GaugeVec, queryname, queryurl, traceidwithskywalking, traceidwiithaloudata, errormsg, starttime string, value float64) error {
	RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwithskywalking, traceidwiithaloudata, errormsg, starttime).Set(value)
	return nil
}

func GetdatetargetlogMetricquerystage(RequestLogFilter *prometheus.GaugeVec, queryurl, queryname, traceidwiithaloudata, responsecode, stages, starttime string, value float64) error {
	RequestLogFilter.WithLabelValues(queryname, queryurl, traceidwiithaloudata, responsecode, stages, starttime).Set(value)
	return nil
}
