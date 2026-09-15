package migrations

import (
	"database/sql"
	"log/slog"

	"github.com/DerKnerd/gorp/v4"
)

func BaseMigration(conn *sql.DB) {
	slog.Info("Create or replace function to setup foreign keys")
	_, err := conn.Exec(`
create or replace function add_foreign_key_if_not_exists(from_table text, from_column text, to_table text, to_column text)
returns void language plpgsql as
$$
declare 
   fk_exists boolean;
begin
    fk_exists := case when exists (select true
	from information_schema.table_constraints tc
		inner join information_schema.constraint_column_usage ccu
			using (constraint_catalog, constraint_schema, constraint_name)
		inner join information_schema.key_column_usage kcu
			using (constraint_catalog, constraint_schema, constraint_name)
	where constraint_type = 'FOREIGN KEY'
	  and ccu.table_name = to_table
	  and ccu.column_name = to_column
	  and tc.table_name = from_table
	  and kcu.column_name = from_column) then true else false end;
	if not fk_exists then
		execute format('alter table %s add constraint %s_%s_%s_%s_fkey foreign key (%s) references %s(%s) on delete cascade', from_table, from_table, from_column, to_table, to_column, from_column, to_table, to_column);
	end if;
end
$$;
`)
	if err != nil {
		panic(err)
	}
}

func CreateAdditionalTables(conn *sql.DB) {
	_, err := conn.Exec(`
create table if not exists configuration
(
    id                        serial primary key,
    current_frontend_theme_id integer null,
    constraint "uniq_configuration_current_frontend_theme_id"
        unique (current_frontend_theme_id),
    constraint "fk_theme_configuration"
        foreign key (current_frontend_theme_id) references theme (id)
);`)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`
create table if not exists file_tag_file
(
    file_tag_id int not null,
    file_id     int not null,
    primary key (file_id, file_tag_id),
    constraint FK_file_tag_file_file_id
        foreign key (file_id) references file (id)
            on delete cascade,
    constraint FK_file_tag_file_file_tag_id
        foreign key (file_tag_id) references file_tag (id)
            on delete cascade
)`)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`
create table if not exists ip_address
(
	address_type text,
	ip_start     inet,
	ip_end       inet,
	country      text,
	city         text
)`)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(`
create table if not exists migration_state
(
    name text primary key
);`)
	if err != nil {
		panic(err)
	}
}

func CheckMigration(executor gorp.SqlExecutor, name string) bool {
	count, _ := executor.SelectInt(`
select count(*)
from migration_state
where name = $1`, name)

	return count == 1
}

func MarkMigration(executor gorp.SqlExecutor, name string) error {
	_, err := executor.Exec(`
insert into migration_state (name) values ($1)`, name)

	return err
}
