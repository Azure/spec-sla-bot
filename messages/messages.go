package messages

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

// InfraDB is a connection to the infra monitor database to be used
// throughout the application.
var InfraDB *sql.DB

func init() {
	conn, err := sql.Open("mssql", os.Getenv("CUSTOMCONNSTR_INFRA_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	InfraDB = conn
}
