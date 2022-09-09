package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go-rest-api-boilerplate/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type postgre struct {
	host   string
	port   string
	dbName string
	user   string
	pass   string
	Conn   *sql.DB
}

func NewPostgreeDb(host, port, dbName, user, pass string) DB {
	return &postgre{
		host:   host,
		port:   port,
		dbName: dbName,
		user:   user,
		pass:   pass,
	}
}

func (d *postgre) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", d.host, d.port, d.user, d.dbName, d.pass)
	return dsn
}

func (d *postgre) Connect() (*sql.DB, error) {
	//db, err := sql.Open("postgres", d.DSN())

	db, err := otelsql.Open("postgres", d.DSN(),
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName(d.dbName))

	if err != nil {
		log.WithError(err).Fatal("unable to open postgres database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("unable to ping connection postgres database")
		return nil, err
	}

	d.Conn = db
	return d.Conn, nil
}

func initMigrator(dbConn *sql.DB) (*migrate.Migrate, error) {
	//driver, err := mysql.WithInstance(dbConn, &mysql.Config{})

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	//TODO
	mountPath := "file://./migrations"
	if strings.HasSuffix(os.Args[0], ".test") {
		pathMigration, err := os.Getwd()
		if err != nil {
			log.WithError(err).Fatal("error get working directory")
		}
		mountPath = "file://" + pathMigration + "/../../migrations"
	}

	log.Infof("source files migration : %s", mountPath)
	return migrate.NewWithDatabaseInstance(
		mountPath,
		config.App.DbName,
		driver,
	)
}

func (d *postgre) MigrateUp() error {
	m, err := initMigrator(d.Conn)
	if err != nil {
		log.WithError(err).Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Errorf("err up migrating database: %s", err.Error())
		return err
	}

	log.Infoln("migrate database has been successfully")
	return nil
}

func (d *postgre) MigrateDown() error {
	m, err := initMigrator(d.Conn)
	if err != nil {
		log.WithError(err).Fatal(err.Error())
	}

	err = m.Down()
	if err != nil {
		log.WithError(err).Errorf("err down migrating database: %s", err.Error())
		return err
	}

	log.Info("migrate down database has been successfully")
	return nil
}
