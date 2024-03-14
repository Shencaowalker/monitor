package apis

import (
	_ "github.com/lib/pq"
)

// var (
// 	dbHost     = "10.5.100.3"
// 	dbPort     = "5432"
// 	dbUser     = "postgres"
// 	dbPassword = "postgres"
// 	dbName     = "postgres"
// )

// func main1() {
// 	viper.SetConfigName("metrics")
// 	// 设置配置文件的类型
// 	viper.SetConfigType("yaml")
// 	// 添加配置文件的路径，指定 config 目录下寻找
// 	viper.AddConfigPath("./apis")
// 	// 寻找配置文件并读取
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %w", err))
// 	}
// 	var (
// 		dbHost     = viper.Get("postgres.server")
// 		dbPort     = viper.Get("postgres.service_port")
// 		dbUser     = viper.Get("postgres.username")
// 		dbPassword = viper.Get("postgres.password")
// 		dbName     = "postgres"
// 	)

// 	// Connect to PostgreSQL database
// 	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// for _, metricaaa := range viper.GetStringMap("postgres.metrics") {
// 	// 	fmt.Println(metricaaa)
// 	// }

// 	var metrics_all []*prometheus.GaugeVec
// 	for _, metricaa := range viper.GetStringMap("postgres.metrics") {
// 		// fmt.Println(metricaa.(map[string]interface{})["name"])
// 		newGaugeVecinstance := newGaugeVec(metricaa.(map[string]interface{})["name"].(string), metricaa.(map[string]interface{})["help"].(string), metricaa.(map[string]interface{})["sql_content"].(string))
// 		metrics_all = append(metrics_all, newGaugeVecinstance)
// 	}

// 	// Define and register metrics
// 	// metrics := []*prometheus.GaugeVec{
// 	// 	newGaugeVec("metric_name_1", "Metric 1 description"),
// 	// 	newGaugeVec("metric_name_2", "Metric 2 description"),
// 	// 	// Add more metrics as needed
// 	// }

// 	// Register metrics with Prometheus
// 	for _, metric := range metrics_all {
// 		prometheus.MustRegister(metric)
// 	}

// 	// Start HTTP server to expose metrics
// 	// http.Handle("/metrics", promhttp.Handler())
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":8080", nil))
// 	}()

// 	// Query and update metrics periodically
// 	for {
// 		err := queryAndUpdateMetrics(db, metrics_all)
// 		if err != nil {
// 			log.Println("Error querying database:", err)
// 		}
// 		time.Sleep(1 * time.Minute) // Adjust the interval as needed
// 	}
// }

// func queryAndUpdateMetrics(db *sql.DB, metrics []*prometheus.GaugeVec) error {
// 	// Execute SQL query to fetch data from database
// 	// rows, err := db.Query("select count(*) as label1 from pg_stat_activity")
// 	// for _, newGaugeVecinstance := range metrics{
// 	// rows, err := db.Query(newGaugeVecinstance[])
// 	rows, err := db.Query("show max_connections")
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	// Iterate over rows and update metric values
// 	for rows.Next() {
// 		var aaaa int
// 		// var value1, value2 int
// 		// if err := rows.Scan(&value1, &value2); err != nil {
// 		// 	return err
// 		// }
// 		// Update each metric with fetched values
// 		if err := rows.Scan(&aaaa); err != nil {
// 			return err
// 		}
// 		for _, metric := range metrics {
// 			// metric.With(prometheus.Labels{"label1": fmt.Sprintf("%d", value1), "label2": fmt.Sprintf("%d", value2)}).Set(float64(value1 + value2))
// 			metric.With(prometheus.Labels{"host": "aaaaaa"}).Set(float64(aaaa))
// 		}
// 	}

// 	return nil
// }

// func newGaugeVec(name, help, sql_content string) *prometheus.GaugeVec {
// 	return prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: name,
// 			Help: help,
// 			// Sql_content: sql_content,
// 		},
// 		[]string{"host"},
// 	)
// }

// func queryAndUpdateMetrics(db *sql.DB, metrics []*prometheus.GaugeVec) error {
// 	// Execute SQL query to fetch data from database
// 	// rows, err := db.Query("select count(*) as label1 from pg_stat_activity")
// 	rows, err := db.Query("show max_connections")
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	// Iterate over rows and update metric values
// 	for rows.Next() {
// 		var aaaa int
// 		// var value1, value2 int
// 		// if err := rows.Scan(&value1, &value2); err != nil {
// 		// 	return err
// 		// }
// 		// Update each metric with fetched values
// 		if err := rows.Scan(&aaaa); err != nil {
// 			return err
// 		}
// 		for _, metric := range metrics {
// 			// metric.With(prometheus.Labels{"label1": fmt.Sprintf("%d", value1), "label2": fmt.Sprintf("%d", value2)}).Set(float64(value1 + value2))
// 			metric.With(prometheus.Labels{"host": dbHost}).Set(float64(aaaa))
// 		}
// 	}

// 	return nil
// }

// func newGaugeVec(name, help string) *prometheus.GaugeVec {
// 	return prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: name,
// 			Help: help,
// 		},
// 		[]string{"host"},
// 	)
// }
