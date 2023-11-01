package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"auth-service/config"
	"auth-service/db/repo"
	"auth-service/endpoint"
	"auth-service/language"
	"auth-service/server"
	"auth-service/token"
	"auth-service/transport"
	"auth-service/utils"

	"github.com/stretchr/testify/require"
)

var random *utils.UtilRandom

func TestMain(m *testing.M) {
	random = utils.NewUtilRandom()
	os.Exit(m.Run())
}

func requireBodyMatchError(t *testing.T, body *bytes.Buffer) *endpoint.Response {
	var res = new(endpoint.Response)
	err := json.NewDecoder(body).Decode(&res)
	require.NoError(t, err)
	return res
}

func newTestHttpHandler(t *testing.T, repo repo.Repo, tokenKey *string) http.Handler {
	var key string
	if tokenKey != nil {
		key = *tokenKey
	} else {
		key = random.RandomString(32)
	}
	cfgToken := config.Token{
		SecretKey:            key,
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Minute,
	}
	config := &config.Config{
		Auth: config.Auth{Token: cfgToken},
	}

	multiLocalizer := language.LoadAllFileLanguage("../language/localizations_src/")
	tokenMaker, err := token.NewJWTMaker(cfgToken.SecretKey, cfgToken.AccessTokenDuration, cfgToken.RefreshTokenDuration)
	require.NoError(t, err)

	ep := endpoint.NewEndpointModule(config, repo, tokenMaker)

	server, _ := server.NewServer(config, repo, &multiLocalizer)
	trans := transport.NewModule(config, repo, &multiLocalizer, server.TokenMaker, ep)
	return NewHttpRouter(server, trans, ep)
}

func newTestServer(t *testing.T, repo repo.Repo) *server.Server {
	config := &config.Config{
		Auth: config.Auth{
			Token: config.Token{
				SecretKey:            random.RandomString(32),
				AccessTokenDuration:  time.Minute,
				RefreshTokenDuration: time.Minute,
			},
		},
	}
	multiLocalizer := language.LoadAllFileLanguage("../language/localizations_src/")
	server, err := server.NewServer(config, repo, &multiLocalizer)
	require.NoError(t, err)

	return server
}

func newAuthToken(t *testing.T) (tokenKey, authorizationHeader string) {
	tokenKey = random.RandomString(32)
	tokenMaker, err := token.NewJWTMaker(tokenKey, time.Minute, time.Hour)
	require.NoError(t, err)

	accressToken, _, _, _, err := tokenMaker.CreateTokenPair(
		token.USER_APP,
		token.ANDROID,
		random.RandomString(12),
		random.RandomString(10),
		random.RandomEmail(),
		random.RandomPhone(),
		make(map[string]interface{}, 0),
	)
	require.NoError(t, err)
	authorizationHeader = fmt.Sprintf("Bearer %s", accressToken)
	return
}
