package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"monitor/methods"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Tag "成功"
func TransferToLoki(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		request_header := request.Header
		// fmt.Println("Header全部数据:", header["Origin"][0])
		if request.Method == "OPTIONS" {
			writer.Header().Set("Access-Control-Allow-Origin", request_header["Origin"][0])
			// writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Content-Type", "application/json")
			writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "content-type")
			writer.Header().Set("Access-Control-Max-Age", "31536000")
			// writer.Header().Set("Referer", request.RemoteAddr)
			// writer.Header().Set("Access-Control-Allow-Headers", "content-type")
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			// writer.Header().Set("allowedCredentials", "true")
			writer.WriteHeader(200)
			writer.Write([]byte("no content\n"))
			log.Println("Received options request")
			return
		}
		var lokipushdata methods.LokiPushData
		var receiveData methods.ReceiveDataUploadedToLoki
		if err := json.NewDecoder(request.Body).Decode(&receiveData); err != nil {
			request.Body.Close()
			log.Println(err)
			// result := methods.Sendmessage(config, alarmdata)
			// if err := json.NewEncoder(writer).Encode(result); err != nil {
			// 	log.Fatal(err)
			// methods.Sendmessage(config, alarmdata)
		}
		nowtime := time.Now()
		currenttimeformat := nowtime.Format("2006-01-02 15:04:05.999")
		currenttimestamp := strconv.FormatInt(nowtime.UnixNano(), 10)
		content := currenttimeformat + " " + receiveData.Type + " " + receiveData.Content
		Streams := make([]methods.Streams, 1)
		Value := make([]string, 2)
		Values := make([][]string, 1)
		Value[0] = currenttimestamp
		Value[1] = content
		Values[0] = Value
		Streams[0].Stream = receiveData.Labels
		Streams[0].Values = Values

		lokipushdata.Streams = Streams
		// fmt.Println(lokipushdata)
		jsonlokipushdata, err := json.Marshal(lokipushdata)
		if err != nil {
			fmt.Println("JSON marshaling failed:", err)
			return
		}
		fmt.Println(string(jsonlokipushdata))
		url := "http://" + config.GetString("airerrorserverlogtopg.lokiipport") + "/loki/api/v1/push"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonlokipushdata))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			writer.Header().Set("Access-Control-Allow-Origin", request_header["Origin"][0])
			writer.Header().Set("Content-Type", "application/json")
			writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "content-type")
			writer.Header().Set("Access-Control-Max-Age", "31536000")
			writer.WriteHeader(403)
			writer.Write([]byte(err.Error()))
		}
		writer.Header().Set("Access-Control-Allow-Origin", request_header["Origin"][0])
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "content-type")
		writer.Header().Set("Access-Control-Max-Age", "31536000")
		// writer.Header().Set("allowedCredentials", "true")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.WriteHeader(200)
		writer.Write([]byte("seccess\n"))
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
	}
}
