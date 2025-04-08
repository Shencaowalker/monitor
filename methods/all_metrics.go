package methods

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

//	func newGaugeVec(name, help string) *prometheus.GaugeVec {
//		return prometheus.NewGaugeVec(
//			prometheus.GaugeOpts{
//				Name: name,
//				Help: help,
//				// Sql_content: sql_content,
//			},
//			[]string{"host"},
//		)
//	}
func MustRegisterOnce() http.Handler {
	return promhttp.Handler()
}

func PostgresqlMustRegisterOnce(config *viper.Viper, db *sql.DB) {
	// defer db.Close()
	fmt.Println("postgres.metrics is ", config.GetStringSlice("postgres.server")[0])
	for _, metricaa := range config.GetStringMap("postgres.metrics.items") {
		name := metricaa.(map[string]interface{})["name"].(string)
		help := metricaa.(map[string]interface{})["help"].(string)
		sql_content := metricaa.(map[string]interface{})["sql_content"].(string)
		fmt.Println("metricaa is ", metricaa)
		// newGaugeVecinstance := newGaugeVec(metricaa.(map[string]interface{})["name"].(string), metricaa.(map[string]interface{})["help"].(string))
		newGaugeVecinstance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
		go func() error {
			for {
				err := PgQueryAndUpdateMetric(newGaugeVecinstance, db, sql_content)
				if err != nil {
					return err
				}

				duration, err := time.ParseDuration(config.GetString("postgres.metrics.cycle"))
				// fmt.Println(config.GetString("postgres.metrics.cycle"))
				if err != nil {
					panic(err)

				}
				time.Sleep(duration)
			}
		}()
		prometheus.MustRegister(newGaugeVecinstance)
	}

}

func ReLogNumMustRegisterOnce(config *viper.Viper) {
	// defer db.Close()
	log.Println("调用接口产出日志")
	// logs_lists := config.GetStringMap("mixedformat.requestlogs")
	// 获取待收集日志指标项目
	logs_lists := config.GetStringSlice("mixedformat.requestlogs")
	// 日志指标项统计区间
	collectionscopeseconds := config.GetString("mixedformat.collectionscopeseconds")
	// 设置延迟时间，主要是为了防止loli日志不能及时被查询，延迟时间间隔进行日志收取
	latencycollectionseconds := config.GetString("mixedformat.latencycollectionseconds")
	// loki地址
	lokiipport := config.GetString("mixedformat.lokiipport")
	RequestLogFilter := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "request_log_num",
			Help: "Log query based on labels",
		},
		[]string{"name", "url"},
	)

	go func() error {
		for {
			for _, i := range logs_lists {
				// log_pro := i
				label_list := config.GetStringMap(i + ".label_list")
				lokire := config.GetStringSlice(i + ".lokire")
				lokiexclre := config.GetStringSlice(i + ".lokiexclre")
				lokire_string := ""
				// 是否进行string模糊匹配
				if len(lokire) != 0 {
					for _, j := range lokire {
						lokire_string += "|~`" + j + "`"
					}
				}
				// 是否进行string 剔除匹配
				if len(lokiexclre) != 0 {
					for _, j := range lokiexclre {
						lokire_string += "!~`" + j + "`"
					}
				}
				// 获取每次查询loki的日志条数限制
				recordslimit := config.GetString(i + ".recordslimit")
				// 进行loki查询
				data, url, _ := Getlogandurlfromloki(lokiipport, label_list, time.Now(), latencycollectionseconds, collectionscopeseconds, recordslimit, lokire_string)

				log.Println(i, "的数据是:", len(data))
				err := GetdatetargetlogMetric(RequestLogFilter, url, i, float64(len(data)))
				if err != nil {
					return err
				}
			}

			duration, err := time.ParseDuration(config.GetString("mixedformat.cycle"))
			// fmt.Println(config.GetString("postgres.metrics.cycle"))
			if err != nil {
				panic(err)
			}
			time.Sleep(duration)
			RequestLogFilter.Reset()
		}
	}()

	prometheus.MustRegister(RequestLogFilter)
}

// 获取request请求值
func ReLogNumMustRegisterOncewithurl(config *viper.Viper) {
	// defer db.Close()
	log.Println("调用接口产出日志")
	// logs_lists := config.GetStringMap("mixedformat.requestlogs")
	// 获取待收集日志指标项目
	logs_lists := config.GetStringSlice("mixedformat.requestlogs")
	// 日志指标项统计区间
	collectionscopeseconds := config.GetString("mixedformat.collectionscopeseconds")
	// 设置延迟时间，主要是为了防止loli日志不能及时被查询，延迟时间间隔进行日志收取
	latencycollectionseconds := config.GetString("mixedformat.latencycollectionseconds")
	// loki地址
	lokiipport := config.GetString("mixedformat.lokiipport")
	RequestLogFilternew := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "request_log_num_with_interface",
			Help: "Log query based on labels",
		},
		[]string{"name", "url", "traceidwithskywalking", "traceidwiithaloudata", "responsecode", "starttime", "apiurl"},
	)

	go func() error {
		for {
			for _, i := range logs_lists {
				// log_pro := i
				label_list := config.GetStringMap(i + ".label_list")
				lokire := config.GetStringSlice(i + ".lokire")
				lokiexclre := config.GetStringSlice(i + ".lokiexclre")
				lokire_string := ""
				// 是否进行string模糊匹配
				if len(lokire) != 0 {
					for _, j := range lokire {
						lokire_string += "|~`" + j + "`"
					}
				}
				// 是否进行string 剔除匹配
				if len(lokiexclre) != 0 {
					for _, j := range lokiexclre {
						lokire_string += "!~`" + j + "`"
					}
				}
				// 获取每次查询loki的日志条数限制
				recordslimit := config.GetString(i + ".recordslimit")
				// 进行loki查询,返回匹配的值列表
				data, url, _ := Getlogandurlfromloki(lokiipport, label_list, time.Now(), latencycollectionseconds, collectionscopeseconds, recordslimit, lokire_string)

				log.Println(i, "的数据是:", len(data))
				for _, jj := range data {
					var aaa map[string]string
					aaa = make(map[string]string)

					// reg, err = regexp.Compile("(\\d+-\\d+-\\d+\\s\\S+)\\s\\[(.*?)\\]\\s\\[.*?\\]\\s(\\w+)\\s+(\\S+)\\s-\\s(.*)")
					regex := config.GetString(i + ".regex")
					reg, err := regexp.Compile(regex)
					if err != nil {
						fmt.Printf("正则表达式编译错误: %v\n", err)
						continue
					}
					// 获取子匹配总数
					submatchCount := reg.NumSubexp()
					if submatchCount == 0 {
						fmt.Println("警告: 正则表达式不包含捕获组")
						continue
					}
					label_name := config.GetString("mixedformat.requestlogslabel_name")
					value_label := config.GetString(i + ".value_label")

					customNames := strings.Split(label_name, " ")
					// 验证名称数量
					if len(customNames) != submatchCount {
						fmt.Printf("错误: 需要%d个标签，但输入了%d个\n", submatchCount, len(customNames))
						continue
					}
					matches := reg.FindAllStringSubmatch(jj, len(label_name)+1)
					if matches == nil {
						fmt.Println("没有找到匹配项")
						continue
					}
					for _, m := range matches {
						for i := 1; i < len(m); i++ {
							name := customNames[i-1]
							aaa[name] = m[i]
						}
					}
					value, _ := strconv.Atoi((aaa[value_label]))
					valuefloat64 := float64(value)
					err = GetdatetargetlogMetricnew(RequestLogFilternew, url, i, aaa["traceidwithskywalking"], aaa["traceidwiithaloudata"], aaa["responsecode"], aaa["starttime"], aaa["url"], valuefloat64)
					if err != nil {
						panic(err)
					}
				}

			}

			duration, err := time.ParseDuration(config.GetString("mixedformat.cycle"))
			// fmt.Println(config.GetString("postgres.metrics.cycle"))
			if err != nil {
				panic(err)
			}
			time.Sleep(duration)
			RequestLogFilternew.Reset()
		}
	}()

	prometheus.MustRegister(RequestLogFilternew)
}

// 获取报错记录数
func ReLogNumErrorwithurl(config *viper.Viper) {
	// defer db.Close()
	log.Println("调用接口产出日志")
	// logs_lists := config.GetStringMap("mixedformat.requestlogs")
	// 获取待收集日志指标项目
	logs_lists := config.GetStringSlice("mixedformat.errorlogs")
	// 日志指标项统计区间
	collectionscopeseconds := config.GetString("mixedformat.collectionscopeseconds")
	// 设置延迟时间，主要是为了防止loli日志不能及时被查询，延迟时间间隔进行日志收取
	latencycollectionseconds := config.GetString("mixedformat.latencycollectionseconds")
	// loki地址
	lokiipport := config.GetString("mixedformat.lokiipport")
	RequestLogFiltererror := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "error_logs",
			Help: "collect error logs",
		},
		[]string{"name", "url", "traceidwithskywalking", "traceidwiithaloudata", "errormsg", "starttime"},
	)

	go func() error {
		for {
			for _, i := range logs_lists {
				// log_pro := i
				label_list := config.GetStringMap(i + ".label_list")
				lokire := config.GetStringSlice(i + ".lokire")
				lokiexclre := config.GetStringSlice(i + ".lokiexclre")
				lokire_string := ""
				// 是否进行string模糊匹配
				if len(lokire) != 0 {
					for _, j := range lokire {
						lokire_string += "|~`" + j + "`"
					}
				}
				// 是否进行string 剔除匹配
				if len(lokiexclre) != 0 {
					for _, j := range lokiexclre {
						lokire_string += "!~`" + j + "`"
					}
				}
				// 获取每次查询loki的日志条数限制
				recordslimit := config.GetString(i + ".recordslimit")
				// 进行loki查询,返回匹配的值列表
				data, url, _ := Getlogandurlfromloki(lokiipport, label_list, time.Now(), latencycollectionseconds, collectionscopeseconds, recordslimit, lokire_string)

				log.Println(i, "的数据是:", len(data))
				for _, jj := range data {
					var aaa map[string]string
					aaa = make(map[string]string)

					// reg, err = regexp.Compile("(\\d+-\\d+-\\d+\\s\\S+)\\s\\[(.*?)\\]\\s\\[.*?\\]\\s(\\w+)\\s+(\\S+)\\s-\\s(.*)")
					regex := config.GetString(i + ".regex")
					reg, err := regexp.Compile(regex)
					if err != nil {
						fmt.Printf("正则表达式编译错误: %v\n", err)
						continue
					}
					// 获取子匹配总数
					submatchCount := reg.NumSubexp()
					if submatchCount == 0 {
						fmt.Println("警告: 正则表达式不包含捕获组")
						continue
					}
					label_name := config.GetString("mixedformat.errorlogslabel_name")
					// value_label := config.GetString(i + ".value_label")

					customNames := strings.Split(label_name, " ")
					// 验证名称数量
					if len(customNames) != submatchCount {
						fmt.Printf("错误: 需要%d个标签，但输入了%d个\n", submatchCount, len(customNames))
						continue
					}
					matches := reg.FindAllStringSubmatch(jj, len(customNames)+1)
					if matches == nil {
						fmt.Println("没有找到匹配项")
						continue
					}
					for _, m := range matches {
						for i := 1; i < len(m); i++ {
							name := customNames[i-1]
							aaa[name] = m[i]
						}
					}
					valuefloat64 := float64(1)
					err = GetdatetargetlogMetricerror(RequestLogFiltererror, i, url, aaa["traceidwithskywalking"], aaa["traceidwiithaloudata"], aaa["errormsg"], aaa["starttime"], valuefloat64)
					if err != nil {
						panic(err)
					}
				}

			}

			duration, err := time.ParseDuration(config.GetString("mixedformat.cycle"))
			// fmt.Println(config.GetString("postgres.metrics.cycle"))
			if err != nil {
				panic(err)
			}
			time.Sleep(duration)
			RequestLogFiltererror.Reset()
		}
	}()
	prometheus.MustRegister(RequestLogFiltererror)
}

// 获取查询各个阶段的耗时指标
func ReLogCostquerystageswithurl(config *viper.Viper) {
	// defer db.Close()
	log.Println("调用接口产出日志")
	// logs_lists := config.GetStringMap("mixedformat.requestlogs")
	// 获取待收集日志指标项目
	logs_lists := config.GetStringSlice("mixedformat.queryphasetimelogs")
	// 日志指标项统计区间
	collectionscopeseconds := config.GetString("mixedformat.collectionscopeseconds")
	// 设置延迟时间，主要是为了防止loli日志不能及时被查询，延迟时间间隔进行日志收取
	latencycollectionseconds := config.GetString("mixedformat.latencycollectionseconds")
	// loki地址
	lokiipport := config.GetString("mixedformat.lokiipport")
	RequestLogFilter := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "query_stage_cost",
			Help: "Collect error logs",
		},
		[]string{"name", "url", "traceidwithskywalking", "traceidwiithaloudata", "stages", "starttime"},
	)

	go func() error {
		for {
			for _, i := range logs_lists {
				// log_pro := i
				label_list := config.GetStringMap(i + ".label_list")
				lokire := config.GetStringSlice(i + ".lokire")
				lokiexclre := config.GetStringSlice(i + ".lokiexclre")
				lokire_string := ""
				// 是否进行string模糊匹配
				if len(lokire) != 0 {
					for _, j := range lokire {
						lokire_string += "|~`" + j + "`"
					}
				}
				// 是否进行string 剔除匹配
				if len(lokiexclre) != 0 {
					for _, j := range lokiexclre {
						lokire_string += "!~`" + j + "`"
					}
				}
				// 获取每次查询loki的日志条数限制
				recordslimit := config.GetString(i + ".recordslimit")
				// 进行loki查询,返回匹配的值列表
				data, url, _ := Getlogandurlfromloki(lokiipport, label_list, time.Now(), latencycollectionseconds, collectionscopeseconds, recordslimit, lokire_string)

				log.Println(i, "的数据是:", len(data))
				for _, jj := range data {
					var aaa map[string]string
					aaa = make(map[string]string)

					// reg, err = regexp.Compile("(\\d+-\\d+-\\d+\\s\\S+)\\s\\[(.*?)\\]\\s\\[.*?\\]\\s(\\w+)\\s+(\\S+)\\s-\\s(.*)")
					regex := config.GetString(i + ".regex")
					reg, err := regexp.Compile(regex)
					if err != nil {
						fmt.Printf("正则表达式编译错误: %v\n", err)
						continue
					}
					// 获取子匹配总数
					submatchCount := reg.NumSubexp()
					if submatchCount == 0 {
						fmt.Println("警告: 正则表达式不包含捕获组")
						continue
					}
					label_name := config.GetString("mixedformat.queryphasetimelogslabel_name")
					value_label := config.GetString(i + ".value_label")

					customNames := strings.Split(label_name, " ")
					// 验证名称数量
					if len(customNames) != submatchCount {
						fmt.Printf("错误: 需要%d个标签，但输入了%d个\n", submatchCount, len(customNames))
						continue
					}
					matches := reg.FindAllStringSubmatch(jj, len(label_name)+1)
					if matches == nil {
						fmt.Println("没有找到匹配项")
						continue
					}
					for _, m := range matches {
						for i := 1; i < len(m); i++ {
							name := customNames[i-1]
							aaa[name] = m[i]
						}
					}

					value, _ := strconv.Atoi((aaa[value_label]))
					valuefloat64 := float64(value)
					err = GetdatetargetlogMetricquerystage(RequestLogFilter, url, i, aaa["traceidwithskywalking"], aaa["traceidwiithaloudata"], aaa["stages"], aaa["starttime"], valuefloat64)
					if err != nil {
						panic(err)
					}
				}

			}

			duration, err := time.ParseDuration(config.GetString("mixedformat.cycle"))
			// fmt.Println(config.GetString("postgres.metrics.cycle"))
			if err != nil {
				panic(err)
			}
			time.Sleep(duration)
			RequestLogFilter.Reset()
		}
	}()

	prometheus.MustRegister(RequestLogFilter)
}

func NacosMustRegisterOnce(config *viper.Viper) {
	servicelist := config.GetStringSlice("global.servicelist")
	NacosServerstatusDuration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_status_on_nacos",
			Help: "The number of producers serving on Nacos does not meet expectations at 0, but meets expectations at 1",
		},
		[]string{"currentcount", "normalcount", "servicename"},
	)

	NacosProducernumDuration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nacos_producer_count",
			Help: "List the status of all nacos producers. 1 is not a sufficient quantity but available, 2 is expected, and 0 is unavailable.",
		},
		[]string{"healthytcount", "producer_name", "normalcount", "servicename"},
	)

	NacosNormalProducersumDuration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_normal_producer_sum",
			Help: "Number of service normal producers",
		},
		[]string{"servicename"},
	)
	NacosCurrentProducersumDuration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_current_producer_sum",
			Help: "Number of service current producers",
		},
		[]string{"servicename"},
	)

	go func() error {
		for {
			for _, j := range servicelist {
				// NacosServerstatusDuration = prometheus.NewGaugeVec(
				// 	prometheus.GaugeOpts{
				// 		Name: "server_status_on_nacos",
				// 		Help: "The number of producers serving on Nacos does not meet expectations at 0, but meets expectations at 1",
				// 	},
				// 	[]string{"currentcount", "normalcount", "servicename"},
				// )

				currentserviceproducer := GetserviceproducerlistJson(config, j)
				serviceProducerList := config.GetStringMapString(j + ".serviceList")
				status, currentcount, normalcount := NacosServerContrast(currentserviceproducer, serviceProducerList)
				SetNacosServerstatus(NacosServerstatusDuration, currentcount, normalcount, status, j)
				SetNacosNormalProducersum(NacosNormalProducersumDuration, j, normalcount)
				SetNacosCurrentProducersum(NacosCurrentProducersumDuration, j, currentcount)

				for k, f := range serviceProducerList {
					// NacosProducernumDuration = prometheus.NewGaugeVec(
					// 	prometheus.GaugeOpts{
					// 		Name: "nacos_producer_count",
					// 		Help: "List the status of all nacos producers. 1 is not a sufficient quantity but available, 2 is expected, and 0 is unavailable.",
					// 	},
					// 	[]string{"healthytcount", "producer_name", "normalcount", "servicename"},
					// )
					status1, healthytcount1, normalcount1, producer_name1 := NacosServerProducerContrast(currentserviceproducer, k, f)
					// fmt.Println("当前生产者名称是", producer_name1)
					SetNacosProducernum(NacosProducernumDuration, healthytcount1, normalcount1, status1, producer_name1, j)
				}

			}
			duration, err := time.ParseDuration(config.GetString("mixedformat.cycle"))
			// fmt.Println(config.GetString("postgres.metrics.cycle"))
			if err != nil {
				panic(err)
			}
			time.Sleep(duration)
			NacosServerstatusDuration.Reset()
			NacosProducernumDuration.Reset()
			NacosNormalProducersumDuration.Reset()
			NacosCurrentProducersumDuration.Reset()
		}
	}()
	prometheus.MustRegister(NacosServerstatusDuration)
	prometheus.MustRegister(NacosProducernumDuration)
	prometheus.MustRegister(NacosNormalProducersumDuration)
	prometheus.MustRegister(NacosCurrentProducersumDuration)

}
