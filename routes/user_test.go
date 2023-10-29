package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "auth-service/db/gen"
	mockrepo "auth-service/db/mock"
	"auth-service/endpoint"
	"auth-service/utils"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) *db.User {
	random := utils.NewUtilRandom()
	hashedPassword, err := utils.HashPassword(random.RandomString(6))
	require.NoError(t, err)
	birDay := random.RandomBirthday()
	email := random.RandomEmail()
	customer := &db.User{
		ID:       uuid.New().String(),
		FullName: random.RandomName(),
		Email:    &email,
		Birthday: &birDay,
		Phone:    random.RandomPhone(),
		Password: &hashedPassword,
		Address:  random.RandomStringP(30),
	}
	return customer
}

func TestSendOTPAPI(t *testing.T) {

	testCases := []struct {
		name          string
		phone         string
		buildStubs    func(store *mockrepo.MockRepo)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "New User",
			phone: "+84909000999",
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					GetOTPAuth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, pgx.ErrNoRows)

				store.EXPECT().
					CreateOTPAuth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Failed Resend User",
			phone: "+84909000999",
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					GetOTPAuth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&db.OtpAuthentication{
						ID:       1,
						Phone:    "+84909000999",
						Otp:      "1111",
						ResendAt: time.Now().Add(time.Minute),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// {
		// 	name:       "NotFound",
		// 	customerId: customer.ID,
		// 	buildStubs: func(store *mockrepo.MockRepo) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Eq(customer.ID)).
		// 			Times(1).
		// 			Return(nil, sql.ErrNoRows)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusOK, recorder.Code)
		// 		result := requireBodyMatchError(t, recorder.Body)
		// 		require.Nil(t, result.Data)
		// 		require.NotEmpty(t, result.Message)
		// 	},
		// },
		// {
		// 	name:       "InternalError",
		// 	customerId: customer.ID,
		// 	buildStubs: func(store *mockrepo.MockRepo) {
		// 		store.EXPECT().
		// 			GetUser(gomock.Any(), gomock.Eq(customer.ID)).
		// 			Times(1).
		// 			Return(nil, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusOK, recorder.Code)
		// 		result := requireBodyMatchError(t, recorder.Body)
		// 		require.Nil(t, result.Data)
		// 		require.NotEmpty(t, result.Message)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			recorder := httptest.NewRecorder()

			httpRouter := newTestHttpHandler(t, repo, nil)
			// Marshal body data to JSON
			data, err := json.Marshal(map[string]string{"phone": tc.phone})
			require.NoError(t, err)

			url := "/otp"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			httpRouter.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUserByIDAPI(t *testing.T) {
	customer := randomUser(t)

	testCases := []struct {
		name          string
		customerId    string
		buildStubs    func(store *mockrepo.MockRepo)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			customerId: customer.ID,
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				result := requireBodyMatchUser(t, recorder.Body)
				requireUser(t, customer, result)
			},
		},
		{
			name:       "NotFound",
			customerId: customer.ID,
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				result := requireBodyMatchError(t, recorder.Body)
				require.Nil(t, result.Data)
				require.NotEmpty(t, result.Meta.Message)
			},
		},
		{
			name:       "InternalError",
			customerId: customer.ID,
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				result := requireBodyMatchError(t, recorder.Body)
				require.Nil(t, result.Data)
				require.NotEmpty(t, result.Meta.Message)
			},
		},
	}

	tokenKey, authorizationHeader := newAuthToken(t)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			recorder := httptest.NewRecorder()

			httpRouter := newTestHttpHandler(t, repo, &tokenKey)

			url := fmt.Sprintf(`/admin/customer/%s`, tc.customerId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			request.Header.Set("authorization", authorizationHeader)
			require.NoError(t, err)

			httpRouter.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListUserAPI(t *testing.T) {
	n := 10
	customers := make([]*db.User, n)
	for i := 0; i < n; i++ {
		item := randomUser(t)
		customers[i] = item
	}

	type Query struct {
		page    int
		perPage int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockrepo.MockRepo)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: Query{page: 1, perPage: 10},
			buildStubs: func(store *mockrepo.MockRepo) {
				param := &db.ListUsersParams{Limit: int32(n), Offset: 0}
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(param)).
					Times(1).
					Return(customers, nil)
				store.EXPECT().
					CountUsers(gomock.Any()).
					Times(1).
					Return(int64(20), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res struct {
					endpoint.PagingResponse
					Data []db.User `json:"data"`
				}
				err := json.NewDecoder(recorder.Body).Decode(&res)
				require.NoError(t, err)
				require.NotEmpty(t, res.Data)
				require.Len(t, res.Data, 10)
				require.Equal(t, res.TotalResult, 20)
			},
		},
		{
			name:  "InternalError",
			query: Query{page: 1, perPage: 5},
			buildStubs: func(store *mockrepo.MockRepo) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, utils.ErrorFailed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res endpoint.Response
				err := json.NewDecoder(recorder.Body).Decode(&res)
				require.NoError(t, err)
				require.Empty(t, res.Data)
				require.Equal(t, http.StatusInternalServerError, res.Meta.Code)
			},
		},
	}

	tokenKey, authorizationHeader := newAuthToken(t)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf(`/admin/customer?page=%v&perPage=%v`, tc.query.page, tc.query.perPage)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			request.Header.Set("authorization", authorizationHeader)
			require.NoError(t, err)
			httpRouter := newTestHttpHandler(t, repo, &tokenKey)

			httpRouter.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer) *db.User {
	var res struct {
		StatusCode int      `json:"code"`
		Message    string   `json:"message"`
		Data       *db.User `json:"data,omitempty"`
	}
	err := json.NewDecoder(body).Decode(&res)
	require.NoError(t, err)

	return res.Data
}

func requireUser(t *testing.T, expected, result *db.User) {
	require.NotNil(t, result)
	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.FullName, result.FullName)
	require.Equal(t, expected.Email, result.Email)
}

// func TestLoginUserAPI(t *testing.T) {
// 	random := utils.NewUtilRandom()
// 	pwd := random.RandomString(12)
// 	customer := randomUserWithPwd(t, pwd)

// 	type Body struct {
// 		Email    string
// 		Password string
// 		ByPhone  bool `json:"byPhone"`
// 	}

// 	testCases := []struct {
// 		name          string
// 		body          Body
// 		buildStubs    func(store *mockrepo.MockRepo)
// 		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: Body{Email: customer.Email, Password: pwd},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 				store.EXPECT().
// 					GetUserByEmail(gomock.Any(), gomock.Eq(customer.Email)).
// 					Times(1).
// 					Return(customer, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res struct {
// 					endpoint.Response
// 					Data *endpoint.LoginUserResponse `json:"data"`
// 				}
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, res.Data)
// 				require.Empty(t, res.Message)
// 				require.Equal(t, res.Data.Email, customer.Email)
// 				require.NotEmpty(t, res.Data.Token)
// 				require.NotEmpty(t, res.Data.RefreshToken)
// 			},
// 		},
// 		{
// 			name: "OK Phone",
// 			body: Body{Email: customer.Phone, Password: pwd, ByPhone: true},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 				store.EXPECT().
// 					GetUserByPhone(gomock.Any(), gomock.Eq(customer.Phone)).
// 					Times(1).
// 					Return(customer, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res struct {
// 					endpoint.Response
// 					Data *endpoint.LoginUserResponse `json:"data"`
// 				}
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, res.Data)
// 				require.Empty(t, res.Message)
// 				require.Equal(t, res.Data.Email, customer.Email)
// 				require.NotEmpty(t, res.Data.Token)
// 				require.NotEmpty(t, res.Data.RefreshToken)
// 			},
// 		},
// 		{
// 			name: "ValidateEmail",
// 			body: Body{Email: "eweqewqeq", Password: "313123"},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res endpoint.Response
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.Nil(t, res.Data)
// 				require.Equal(t, http.StatusBadRequest, res.StatusCode)
// 				require.NotEmpty(t, res.Message)
// 			},
// 		},
// 		{
// 			name: "NotFound",
// 			body: Body{Email: "cxcsadz3dsadsad2131@gmail.com", Password: "313123"},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 				store.EXPECT().
// 					GetUserByEmail(gomock.Any(), gomock.Eq("cxcsadz3dsadsad2131@gmail.com")).
// 					Times(1).
// 					Return(nil, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res endpoint.Response
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.Nil(t, res.Data)
// 				require.Equal(t, http.StatusNotFound, res.StatusCode)
// 				require.NotEmpty(t, res.Message)
// 			},
// 		},
// 		{
// 			name: "WrongPassword",
// 			body: Body{Email: customer.Email, Password: "313123"},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 				store.EXPECT().
// 					GetUserByEmail(gomock.Any(), gomock.Eq(customer.Email)).
// 					Times(1).
// 					Return(customer, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res endpoint.Response
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.Nil(t, res.Data)
// 				require.NotEmpty(t, res.Message)
// 			},
// 		},
// 		{
// 			name: "BlockedUser",
// 			body: Body{Email: customer.Email, Password: customer.Password},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 				blockedUser := customer
// 				curTime := time.Now()
// 				blockedUser.DeletedAt = &curTime
// 				store.EXPECT().
// 					GetUserByEmail(gomock.Any(), gomock.Eq(customer.Email)).
// 					Times(1).
// 					Return(blockedUser, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res endpoint.Response
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.Nil(t, res.Data)
// 				require.Equal(t, http.StatusLocked, res.StatusCode)
// 				require.NotEmpty(t, res.Message)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := mockrepo.NewMockRepo(ctrl)
// 			tc.buildStubs(repo)

// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/public/login"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			httpRouter := newTestHttpHandler(t, repo, nil)

// 			httpRouter.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

// func TestRenewTokenAPI(t *testing.T) {
// 	random := utils.NewUtilRandom()
// 	pwd := random.RandomString(12)
// 	customer := randomUserWithPwd(t, pwd)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	repo := mockrepo.NewMockRepo(ctrl)
// 	server := newTestServer(t, repo)

// 	var refreshToken string
// 	t.Run("CreateUser", func(t *testing.T) {
// 		repo.EXPECT().
// 			GetUserByEmail(gomock.Any(), gomock.Eq(customer.Email)).
// 			Times(1).
// 			Return(customer, nil)

// 		server.Repo = repo
// 		recorder := httptest.NewRecorder()

// 		// Marshal body data to JSON
// 		jsBody := fmt.Sprintf(`{"email": "%s","password": "%s"}`, customer.Email, pwd)

// 		url := "/public/login"
// 		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsBody)))
// 		require.NoError(t, err)

// 		httpRouter := newTestHttpHandler(t, repo, nil)
// 		httpRouter.ServeHTTP(recorder, request)

// 		var res struct {
// 			endpoint.Response
// 			Data *endpoint.LoginUserResponse `json:"data"`
// 		}

// 		err = json.NewDecoder(recorder.Body).Decode(&res)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, res.Data)
// 		require.NotEmpty(t, res.Data.RefreshToken)
// 		refreshToken = res.Data.RefreshToken
// 	})

// 	type Body struct {
// 		RefreshToken string `json:"refreshToken"`
// 	}
// 	testCases := []struct {
// 		name          string
// 		body          Body
// 		buildStubs    func(store *mockrepo.MockRepo)
// 		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: Body{RefreshToken: refreshToken},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res struct {
// 					endpoint.Response
// 					Data *endpoint.LoginUserResponse `json:"data"`
// 				}
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.NotEmpty(t, res.Data)
// 				require.Empty(t, res.Message)
// 				require.NotEmpty(t, res.Data.Token)
// 				require.NotEmpty(t, res.Data.RefreshToken)
// 			},
// 		},
// 		{
// 			name: "EmptyToken",
// 			body: Body{RefreshToken: ""},
// 			buildStubs: func(store *mockrepo.MockRepo) {
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				var res endpoint.Response
// 				err := json.NewDecoder(recorder.Body).Decode(&res)
// 				require.NoError(t, err)
// 				require.Equal(t, http.StatusBadRequest, res.StatusCode)
// 				require.Nil(t, res.Data)
// 				require.NotEmpty(t, res.Message)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {

// 			tc.buildStubs(repo)
// 			server.Repo = repo

// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/public/renew"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			httpRouter := newTestHttpHandler(t, repo, nil)
// 			httpRouter.ServeHTTP(recorder, request)

// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }
