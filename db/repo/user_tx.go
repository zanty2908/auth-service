package repo

import (
	db "auth-service/db/gen"
	"context"

	"github.com/rs/zerolog/log"
)

func (s *RepoImpl) CreateUserTx(ctx context.Context, otpAuthId int32, arg *db.CreateUserParams) (*db.User, error) {
	var user *db.User
	err := s.ExecTx(ctx, func(q *db.Queries) error {
		var err error

		err = s.DeleteOTPAuthByID(ctx, otpAuthId)
		if err != nil {
			log.Err(err).Msg("CreateUserTx: delete otp auth error")
			return err
		}

		user, err = q.CreateUser(ctx, arg)
		if err != nil {
			log.Err(err).Msg("CreateUserTx: create user error")
			return err
		}

		return nil
	})

	return user, err
}
