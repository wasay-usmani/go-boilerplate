package dbkit

import (
	"database/sql"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql/driver"
	"github.com/wasay-usmani/go-boilerplate/pkg/logkit"
)

func LoadMySqlConn(mysqlURL string, debug bool, l logkit.Logger) (*SqlConn, error) {
	sqldb, err := sql.Open("mysql", mysqlURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect mysql: %w", err)
	}

	boil.DebugMode = debug
	err = sqldb.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping mysql: %w", err)
	}

	return &SqlConn{sqldb, l}, nil
}
