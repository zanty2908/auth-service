package endpoint

import (
	"context"

	db "auth-service/db/gen"

	"github.com/rs/zerolog/log"
)

func (s *Module) UserUpdate(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*db.UpdateUserParams)

	user, err := s.repo.UpdateUser(c, req)
	if err != nil {
		log.Err(err).Msg("Update user: update user error")
		return nil, err
	}

	return mapUserResponse(user), nil
}
