package methods

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"

	"github.com/spf13/viper"
)

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// 注册用监控项内容
type Registration_nformation struct {
	Id       string `json:"id"`
	Group    string `json:"group"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Tags     string `json:"tags"`
	Phone    string `json:"phone"`
	Env      string `json:"env"`
	M_type   string `json:"m_type"`
	App_type string `json:"App_type"`
}

// 注册用监控项内容列表
type Registration_nformations struct {
	Values []Registration_nformation `json:"values"`
}

// 下线监控项列表
type Ddownline_nformations struct {
	Itemids []string `"itemids"`
}

// 下线高境项列表
type Ddownline_alarms struct {
	Alarmids []string `"alarmids"`
}

// 注册告警项内容
type Registration_Alarm struct {
	Alert  string `json:"alert"`
	Expr   string `json:"expr"`
	For    string `json:"for"`
	Labels struct {
		Severity string `json:"severity"`
		Env      string `json:"env"`
		Project  string `"json:"project"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
}

// 注册告警项内容
type Registration_Alarms struct {
	Values []Registration_Alarm `json:"values"`
}

// swagger:route isPhoneNum
//
//	手机号格式校验
//
// This will show all available pets by default.
//
// return  true or false
func isPhoneNum(s string) bool {
	// 手机号正则表达式
	pattern := `^1[3456789]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(s)
}

// swagger:route ConsulregisterItem
//
//	获取请求json 拆分后进行consul 接口内容重组并注册监控项
//
// This will show all available pets by default.
//
//	    Schemes: http, https
//
//	    Responses:
//	      503: 执行注册consul item任务项 任务失败
//	      200: 注册item成功/返回consul自己返回的返回值
//		  504: 等待进程退出失败
func ConsulregisterItem(config *viper.Viper, information Registration_nformation) (result Resp) {
	var json_value string
	if isPhoneNum(information.Phone) {
		json_value = "{\"id\": \"" + information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + "\",\"name\": \"" + information.Group + "_" + information.Tags + "\",\"address\": \"" + information.Address + "\",\"port\": " + information.Port + ",\"tags\": [\"" + information.Tags + "\"],\"meta\":{\"env\":\"" + information.Env + "\",\"m_type\":\"" + information.M_type + "\",\"app\":\"" + information.App_type + "\",\"phone\":\"" + information.Phone + "\"},\"checks\": [{\"" + information.M_type + "\": \"" + information.Address + ":" + information.Port + "\", \"interval\": \"60s\"}]}"
	} else {
		json_value = "{\"id\": \"" + information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + "\",\"name\": \"" + information.Group + "_" + information.Tags + "\",\"address\": \"" + information.Address + "\",\"port\": " + information.Port + ",\"tags\": [\"" + information.Tags + "\"],\"meta\":{\"env\":\"" + information.Env + "\",\"m_type\":\"" + information.M_type + "\",\"app\":\"" + information.App_type + "\"},\"checks\": [{\"" + information.M_type + "\": \"" + information.Address + ":" + information.Port + "\", \"interval\": \"60s\"}]}"
	}
	log.Println("注册consul item 的json:" + json_value)
	registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/agent/service/register")
	stdout, _ := registrationcmd.StdoutPipe()
	err := registrationcmd.Start()
	if err != nil {
		fmt.Println("执行注册consul item 任务失败,错误详情：", err)
		result.Msg = "执行注册consul item 任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "注册item成功"
		}
	}
	if err := registrationcmd.Wait(); err != nil { //避免僵尸进程
		fmt.Println("等待进程退出失败,错误详情：", err)
		result.Msg = "等待进程退出失败"
		result.Code = "504"
	}
	log.Println(information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + result.Msg)
	return
}

// swagger:route ConsuldownlineItem
//
//	获取consul 注册的单个任务id ,进行监控任务项下线
//
// This will show all available pets by default.
//
//	    Schemes: http, https
//
//	    Responses:
//	      503: "执行删除items " + id + "任务失败"
//	      200: 删除consul任务成功/返回consul自己返回的返回值
//		  504: 等待进程退出失败
func ConsuldownlineItem(config *viper.Viper, id string) (result Resp) {
	deleteservicecmd := exec.Command("curl", "-XPUT", "http://"+config.GetString("global.consulipport")+"/v1/agent/service/deregister/"+id)
	stdout, _ := deleteservicecmd.StdoutPipe()
	defer stdout.Close()
	err := deleteservicecmd.Start()
	if err != nil {
		fmt.Println("执行删除items "+id+"任务失败，错误详情：", err)
		result.Msg = "执行删除items " + id + "任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "删除consul任务成功"
		}
	}
	if err := deleteservicecmd.Wait(); err != nil { //避免僵尸进程
		fmt.Println("等待进程退出失败,错误详情：", err)
		result.Msg = "等待进程退出失败"
		result.Code = "504"
	}
	log.Println(id + result.Msg)
	return
}

// swagger:route ConsulregisterAlarm
//
//	获取请求json 拆分后进行consul 接口内容重组并注册告警指标
//
// This will show all available pets by default.
//
//	    Schemes: http, https
//
//	    Responses:
//	      503: 执行注册consul alarm 任务失败
//	      200: 注册alarm成功/返回consul自己返回的返回值
//		  504: 等待进程退出失败
func ConsulregisterAlarm(config *viper.Viper, information Registration_Alarm) (result Resp) {
	json_value := "{\"alert\": \"" + information.Alert + "\",\"expr\": \"" + information.Expr + "\",\"for\": \"" + information.For + "\",\"labels\":{\"severity\":\"" + information.Labels.Severity + "\",\"env\":\"" + information.Labels.Env + "\",\"project\":\"" + information.Labels.Project + "\"},\"annotations\":{\"description\":\"" + information.Annotations.Description + "\",\"summary\":\"" + information.Annotations.Summary + "\"}}"
	log.Println("注册consul alarn 的json:" + json_value)
	registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/kv/prometheus/rules/"+information.Alert)
	stdout, _ := registrationcmd.StdoutPipe()
	err := registrationcmd.Start()
	if err != nil {
		fmt.Println("执行注册consul alarm 任务失败,错误详情：", err)
		result.Msg = "执行注册consul alarm 任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "注册alarm成功"
		}
	}
	if err := registrationcmd.Wait(); err != nil { //避免僵尸进程
		fmt.Println("等待进程退出失败,错误详情：", err)
		result.Msg = "等待进程退出失败"
		result.Code = "504"
	}
	log.Println(information.Alert + " is " + result.Msg)
	return
}

// swagger:route ConsuldownlineAlarm
//
//	获取consul注册的单个告警项目 ,进行告警任务项下线
//
// This will show all available pets by default.
//
//	    Schemes: http, https
//
//	    Responses:
//	      503: "执行删除alarm " + id + "任务失败"
//	      200: "删除alarm " + id + " 成功"
//		  504: 等待进程退出失败
func ConsuldownlineAlarm(config *viper.Viper, id string) (result Resp) {
	deleteservicecmd := exec.Command("curl", "-XDELETE", "http://"+config.GetString("global.consulipport")+"/v1/kv/prometheus/rules/"+id)
	stdout, _ := deleteservicecmd.StdoutPipe()
	defer stdout.Close()
	err := deleteservicecmd.Start()
	if err != nil {
		fmt.Println("执行删除alarm "+id+"任务失败，错误详情：", err)
		result.Msg = "执行删除alarm " + id + "任务失败"
		result.Code = "503"
	} else {
		// res, _ := ioutil.ReadAll(stdout)
		// resdata := string(res)
		// if resdata != "" {
		// 	result.Code = "200"
		// 	result.Msg = resdata
		// } else {
		result.Code = "200"
		result.Msg = "删除alarm " + id + " 成功"
		// }
	}
	if err := deleteservicecmd.Wait(); err != nil { //避免僵尸进程
		fmt.Println("等待进程退出失败,错误详情：", err)
		result.Msg = "等待进程退出失败"
		result.Code = "504"
	}
	log.Println("delete alarm " + id + " " + result.Msg)
	return
}
