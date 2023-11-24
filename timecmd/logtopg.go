package timecmd

import (
	"fmt"
	"monitor/methods"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/viper"
)

func Serverlogtopg(config *viper.Viper, project string) {
	exp2 := config.GetStringMap(project + ".listmetrics")
	// fmt.Println("exp2 is:", exp2)
	var values [][]interface{}
	timenow := time.Now()
	for ii, jj := range exp2 {
		loglists := methods.Getlogfromloki(config.GetString(project+".lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), timenow, jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string), jj.(map[string]interface{})["recordslimit"].(string), (jj.(map[string]interface{})["lokire"]).(string))
		for _, j := range loglists {
			// fmt.Println(j)
			raw := methods.SplitoneLineforreFilter(j, jj.(map[string]interface{})["regex"].(string))
			if len(raw) != 5 {
				fmt.Println("匹配失败")
			} else if !utf8.ValidString(raw[4].(string)) {
				fmt.Println("字段中包含特殊字符，忽略")
			} else {
				raw = append(raw[1:], ii)
				values = append(values, raw)
			}
		}
	}
	// fmt.Println(values)
	db := methods.InitDB(config, project)
	methods.InsertintoDB(db, config, project, values)
}

func Servertypelogetopg(config *viper.Viper, logtype string, project string) {
	exp2 := config.GetStringMap(project + ".listmetrics")
	// fmt.Println("exp2 is:", exp2)
	var values [][]interface{}
	timenow := time.Now()
	for ii, jj := range exp2 {
		loglists := methods.Getlogfromloki(config.GetString(project+".lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), timenow, jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string), jj.(map[string]interface{})["recordslimit"].(string), (jj.(map[string]interface{})["lokire"]).(string))
		for _, j := range loglists {
			// fmt.Println(j)
			raw := methods.SplitoneLineforreFilter(j, jj.(map[string]interface{})["regex"].(string))
			if len(raw) != 5  {
				fmt.Println("匹配失败")
				fmt.Println(j)
			} else if !utf8.ValidString(raw[4].(string)) {
				fmt.Println("字段中包含特殊字符，忽略")
			} else if strings.ToLower(raw[2].(string)) == strings.ToLower(logtype) {
				{
					raw = append(raw[1:], ii)
					values = append(values, raw)
				}
			}
		}
	}
	fmt.Println(len(values))
	db := methods.InitDB(config, project)
	methods.InsertintoDB(db, config, project, values)
}
