package migrations

import (
	"log/slog"

	"github.com/DerKnerd/gorp/v4"
)

func Migrate202609160005convertTimestampToDate(executor gorp.SqlExecutor) {
	if !CheckMigration(executor, "Migrate202609160005convertTimestampToDate") {
		slog.Info("Convert timestamp column to date")
		_, err := executor.Exec(`
alter table analytics
    alter column timestamp type date using timestamp::date;`)
		if err != nil {
			panic(err)
		}

		err = MarkMigration(executor, "Migrate202609160005convertTimestampToDate")
		if err != nil {
			panic(err)
		}
	}
}
