package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	migrationsPath = "./dao/postgres/migrations"

	DuplicateError = "DB duplicate error"
	NoRowsError    = "DB no rows in resultset"
)

type Postgres struct {
	cfg config.Postgres
	db  *sqlx.DB
}

func NewPostgres(cfg config.Postgres) (*Postgres, error) {
	conn, err := makeConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("makeConn: %s", err.Error())
	}
	db := &Postgres{
		cfg: cfg,
		db:  sqlx.NewDb(conn, "postgres"),
	}
	err = db.makeMigration(conn, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("makeMigration: %s", err.Error())
	}
	return db, nil
}

func makeConn(cfg config.Postgres) (*sql.DB, error) {
	s := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Database, cfg.SSLMode)
	return sql.Open("postgres", s)
}

func (db Postgres) find(dest interface{}, sb squirrel.SelectBuilder) error {
	sqlStatement, args, err := sb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	err = db.db.Select(dest, sqlStatement, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db Postgres) first(dest interface{}, sb squirrel.SelectBuilder) error {
	sqlStatement, args, err := sb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	err = db.db.Get(dest, sqlStatement, args...)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "sql: no rows in result set" {
			return fmt.Errorf(NoRowsError)
		}
		return err
	}
	return nil
}

func (db Postgres) insert(sb squirrel.InsertBuilder, primaryKey ...string) (lastID uint64, err error) {
	sqlStatement, args, err := sb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}
	if len(primaryKey) > 0 {
		sqlStatement = fmt.Sprintf("%s RETURNING %s ", sqlStatement, primaryKey[0])
		err = db.db.QueryRow(sqlStatement, args...).Scan(&lastID)
		if err != nil {
			errMsg := err.Error()
			if len(errMsg) > 50 {
				if errMsg[:50] == "pq: duplicate key value violates unique constraint" {
					return lastID, fmt.Errorf(DuplicateError)
				}
			}
			return 0, err
		}
	} else {
		_, err = db.db.Exec(sqlStatement, args...)
		if err != nil {
			errMsg := err.Error()
			if len(errMsg) > 50 {
				if errMsg[:50] == "pq: duplicate key value violates unique constraint" {
					return lastID, fmt.Errorf(DuplicateError)
				}
			}
			return 0, err
		}
	}
	return lastID, nil
}

func (db Postgres) update(sb squirrel.UpdateBuilder) (err error) {
	sqlStatement, args, err := sb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = db.db.Exec(sqlStatement, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db Postgres) delete(sb squirrel.DeleteBuilder) (err error) {
	sqlStatement, args, err := sb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = db.db.Exec(sqlStatement, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db Postgres) makeMigration(conn *sql.DB, migrationDir string) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{
		DatabaseName: db.cfg.Database,
	})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %s", err.Error())
	}
	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		db.cfg.Database, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %s", err.Error())
	}
	if err = mg.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
