package dbkit

import (
	"database/sql"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/lib/pq"
	"github.com/wasay-usmani/go-boilerplate/pkg/logkit"
)

// The LoadPostgresSqlConn function can be used for redshift as well as postgres
func LoadPostgresSqlConn(postgresURL string, debug bool, l logkit.Logger) (*SqlConn, error) {
	sqldb, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect postgres: %w", err)
	}

	err = sqldb.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping postgres: %w", err)
	}

	if debug {
		boil.DebugMode = true
	}

	return &SqlConn{sqldb, l}, nil
}
