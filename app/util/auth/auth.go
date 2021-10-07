package auth

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/hash"
	"brucosijada.kset.org/app/providers/session"
	"brucosijada.kset.org/app/repo"
)

type authUtil struct{}

func AuthProvider() authUtil {
	return authUtil{}
}

func (a authUtil) UserIdKey() string {
	return "userId"
}

func (a authUtil) GetAuthUserId(ctx *fiber.Ctx) (id uint, valid bool) {
	store, err := session.SessionProvider().Client().Get(ctx)
	if err != nil {
		panic(err)
	}

	val := store.Get(a.UserIdKey())

	if val == nil {
		return 0, false
	}

	userID := val.(uint)

	return userID, userID != 0
}

func (a authUtil) IsAuthenticated(ctx *fiber.Ctx) (authenticated bool) {
	_, valid := a.GetAuthUserId(ctx)

	return valid
}

func (a authUtil) ValidateUser(identity string, password string) *models.User {
	user := repo.User().GetByUsername(identity)
	emptyPassword := hash.HashProvider().EmptyPassword()

	if user.Password == "" {
		user.Password = emptyPassword
	}

	passwordValid := hash.HashProvider().VerifyPassword(password, user.Password)

	if passwordValid {
		return &user
	} else {
		return nil
	}
}
