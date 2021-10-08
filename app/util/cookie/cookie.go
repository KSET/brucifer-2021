package cookie

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/session"
	"brucosijada.kset.org/app/util/auth"
)

type cookieUtil struct {
}

func Cookie() cookieUtil {
	return cookieUtil{}
}

func (c cookieUtil) SetAuthCookie(user *models.User, ctx *fiber.Ctx) {
	store, err := session.SessionProvider().Client().Get(ctx)
	if err != nil {
		panic(err)
	}

	store.Set(auth.AuthProvider().UserIdKey(), user.ID)
	if err := store.Save(); err != nil {
		panic(err)
	}
}

func (c cookieUtil) RemoveAuthCookie(ctx *fiber.Ctx) bool {
	auth := auth.AuthProvider()
	session := session.SessionProvider()

	if !auth.IsAuthenticated(ctx) {
		return false
	}

	store, _ := session.Client().Get(ctx)
	store.Destroy()

	ctx.ClearCookie(session.CookieName())

	return true
}
