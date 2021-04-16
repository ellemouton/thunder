package server

import "database/sql"

type Backends interface {
        GetDB() *sql.DB
}
