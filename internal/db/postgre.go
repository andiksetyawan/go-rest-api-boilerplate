package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type postgre struct {
	host   string
	port   string
	dbName string
	user   string
	pass   string
	conn   *sql.DB
}

func NewPostgreeDb(host, port, dbName, user, pass string) *postgre {
	return &postgre{
		host:   host,
		port:   port,
		dbName: dbName,
		user:   user,
		pass:   pass,
	}
}

func NewPostgreeTestContainerDb() *postgre {
	ctx := context.Background()
	sqlDb := postgre{
		user:   "postgres",
		pass:   "postgres",
		dbName: "test_db",
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       sqlDb.dbName,
			"POSTGRES_USER":     sqlDb.user,
			"POSTGRES_PASSWORD": sqlDb.pass,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithPollInterval(1 * time.Second),
	}

	containers, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.WithError(err).Fatal("failed to create test container, please check 5432 port on this OS")
	}

	port, err := containers.Ports(ctx)
	if err != nil {
		log.WithError(err).Fatal("unable to get host test container")
	}
	sqlDb.host = port["5432/tcp"][0].HostIP
	sqlDb.port = port["5432/tcp"][0].HostPort

	return &sqlDb
}

func (d *postgre) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", d.host, d.port, d.user, d.dbName, d.pass)
	return dsn
}

func (d *postgre) GetConnection() *sql.DB {
	return d.conn
}

func (d *postgre) Connect() *postgre {
	//db, err := sql.Open("postgres", d.DSN())
	db, err := otelsql.Open("postgres", d.DSN(),
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName(d.dbName))

	if err != nil {
		log.WithError(err).Fatal("unable to open postgres database")
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("unable to ping connection postgres database")
		return nil
	}

	d.conn = db
	return d
}

func initMigrator(d *postgre) (*migrate.Migrate, error) {
	//driver, err := mysql.WithInstance(dbConn, &mysql.Config{})

	driver, err := postgres.WithInstance(d.conn, &postgres.Config{})
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
		d.dbName,
		driver,
	)
}

func (d *postgre) MigrateUp() error {
	m, err := initMigrator(d)
	if err != nil {
		log.WithError(err).Fatal(err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Infof("migrate up database: %s", err.Error())
			return nil
		}
		log.Errorf("err up migrating database: %s", err.Error())
		return err
	}

	log.Infoln("migrate database has been successfully")
	return nil
}

func (d *postgre) MigrateDown() error {
	m, err := initMigrator(d)
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
