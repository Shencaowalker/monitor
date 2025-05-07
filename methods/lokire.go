package methods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Getlogfromlokippre(lokiipport string, label_list map[string]interface{}, timenow time.Time, latencycollectionseconds string, collectionscopeseconds string, recordslimit string, lokire_string string) (loglists []string, err error) {
	latencyseconds, _ := time.ParseDuration("-" + latencycollectionseconds + "s")
	latencycollectiontime := timenow.Add(1 * latencyseconds)
	latencycollectiontimeformat := ((latencycollectiontime.UnixNano()) / 1000000000) * 1000000000

	scopeseconds, _ := time.ParseDuration(collectionscopeseconds + "s")
	collectionscopetime := ((latencycollectiontime.Add(1*scopeseconds).UnixNano())/1000000000)*1000000000 - 1

	log.Println("开始时间：", latencycollectiontimeformat, "结束时间：", collectionscopetime)

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
		log.Println("访问的loki接口是", a)
	}
	var url string

	if lokire_string != "" {
		url = "http://" + lokiipport + "/loki/api/v1/query_range?query={" + a[:len(a)-1] + "}" + lokire_string + "&start=" + strconv.Itoa(int(latencycollectiontimeformat)) + "&end=" + strconv.Itoa(int(collectionscopetime)) + "&limit=" + recordslimit
		log.Println("访问的loki接口是", url)
	} else {
		url = "http://" + lokiipport + "/loki/api/v1/query_range?query={" + a[:len(a)-1] + "}&start=" + strconv.Itoa(int(latencycollectiontimeformat)) + "&end=" + strconv.Itoa(int(collectionscopetime)) + "&limit=" + recordslimit
		log.Println("访问的loki接口是", url)
	}
	// url := "http://" + lokiipport + "/loki/api/v1/query_range?query={job=\"" + "chaos" + "\"}&start=1672816117813000000&end=1672902517813000000&limit=8000"
	client := &http.Client{}
	var lokiquery_range LokiQuery_range
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("wraps: create ", url, "request  error.")
		return loglists, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Get", url, " err.")
		return loglists, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("read resp error log read err.")
		return loglists, err
	}

	json.Unmarshal(body, &lokiquery_range)
	for _, i := range lokiquery_range.Data.Result {
		for _, j := range i.Values {
			loglists = append(loglists, j[1])
		}
	}
	fmt.Println(len(loglists))
	return loglists, nil
}

// func

func Getlogandurlfromloki(lokiipport string, label_list map[string]interface{}, timenow time.Time, latencycollectionseconds string, collectionscopeseconds string, recordslimit string, lokire_string string) (loglists []string, url string, err error) {
	latencyseconds, _ := time.ParseDuration("-" + latencycollectionseconds + "s")
	latencycollectiontime := timenow.Add(1 * latencyseconds)
	latencycollectiontimeformat := ((latencycollectiontime.UnixNano()) / 1000000000) * 1000000000

	scopeseconds, _ := time.ParseDuration(collectionscopeseconds + "s")
	collectionscopetime := ((latencycollectiontime.Add(1*scopeseconds).UnixNano())/1000000000)*1000000000 - 1

	log.Println("开始时间：", latencycollectiontimeformat, "结束时间：", collectionscopetime)

	// fmt.Println("latencycollectiontime:", latencycollectiontimeformat)
	// fmt.Println("collectionscopetime:", collectionscopetime)
	// fmt.Println("nowtime:", time.Now().UnixNano())

	// fmt.Println("environment:", label_list["environment"])
	// fmt.Println("job:", label_list["job"])
	// fmt.Println("filename:", label_list["filename"])
	// fmt.Println("latencycollectiontime:", latencycollectiontimeformat)

	var a string
	var keys []string
	for key := range label_list {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		a = a + key + "=\"" + label_list[key].(string) + "\","
		log.Println("访问的loki接口是", a)
	}

	if lokire_string != "" {
		url = "http://" + lokiipport + "/loki/api/v1/query_range?query={" + a[:len(a)-1] + "}" + lokire_string + "&start=" + strconv.Itoa(int(latencycollectiontimeformat)) + "&end=" + strconv.Itoa(int(collectionscopetime)) + "&limit=" + recordslimit
		log.Println("访问的loki接口是", url)
	} else {
		url = "http://" + lokiipport + "/loki/api/v1/query_range?query={" + a[:len(a)-1] + "}&start=" + strconv.Itoa(int(latencycollectiontimeformat)) + "&end=" + strconv.Itoa(int(collectionscopetime)) + "&limit=" + recordslimit
		log.Println("访问的loki接口是", url)
	}
	// url := "http://" + lokiipport + "/loki/api/v1/query_range?query={job=\"" + "chaos" + "\"}&start=1672816117813000000&end=1672902517813000000&limit=8000"
	client := &http.Client{}
	var lokiquery_range LokiQuery_range
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("wraps: create ", url, "request  error.")
		return loglists, strings.Split(url, "&")[0], err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Get", url, " err.")
		return loglists, strings.Split(url, "&")[0], err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("read resp error log read err.")
		return loglists, url, err
	}

	json.Unmarshal(body, &lokiquery_range)
	for _, i := range lokiquery_range.Data.Result {
		for _, j := range i.Values {
			loglists = append(loglists, j[1])
		}
	}
	fmt.Println(len(loglists))
	return loglists, strings.Split(url, "&")[0], nil
}
