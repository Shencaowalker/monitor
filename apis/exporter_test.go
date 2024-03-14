package apis

// 定义指标
// 先来定义了两个指标数据，一个是Guage类型， 一个是Counter类型。分别代表了CPU温度和磁盘失败次数统计，使用上面的定义进行分类。
// var CpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
// 	Name: "cpu_temperature_celsius",
// 	Help: "Current temperature of the CPU",
// })

// 这里还可以注册其他的参数，比如上面的磁盘失败次数统计上，我们可以同时传递一个device设备名称进去，这样我们采集的时候就可以获得多个不同的指标。
// 每个指标对应了一个设备的磁盘失败次数统计。
// var HdFailures = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "hd_errors_total",
// 		Help: "Number of hard-disk errors.",
// 	},
// 	[]string{"device"},
// )

// 注册指标：
// func init() {
// 	prometheus.MustRegister(CpuTemp)
// 	prometheus.MustRegister(HdFailures)
// }

//使用prometheus.MustRegister是将数据直接注册到Default Registry，
// 就像上面的运行的例子一样，这个Default Registry不需要额外的任何代码就可以将指标传递出去。
// 注册后就可以在程序层面上去使用该指标了，这里我们使用之前定义的指标提供的API（Set和With().Inc）去改变指标的数据内容
// func main() {
// 	CpuTemp.Set(65.3)
// 	HdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

// 	// The Handler function provides a default handler to expose metrics
// 	// via an HTTP server. "/metrics" is the usual endpoint for that.
// 	http.Handle("/metrics", promhttp.Handler())
// 	log.Fatal(http.ListenAndServe(":8080", nil))

// }
