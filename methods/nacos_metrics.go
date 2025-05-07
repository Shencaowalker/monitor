package methods

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func SetNacosProducernum(NacosProducernumDuration *prometheus.GaugeVec, healthytcount, normalcount, status float64, producer_name string, servicename string) {
	NacosProducernumDuration.WithLabelValues(strconv.Itoa(int(healthytcount)), producer_name, strconv.Itoa(int(normalcount)), servicename).Set(status)

}

func SetNacosServerstatus(NacosServerstatusDuration *prometheus.GaugeVec, currentcount, normalcount, status float64, servicename string) {
	NacosServerstatusDuration.WithLabelValues(strconv.Itoa(int(currentcount)), strconv.Itoa(int(normalcount)), servicename).Set(status)
}

func SetNacosNormalProducersum(NacosNormalProducersumDuration *prometheus.GaugeVec, servicename string, status float64) {
	NacosNormalProducersumDuration.WithLabelValues(servicename).Set(status)
}

func SetNacosCurrentProducersum(NacosCurrentProducersumDuration *prometheus.GaugeVec, servicename string, status float64) {
	NacosCurrentProducersumDuration.WithLabelValues(servicename).Set(status)
}
