package migration

import "database/sql"

func CreateLinksTable(db *sql.DB) {
	db.Exec("CREATE TABLE IF NOT EXISTS links (id bigserial NOT NULL PRIMARY KEY, link text NOT NULL UNIQUE, code text NOT NULL, userid text NOT NULL);")
}
