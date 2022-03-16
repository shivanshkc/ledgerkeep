package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/configs"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// Auth middleware verifies if the x-password header in the request matches the one in configs.
func Auth(next http.Handler) http.Handler {
	conf := configs.Get()
	log := logger.Get()

	// Calculating the expected Username and Password.
	expectedUsernameHash := sha256.Sum256([]byte(conf.Auth.Username))
	expectedPasswordHash := sha256.Sum256([]byte(conf.Auth.Password))

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		// Retrieving user provided username and password.
		username, password, ok := request.BasicAuth()
		if !ok {
			httputils.WriteErrAndLog(ctx, writer, errutils.Unauthorized(), log)
			return
		}

		// Hashing the provided username and password for comparison with the expected ones.
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		// Comparing user provided credentials with the expected ones.
		usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
		passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

		// If they don't match, it's 401.
		if !usernameMatch || !passwordMatch {
			httputils.WriteErrAndLog(ctx, writer, errutils.Unauthorized(), log)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
