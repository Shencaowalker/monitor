package methods

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

// 创建pg连接句柄
func Db_Collect_Postgres(config *viper.Viper) *sql.DB {
	dbHost := config.GetStringSlice("postgres.server")[0]
	dbPort := config.GetString("postgres.service_port")
	dbUser := config.GetString("postgres.username")
	dbPassword := config.GetString("postgres.password")
	dbName := "postgres"
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// func newGaugeVec(name, help string) *prometheus.GaugeVec {
// 	return prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: name,
// 			Help: help,
// 			// Sql_content: sql_content,
// 		},
// 		[]string{"host"},
// 	)
// }

func PgQueryAndUpdateMetric(newGaugeVecinstance prometheus.Gauge, db *sql.DB, sql_content string) error {
	rows, err := db.Query(sql_content)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var value float64
		if err := rows.Scan(&value); err != nil {
			log.Println("Error querying database:", err)
			return err
		}
		fmt.Println("value is ", value)
		newGaugeVecinstance.Set(value)
	}
	return nil
}

func PgQueryAndUpdateMetrics(config *viper.Viper, db *sql.DB) error {

	defer db.Close()
	fmt.Println("postgres.metrics is ", config.GetStringSlice("postgres.server")[0])
	for _, metricaa := range config.GetStringMap("postgres.metrics") {
		fmt.Println("metricaa is ", metricaa)
		newGaugeVecinstance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metricaa.(map[string]interface{})["name"].(string),
			Help: metricaa.(map[string]interface{})["help"].(string),
		})
		rows, err := db.Query(metricaa.(map[string]interface{})["sql_content"].(string))
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var value float64
			if err := rows.Scan(&value); err != nil {
				log.Println("Error querying database:", err)
				return err
			}
			fmt.Println("value is ", value)
			newGaugeVecinstance.Set(value)
		}
	}

	return nil
}
