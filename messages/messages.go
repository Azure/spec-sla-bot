package messages

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var InfraDB *sql.DB

func init() {
	log.Print("made it to init in infra")
	conn, err := sql.Open("mssql", os.Getenv("CUSTOMCONNSTR_INFRA_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	InfraDB = conn
}
