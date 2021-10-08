package session

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

type sessionProvider struct {
}

var session_ *session.Store

func SessionProvider() sessionProvider {
	return sessionProvider{}
}

func (p sessionProvider) Client() *session.Store {
	return session_
}

func (p sessionProvider) CookieName() string {
	return "bruc-auth"
}

func (p sessionProvider) Register() error {
	storage := redis.New(
		redis.Config{
			URL: os.Getenv("REDIS_URL"),
		},
	)

	session_ = session.New(
		session.Config{
			Storage: storage,
			KeyLookup: fmt.Sprintf(
				"%s:%s",
				session.SourceCookie,
				p.CookieName(),
			),
		},
	)

	return nil
}
