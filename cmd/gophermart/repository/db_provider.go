package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"io/fs"
	"log"
	"os"
	"sort"
)

const (
	migrationsDir = "./cmd/migrations/postgres/"
)

type DbProvider interface {
	HealthCheck(ctx context.Context) error
	GetConnection() *pgx.Conn
}

type pgProvider struct {
	conn *pgx.Conn
}

func NewPgProvider(ctx context.Context, appConfig *config.AppConfig) (DbProvider, error) {
	if appConfig == nil {
		log.Println("Postgres DB config is empty")
		return nil, errors.New("failed to init pg repository: appConfig is nil")
	}
	pg := &pgProvider{}
	err := pg.connect(ctx, appConfig.DbConnection)
	if err != nil {
		return nil, err
	}
	err = pg.migrationUp(ctx)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *pgProvider) connect(ctx context.Context, dbConfig string) error {
	if dbConfig == "" {
		log.Println("Postgres DB config is empty")
		return errors.New("failed to init pg repository: dbConfig is empty")
	}
	log.Printf("Trying to connect: %s", dbConfig)
	conn, err := pgx.Connect(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf("failed to init pg repository: %v", err)
	}
	p.conn = conn
	return nil
}

func (p *pgProvider) migrationUp(ctx context.Context) error {
	if p.conn == nil {
		return errors.New("failed to start db migration: db connection is empty")
	}
	log.Println("Migrations are started")
	migrator, err := migrate.NewMigrator(ctx, p.conn, "schema_version")
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}
	ReadDir(migrationsDir)
	err = migrator.LoadMigrations(migrationsDir)
	if err != nil {
		log.Fatalf("Unable to load migrations: %v\n", err)
	}
	err = migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Unable to migrate: %v\n", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current schema version: %v", err)
	}
	log.Printf("Migration done. Current schema version: %d", ver)
	return nil
}

func ReadDir(dirname string) ([]fs.FileInfo, error) {
	f, err := os.Open("cmd/migrations/postgres/1_add_users_table.sql")
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func (p *pgProvider) GetConnection() *pgx.Conn {
	return p.conn
}

func (p *pgProvider) HealthCheck(ctx context.Context) error {
	err := p.GetConnection().Ping(ctx)
	if err != nil {
		log.Printf("failed to check connection to Postgres DB: %v", err)
		return err
	}
	log.Println("Postgres DB connection is active")
	return nil
}
