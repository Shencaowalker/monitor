package methods

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// func CreateconnectDB(config *viper.Viper) *sql.DB {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		config.GetString("logtopostgresql.postgresqlip"), config.GetString("logtopostgresql.postgresqlport"), config.GetString("logtopostgresql.postgresqluser"), config.GetString("logtopostgresql.postgresqlpass"), config.GetString("logtopostgresql.postgresqldb"))
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

func InitDB(config *viper.Viper,project string) *sqlx.DB {
	dsn := "host=" + config.GetString(project+".postgresqlip") + " port=" + config.GetString(project+".postgresqlport") + " user=" + config.GetString(project+".postgresqluser") + " password=" + config.GetString(project+".postgresqlpass") + " dbname=" + config.GetString(project+".postgresqldb") + " sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		panic(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(5)
	return db
}
func InsertintoDB(db *sqlx.DB, config *viper.Viper, project string,values [][]interface{}) {
	defer db.Close()
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("Beginx error:", err)
		panic(err)
	}
	stmt, err := tx.Preparex(db.Rebind(config.GetString(project+".sqlmod")))
	if err != nil {
		fmt.Println("Prepare error:", err)
		panic(err)
	}
	for _, value := range values {
		_, err = stmt.Exec(value...)
		if err != nil {
			fmt.Println("Exec error:", err)
			panic(err)
		}
	}
	err = stmt.Close()
	if err != nil {
		fmt.Println("stmt close error:", err)
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("commit error:", err)
		panic(err)
	}
	fmt.Println("insert into db is seccess!")
}
