package repo

import (
	"context"
	"fmt"

	"auth-service/config"
	db "auth-service/db/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repo interface {
	db.Querier
	ExecTx(ctx context.Context, fn func(*db.Queries) error) error
	CreateUserTx(ctx context.Context, otpAuthId int32, arg *db.CreateUserParams) (*db.User, error)
	Caching
}

type RepoImpl struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
	*db.Queries
}

func NewRepo(config *config.Config, sqldb *pgxpool.Pool) Repo {
	rdb := NewRedisClient(config)

	return &RepoImpl{
		DB:      sqldb,
		Redis:   rdb,
		Queries: db.New(sqldb),
	}
}

// ExecTx executes a function within a database transaction
func (s *RepoImpl) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
