package apis

import (
	"monitor/methods"
	"net/http"
)

// 定义指标
// 先来定义了两个指标数据，一个是Guage类型， 一个是Counter类型。分别代表了CPU温度和磁盘失败次数统计，使用上面的定义进行分类。
// var examCpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
// 	Name: "cpu_temperature_celsius",
// 	Help: "Current temperature of the CPU111111",
// })

// 这里还可以注册其他的参数，比如上面的磁盘失败次数统计上，我们可以同时传递一个device设备名称进去，这样我们采集的时候就可以获得多个不同的指标。
// 每个指标对应了一个设备的磁盘失败次数统计。
// var examHdFailures = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "hd_errors_total1111111111",
// 		Help: "Number of hard-disk errors.11111111111",
// 	},
// 	[]string{"device"},
// )

func Exporterexamaaa() http.Handler {
	// methods.PgQueryMetrics(methods.Viper_reader(), methods.Db_Collect_Postgres(methods.Viper_reader()))
	// prometheus.MustRegister(aaa[0])
	// methods.MustRegisterOnce(methods.Viper_reader(), methods.Db_Collect_Postgres(methods.Viper_reader()))
	// examCpuTemp.Set(rand.Float64())
	// examHdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
	// promhttp.HandlerOpts
	// prometheus.MustRegister(examCpuTemp)
	// prometheus.MustRegister(examHdFailures)
	// go func() error {
	// 	for {
	// 		err := methods.PgQueryAndUpdateMetrics(methods.Viper_reader(), methods.Db_Collect_Postgres(methods.Viper_reader()))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		time.Sleep(1 * time.Minute)
	// 	}
	// }()
	return methods.MustRegisterOnce(methods.Viper_reader(), methods.Db_Collect_Postgres(methods.Viper_reader()))

}
