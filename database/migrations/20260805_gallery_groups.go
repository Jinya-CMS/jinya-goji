package migrations

import (
	"database/sql"
	"log/slog"
)

func AddGalleryGroups(conn *sql.DB) {
	slog.Info("Add column group to gallery")
	_, err := conn.Exec(`
alter table gallery
    add column if not exists "group" text`)
	if err != nil {
		panic(err)
	}
}
