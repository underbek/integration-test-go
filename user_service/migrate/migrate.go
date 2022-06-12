package migrate

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations
var Migrations embed.FS

func Migrate(dsn string, path fs.FS) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	goose.SetBaseFS(path)
	return goose.Up(db, "migrations")
}
