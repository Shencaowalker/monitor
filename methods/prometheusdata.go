package methods

import (
	"log"
	"os/exec"

	"github.com/spf13/viper"
)

//更新pushgateway指标信息 改成根据指标来进行
func UpdatePushgatewayMetrics(config *viper.Viper, pushgatewayjobname string) {
	deletecmd := exec.Command("curl", "-XDELETE", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/"+pushgatewayjobname)
	err := deletecmd.Run()
	if err != nil {
		log.Println("删除 ", pushgatewayjobname, " 指标报错 err 继续执行")
	} else {
		log.Println("DELETE ", pushgatewayjobname, " SECCESS")
	}
	pushcmd := exec.Command("curl", "-XPOST", "--data-binary", "@Status.txt", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/"+pushgatewayjobname)
	err = pushcmd.Run()
	if err != nil {
		log.Println("上传pushgateway job ", pushgatewayjobname, " 报错 ", err.Error(), "; 继续执行")
	} else {
		log.Println("UPLOAD ", pushgatewayjobname, " SECCESS")
	}
}
