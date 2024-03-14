package methods

import (
	"fmt"
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

func Rdb_Collect_Redi(config *viper.Viper) redis.Conn {
	dbHost := config.GetStringSlice("redis.server")[0]
	dbPort := config.GetString("redis.service_port")
	dbPassword := config.GetString("redis.password")
	// dbName := config.GetString("redis.dbspace")
	rs, err := redis.Dial("tcp", dbHost+dbPort, redis.DialPassword(dbPassword))
	if err != nil {
		log.Fatalf("connect redis failed, err:%v", err)
		panic(err)
	}
	return rs
}

func getRedisInfo(rs redis.Conn) {
	defer rs.Close()
	rs.Send("MULTI")
	rs.Send("info")
	info, err := redis.Strings(rs.Do("EXEC"))
	if err != nil {
		log.Fatalf("get redis info failed, err:%v", err)
	}
	// 获取info信息，并且切分用于遍历分析
	infoSlice := strings.Split(info[0], "\r\n")
	
	for _, dbInfo := range infoSlice {
		splitInfo := strings.Split(dbInfo, ":")
		fmt.Println(splitInfo[0])
	}

}
