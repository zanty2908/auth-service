package repo

import (
	db "auth-service/db/gen"
	"auth-service/utils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *db.User {
	hashedPassword, err := utils.HashPassword(random.RandomString(6))
	require.NoError(t, err)

	birDay := random.RandomBirthday()
	email := random.RandomEmail()

	arg := &db.CreateUserParams{
		ID:       uuid.New().String(),
		FullName: random.RandomName(),
		Email:    &email,
		Birthday: &birDay,
		Phone:    random.RandomPhone(),
		Password: &hashedPassword,
		Address:  random.RandomStringP(30),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.ID, user.ID)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Address, user.Address)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestTimer(t *testing.T) {
	resendTime := time.Date(2023, 10, 10, 0, 2, 0, 0, time.Local)
	curTime := time.Date(2023, 10, 10, 0, 0, 0, 0, time.Local)

	require.True(t, resendTime.After(curTime))
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Address, user2.Address)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := random.RandomName()
	updatedUser, err := testQueries.UpdateUser(
		context.Background(),
		&db.UpdateUserParams{
			ID:       oldUser.ID,
			FullName: &newFullName,
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.Password, updatedUser.Password)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := random.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(
		context.Background(),
		&db.UpdateUserParams{
			ID:    oldUser.ID,
			Email: &newEmail,
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Password, updatedUser.Password)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := random.RandomString(6)
	newHashedPassword, err := utils.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(
		context.Background(),
		&db.UpdateUserParams{
			ID:       oldUser.ID,
			Password: &newHashedPassword,
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Password, updatedUser.Password)
	require.Equal(t, newHashedPassword, updatedUser.Password)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := random.RandomName()
	newEmail := random.RandomEmail()
	newPassword := random.RandomString(6)
	newHashedPassword, err := utils.HashPassword(newPassword)
	newAvatar := random.RandomString(30)
	newAddress := random.RandomString(30)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(
		context.Background(),
		&db.UpdateUserParams{
			ID:       oldUser.ID,
			FullName: &newFullName,
			Email:    &newEmail,
			Password: &newHashedPassword,
			Avatar:   &newAvatar,
			Address:  &newAddress,
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Password, updatedUser.Password)
	require.Equal(t, newHashedPassword, updatedUser.Password)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Avatar, updatedUser.Avatar)
	require.Equal(t, newAvatar, *updatedUser.Avatar)
	require.NotEqual(t, oldUser.Address, updatedUser.Address)
	require.Equal(t, newAddress, updatedUser.Address)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := &db.ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	result, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Len(t, result, 5)
}
