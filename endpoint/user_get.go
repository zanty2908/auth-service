package endpoint

import (
	"context"
	"net/http"
	"time"

	db "auth-service/db/gen"
	"auth-service/utils"
)

type GetUserByIdParams struct {
	Id string
}

func (s *Module) decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	payload, err := s.tokenMaker.GetTokenPayload(r)
	if err != nil {
		return nil, utils.ErrorBadRequest
	}
	return GetUserByIdParams{Id: payload.Subject}, nil
}

func (s *Module) UserGet(c context.Context, request interface{}) (interface{}, error) {
	req := request.(*GetUserByIdParams)

	user, err := s.repo.GetUser(c, req.Id)
	if err != nil {
		return nil, err
	}

	return mapUserResponse(user), nil
}

type UserResponse struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	FullName  string     `json:"fullName"`
	Phone     string     `json:"phone"`
	Country   string     `json:"country,omitempty"`
	Email     *string    `json:"email,omitempty"`
	LastSign  *time.Time `json:"lastSign,omitempty"`
	Role      string     `json:"role,omitempty"`
	Aud       string     `json:"aud,omitempty"`
}

func mapUserResponse(item *db.User) *UserResponse {
	if item == nil {
		return nil
	}
	return &UserResponse{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		FullName:  item.FullName,
		Phone:     item.Phone,
		Country:   item.Country,
		Email:     item.Email,
		LastSign:  item.LastSign,
		Role:      item.Role,
		Aud:       item.Aud,
	}
}
