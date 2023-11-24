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

func TransferToLoki(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
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
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
	}
}
