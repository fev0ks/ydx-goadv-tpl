package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"log"
)

const (
	migrationsDir = "./migrations/postgres/"
)

type DBProvider interface {
	HealthCheck(ctx context.Context) error
	GetConnection(ctx context.Context) (*pgxpool.Conn, error)
}

type pgProvider struct {
	conn *pgxpool.Pool
}

func NewPgProvider(ctx context.Context, appConfig *config.AppConfig) (DBProvider, error) {
	if appConfig == nil {
		log.Println("Postgres DB config is empty")
		return nil, errors.New("failed to init pg repository: appConfig is nil")
	}
	pg := &pgProvider{}
	err := pg.connect(ctx, appConfig.DBConnection, appConfig.DBMaxConnections)
	if err != nil {
		return nil, err
	}
	err = pg.migrationUp(ctx)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *pgProvider) connect(ctx context.Context, connString string, maxConns int32) error {
	if connString == "" {
		log.Println("Postgres DB config is empty")
		return errors.New("failed to init pg repository: dbConfig is empty")
	}
	log.Printf("Trying to connect: %s", connString)
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}
	dbConfig.MaxConns = maxConns
	conn, err := pgxpool.ConnectConfig(ctx, dbConfig)
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
	acquireConn, err := p.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	log.Println("Migrations are started")
	migrator, err := migrate.NewMigrator(ctx, acquireConn.Conn(), "schema_version")
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}
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

func (p *pgProvider) GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	acquireConn, err := p.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return acquireConn, err
}

func (p *pgProvider) HealthCheck(ctx context.Context) error {
	conn, err := p.GetConnection(ctx)
	if err != nil {
		log.Printf("failed to check connection to Postgres DB: %v", err)
		return err
	}
	defer conn.Release()
	err = conn.Conn().Ping(ctx)
	if err != nil {
		log.Printf("failed to check connection to Postgres DB: %v", err)
		return err
	}
	log.Println("Postgres DB connection is active")
	return nil
}
