package configkit

import (
	"fmt"
	"strings"
)

type DB struct {
	SuperUserDatabaseURL string `validate:"required"`
	Schema               string
	SqlWriteURL          string
	SqlReadURL           string
}

// InitDB loads database specific configs into DB struct.
// Specified dbUrlPrefix helps to load correct env var, because
// this method is common for all services.
// e.g. dbUrlPrefix = DISCOUNT returns WriteURL = DISCOUNT_SQL_WRITE_URL.
func InitDB(c *C, dbUrlPrefix string) DB {
	return DB{
		Schema:               c.Viper.GetString("DB_SCHEMA"),
		SuperUserDatabaseURL: c.Viper.GetString("SUPERUSER_DATABASE_URL"),
		SqlWriteURL:          c.Viper.GetString(fmt.Sprintf("%s_SQL_WRITE_URL", strings.ToUpper(dbUrlPrefix))),
		SqlReadURL:           c.Viper.GetString(fmt.Sprintf("%s_SQL_READ_URL", strings.ToUpper(dbUrlPrefix))),
	}
}
