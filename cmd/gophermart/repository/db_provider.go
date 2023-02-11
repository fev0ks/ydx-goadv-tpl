package repository

import "context"

type DbProvider interface {
	HealthCheck(ctx context.Context) error
}

type pgProvider struct {
}

func NewPgProvider() DbProvider {
	return &pgProvider{}
}

func (p pgProvider) HealthCheck(ctx context.Context) error {
	//TODO implement me
	return nil
}
