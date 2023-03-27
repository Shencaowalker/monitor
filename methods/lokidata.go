package methods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type LokiQuery_range struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream struct {
				ElasticEnv  string `json:"elasticEnv"`
				Environment string `json:"environment"`
				Filename    string `json:"filename"`
				Hostname    string `json:"hostname"`
				Job         string `json:"job"`
				Project     string `json:"project"`
			} `json:"stream"`
			Values [][]string `json:"values"`
		} `json:"result"`
		Stats struct {
			Summary struct {
				BytesProcessedPerSecond int     `json:"bytesProcessedPerSecond"`
				LinesProcessedPerSecond int     `json:"linesProcessedPerSecond"`
				TotalBytesProcessed     int     `json:"totalBytesProcessed"`
				TotalLinesProcessed     int     `json:"totalLinesProcessed"`
				ExecTime                float64 `json:"execTime"`
				QueueTime               float64 `json:"queueTime"`
				Subqueries              int     `json:"subqueries"`
				TotalEntriesReturned    int     `json:"totalEntriesReturned"`
			} `json:"summary"`
			Querier struct {
				Store struct {
					TotalChunksRef        int `json:"totalChunksRef"`
					TotalChunksDownloaded int `json:"totalChunksDownloaded"`
					ChunksDownloadTime    int `json:"chunksDownloadTime"`
					Chunk                 struct {
						HeadChunkBytes    int `json:"headChunkBytes"`
						HeadChunkLines    int `json:"headChunkLines"`
						DecompressedBytes int `json:"decompressedBytes"`
						DecompressedLines int `json:"decompressedLines"`
						CompressedBytes   int `json:"compressedBytes"`
						TotalDuplicates   int `json:"totalDuplicates"`
					} `json:"chunk"`
				} `json:"store"`
			} `json:"querier"`
			Ingester struct {
				TotalReached       int `json:"totalReached"`
				TotalChunksMatched int `json:"totalChunksMatched"`
				TotalBatches       int `json:"totalBatches"`
				TotalLinesSent     int `json:"totalLinesSent"`
				Store              struct {
					TotalChunksRef        int `json:"totalChunksRef"`
					TotalChunksDownloaded int `json:"totalChunksDownloaded"`
					ChunksDownloadTime    int `json:"chunksDownloadTime"`
					Chunk                 struct {
						HeadChunkBytes    int `json:"headChunkBytes"`
						HeadChunkLines    int `json:"headChunkLines"`
						DecompressedBytes int `json:"decompressedBytes"`
						DecompressedLines int `json:"decompressedLines"`
						CompressedBytes   int `json:"compressedBytes"`
						TotalDuplicates   int `json:"totalDuplicates"`
					} `json:"chunk"`
				} `json:"store"`
			} `json:"ingester"`
		} `json:"stats"`
	} `json:"data"`
}

//调接口得到时间区间内日志集合 参数:服务名称、环境、文件名称、
func Getlogfromloki2(config *viper.Viper, servicename string) ProducerJson {
	var producerjson ProducerJson
	url := "http://" + config.GetString("global.nacosip") + ":" + config.GetString("global.nacosport") + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=" + config.GetString("global.pageNo") + "&pageSize=" + config.GetString("global.pageSize") + "&serviceNameParam=providers.*" + servicename + "&groupNameParam=&namespaceId=" + config.GetString("global.namespaceId")
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("wraps ", servicename, "nacos producers url err.\n")
		return producerjson
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get", servicename, "producers err.\n")
		return producerjson
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("read resp error,", servicename, "producers read err.\n")
		return producerjson
	}

	json.Unmarshal(body, &producerjson)
	return producerjson
}

func Getlogfromloki(lokiipport string, label_list map[string]interface{}, latencycollectionseconds string, collectionscopeseconds string) (loglists []string) {
	latencyseconds, _ := time.ParseDuration("-" + latencycollectionseconds + "s")
	latencycollectiontime := time.Now().Add(1 * latencyseconds)
	latencycollectiontimeformat := latencycollectiontime.UnixNano()

	scopeseconds, _ := time.ParseDuration(collectionscopeseconds + "s")
	collectionscopetime := latencycollectiontime.Add(1 * scopeseconds).UnixNano()

	// fmt.Println("latencycollectiontime:", latencycollectiontimeformat)
	// fmt.Println("collectionscopetime:", collectionscopetime)
	// fmt.Println("nowtime:", time.Now().UnixNano())

	// fmt.Println("environment:", label_list["environment"])
	// fmt.Println("job:", label_list["job"])
	// fmt.Println("filename:", label_list["filename"])
	// fmt.Println("latencycollectiontime:", latencycollectiontimeformat)

	var a string
	for i, j := range label_list {
		a = a + i + "=\"" + j.(string) + "\","
	}
	url := "http://" + lokiipport + "/loki/api/v1/query_range?query={" + a[:len(a)-1] + "}&start=" + strconv.Itoa(int(latencycollectiontimeformat)) + "&end=" + strconv.Itoa(int(collectionscopetime)) + "&limit=8000"
	fmt.Println(url)
	// url := "http://" + lokiipport + "/loki/api/v1/query_range?query={job=\"" + "chaos" + "\"}&start=1672816117813000000&end=1672902517813000000&limit=8000"
	client := &http.Client{}
	var lokiquery_range LokiQuery_range
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("wraps: create ", url, "request  error.")
		return loglists
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get", url, " err.\n")
		return loglists
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("read resp error log read err.\n")
		return loglists
	}

	json.Unmarshal(body, &lokiquery_range)
	for _, i := range lokiquery_range.Data.Result {
		for _, j := range i.Values {
			loglists = append(loglists, j[1])
		}
	}
	fmt.Println(len(loglists))
	fmt.Println(url)
	return loglists
}

func SplitoneLinetometrics(metricname string, collist []interface{}, value_col string, line string, re string) string {
	arrs := strings.Split(line, re)
	servicemetric := metricname + "{"
	var drift int
	for i, j := range collist {
		servicemetric = servicemetric + j.(string) + "=\"" + arrs[i] + "\","
		if value_col == j.(string) {
			drift = i
		}
	}
	servicemetric = servicemetric[:len(servicemetric)-2] + "\"} " + arrs[drift] + "\n"
	return servicemetric
}

func SplitoneJoinsightLinetometrics(metricname string, collist []interface{}, value_col string, line string, re string) string {
	arrs := strings.Split(line, re)
	servicemetric := metricname + "{"
	var value_value string
	for i, j := range collist {
		if j.(string) == "jsoncontent" {
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(arrs[i]), &event); err != nil {
				fmt.Println("json Unmarshal error!")
			}
			fmt.Println(event)
			for k, v := range event {
				servicemetric = servicemetric + k + "=\"" + fmt.Sprintf("%v", v) + "\","
			}
			value_value = fmt.Sprintf("%v", event[value_col])
			continue
		}
		servicemetric = servicemetric + j.(string) + "=\"" + arrs[i] + "\","
	}
	servicemetric = servicemetric[:len(servicemetric)-2] + "\"} " + value_value + "\n"
	return servicemetric
}

func SplitoneLinetopostgresql(line string, re string) []interface{} {
	var a []interface{}
	arrs := strings.Split(line, re)
	for _, j := range arrs {
		a = append(a, j)
	}
	return a
}
