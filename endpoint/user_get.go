package endpoint

import (
	"context"
	"net/http"
	"time"

	db "auth-service/db/gen"
	"auth-service/utils"

	kithttp "github.com/go-kit/kit/transport/http"
)

func (s *Module) GetUserEndpoint() *kithttp.Server {
	return s.newHttpEndpoint(
		s.getUser,
		s.decodeGetUserRequest,
	)
}

type GetUserByIdParams struct {
	Id string
}

func (s *Module) decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	customerId := r.Header.Get("customerId")
	if customerId == "" {
		return nil, utils.ErrorBadRequest
	}
	return GetUserByIdParams{Id: customerId}, nil
}

func (s *Module) getUser(c context.Context, request interface{}) (interface{}, error) {
	req := request.(GetUserByIdParams)

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
	Country   string     `json:"country"`
	Email     *string    `json:"email"`
	Birthday  *time.Time `json:"birthday"`
	Avatar    *string    `json:"avatar"`
	Address   *string    `json:"address"`
	Gender    *int16     `json:"gender"`
	Status    int16      `json:"status"`
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
		Birthday:  item.Birthday,
		Avatar:    item.Avatar,
		Address:   item.Address,
		Gender:    item.Gender,
		Status:    item.Status,
	}
}
