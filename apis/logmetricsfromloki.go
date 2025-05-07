package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"monitor/methods"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// swagger:route UpdatetargetlogMetrics
//
// This will show all available pets by default.
//
//	    Schemes: http, https
//
//	    Responses:
//	      503: "执行删除alarm " + id + "任务失败"
//	      200: "删除alarm " + id + " 成功"
//		  504: 等待进程退出失败

func UpdatetargetlogMetrics(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("调用接口产出日志")
		var result methods.Resp

		// logs_lists := config.GetStringMap("mixedformat.requestlogs")
		// 获取待收集日志指标项目
		logs_lists := config.GetStringSlice("mixedformat.requestlogs")
		// 日志指标项统计区间
		collectionscopeseconds := config.GetString("mixedformat.collectionscopeseconds")
		// 设置延迟时间，主要是为了防止loli日志不能及时被查询，延迟时间间隔进行日志收取
		latencycollectionseconds := config.GetString("mixedformat.latencycollectionseconds")
		// loki地址
		lokiipport := config.GetString("mixedformat.lokiipport")

		log.Println(logs_lists)
		var count int
		for _, i := range logs_lists {
			// log_pro := i
			f, _ := os.Create(i + "LogStatus.txt")
			defer f.Close()
			label_list := config.GetStringMap(i + ".label_list")
			lokire := config.GetStringSlice(i + ".lokire")
			lokiexclre := config.GetStringSlice(i + ".lokiexclre")
			value_label := config.GetStringSlice(i + ".value_label")
			showlabelvaluelists := config.GetStringSlice(i + ".showlabelvaluelists")
			label_name := config.GetStringSlice(i + ".label_name")

			label_len, _ := strconv.Atoi(config.GetString("mixedformat.label_len"))
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
			data, err := methods.Getlogfromlokippre(lokiipport, label_list, time.Now(), latencycollectionseconds, collectionscopeseconds, recordslimit, lokire_string)
			if err != nil {
				result.Msg = err.Error()
				result.Code = "400"
				break
			}

			log.Println(i, "的数据是:", len(data))
			if len(data) != 0 {
				var errorloglist string

				// errorlognum := "errorlognum{name=\"" + i + "\",type=\"" + lokire_string + "\"} " + strconv.Itoa(len(data)) + "\n"
				errorlognum := "errorlognum{name=\"" + i + "\",type=\"" + lokire_string + "\""
				// 查看需要生成的指标项的数量
				if len(showlabelvaluelists) != 0 {
					var showlabelvaluedict = make(map[string][]string, len(showlabelvaluelists))
					for _, showlabelvalue := range showlabelvaluelists {
						showlabelvaluedict[showlabelvalue] = make([]string, 0)
					}
					var reg *regexp.Regexp
					for _, jj := range data {
						count++
						// reg, err = regexp.Compile("(\\d+-\\d+-\\d+\\s\\S+)\\s\\[(.*?)\\]\\s\\[.*?\\]\\s(\\w+)\\s+(\\S+)\\s-\\s(.*)")
						regex := config.GetString(i + ".regex")
						reg, err = regexp.Compile(regex)
						matches := reg.FindAllStringSubmatch(jj, len(label_name)+1)

						for _, m := range matches {
							// errorloglist += i + "{re=\"" + lokire_string + "\",flag=\"" + strconv.Itoa(count)
							// if len(m) == len(label_name)+1 {
							// 	for iii, jjj := range m[1:] {
							// 		errorloglist += "\"," + label_name[iii] + "=\"" + methods.KeepFirstTenCharacters(jjj, 128)
							// 	}
							// 	errorloglist += "\"} " + "1" + "\n"
							// }

							errorloglist += i + "{re=\"" + lokire_string + "\",flag=\"" + strconv.Itoa(count) + "\""
							if len(m) == len(label_name)+1 {

								if len(value_label) != 1 {
									for iii, jjj := range m[1:] {

										if methods.StrInSlice(label_name[iii], showlabelvaluelists) {
											if !methods.StrInSlice(jjj, showlabelvaluedict[label_name[iii]]) && jjj != "" {
												showlabelvaluedict[label_name[iii]] = append(showlabelvaluedict[label_name[iii]], jjj)
											}
										}
										errorloglist += "," + label_name[iii] + "=" + methods.KeepFirstTenCharacters(jjj, label_len)
									}
									errorloglist += "} " + "1" + "\n"
								} else {
									var value string
									for iii, jjj := range m[1:] {

										if methods.StrInSlice(label_name[iii], showlabelvaluelists) {
											if !methods.StrInSlice(jjj, showlabelvaluedict[label_name[iii]]) && jjj != "" {
												showlabelvaluedict[label_name[iii]] = append(showlabelvaluedict[label_name[iii]], jjj)
											}
										}
										errorloglist += "," + label_name[iii] + "=" + methods.KeepFirstTenCharacters(jjj, label_len)
										if value_label[0] == label_name[iii] {
											value = jjj
										}
									}
									errorloglist += "} " + value + "\n"
								}

							}

						}
					}
					for lable_name, value_slice := range showlabelvaluedict {
						// errorlognum += "," + lable_name + "=\"" + methods.KeepFirstTenCharacters(fmt.Sprintf(strings.Join(value_slice, ",")),label_len)
						errorlognum += "," + lable_name + "=" + methods.KeepFirstTenCharacters(fmt.Sprintf(strings.Join(value_slice, ",")), label_len)
					}
					errorlognum += "} " + strconv.Itoa(len(data)) + "\n"

				} else {
					var reg *regexp.Regexp
					for _, jj := range data {
						count++
						// reg, err = regexp.Compile("(\\d+-\\d+-\\d+\\s\\S+)\\s\\[(.*?)\\]\\s\\[.*?\\]\\s(\\w+)\\s+(\\S+)\\s-\\s(.*)")
						regex := config.GetString(i + ".regex")
						reg, err = regexp.Compile(regex)
						matches := reg.FindAllStringSubmatch(jj, len(label_name)+1)

						for _, m := range matches {
							// errorloglist += i + "{re=\"" + lokire_string + "\",flag=\"" + strconv.Itoa(count)
							// if len(m) == len(label_name)+1 {
							// 	for iii, jjj := range m[1:] {
							// 		errorloglist += "\"," + label_name[iii] + "=\"" + methods.KeepFirstTenCharacters(jjj, 128)
							// 	}
							// 	errorloglist += "\"} " + "1" + "\n"
							// }

							errorloglist += i + "{re=\"" + lokire_string + "\",flag=\"" + strconv.Itoa(count) + "\""
							if len(m) == len(label_name)+1 {

								if len(value_label) != 1 {
									for iii, jjj := range m[1:] {

										errorloglist += "," + label_name[iii] + "=" + methods.KeepFirstTenCharacters(jjj, label_len)
									}
									errorloglist += "} " + "1" + "\n"
								} else {
									var value string
									for iii, jjj := range m[1:] {

										errorloglist += "," + label_name[iii] + "=" + methods.KeepFirstTenCharacters(jjj, label_len)
										if value_label[0] == label_name[iii] {
											value = jjj
										}
									}
									errorloglist += "} " + value + "\n"
								}

							}

						}
					}
					errorlognum += "} " + strconv.Itoa(len(data)) + "\n"
				}

				_, err = f.Write([]byte(errorlognum))
				if err != nil {
					log.Println("写入" + i + "错误日志数量失败 退出")
					return
				}
				_, err = f.Write([]byte(errorloglist))
				if err != nil {
					log.Println("写入" + i + "错误详细失败 退出")
					return
				}

			}
			err = methods.UpdateNacosMetrics(config, f.Name(), i+"metrics_errors")
			if err != nil {
				result.Msg = err.Error()
				result.Code = "400"
			}
		}

		if result.Code != "400" {
			result.Msg = "调用成功。"
			result.Code = "200"
		}
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}
