package configkit

import (
	"fmt"
	"strings"
)

type DB struct {
	SuperUserDatabaseURL string `validate:"required"`
	Schema               string
	SQLWriteURL          string `validate:"required"`
	SQLReadURL           string `validate:"required"`
}

// InitDB loads database specific configs into DB struct.
// Specified dbUrlPrefix helps to load correct env var, because
// this method is common for all services.
// e.g. dbUrlPrefix = DISCOUNT returns WriteURL = DISCOUNT_SQL_WRITE_URL.
func InitDB(c *C, dbURLPrefix string) DB {
	return DB{
		Schema:               c.Viper.GetString("DB_SCHEMA"),
		SuperUserDatabaseURL: c.Viper.GetString("SUPERUSER_DATABASE_URL"),
		SQLWriteURL:          c.Viper.GetString(fmt.Sprintf("%s_SQL_WRITE_URL", strings.ToUpper(dbURLPrefix))),
		SQLReadURL:           c.Viper.GetString(fmt.Sprintf("%s_SQL_READ_URL", strings.ToUpper(dbURLPrefix))),
	}
}
