package methods

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

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

func MustRegisterOnce(config *viper.Viper, db *sql.DB) http.Handler {
	// defer db.Close()
	fmt.Println("postgres.metrics is ", config.GetStringSlice("postgres.server")[0])
	for _, metricaa := range config.GetStringMap("postgres.metrics.items") {
		name := metricaa.(map[string]interface{})["name"].(string)
		help := metricaa.(map[string]interface{})["help"].(string)
		sql_content := metricaa.(map[string]interface{})["sql_content"].(string)
		fmt.Println("metricaa is ", metricaa)
		// newGaugeVecinstance := newGaugeVec(metricaa.(map[string]interface{})["name"].(string), metricaa.(map[string]interface{})["help"].(string))
		newGaugeVecinstance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
		go func() error {
			for {
				err := PgQueryAndUpdateMetric(newGaugeVecinstance, db, sql_content)
				if err != nil {
					return err
				}
				
				duration, err := time.ParseDuration(config.GetString("postgres.metrics.cycle"))
				// fmt.Println(config.GetString("postgres.metrics.cycle"))
				if err != nil {
					panic(err)

				}
				time.Sleep(duration)
			}
		}()
		prometheus.MustRegister(newGaugeVecinstance)
	}
	return promhttp.Handler()
}

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
