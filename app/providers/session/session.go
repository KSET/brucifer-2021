package session

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
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
			Storage:   storage,
			KeyLookup: fmt.Sprintf("cookie:%s", p.CookieName()),
			KeyGenerator: func() string {
				id, _ := uuid.NewUUID()

				return fmt.Sprintf("brucifer:session:%s", id)
			},
		},
	)

	return nil
}
