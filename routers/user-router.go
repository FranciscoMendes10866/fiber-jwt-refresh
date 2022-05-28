package routers

import (
	"go-refresh/repositories"

	"github.com/gofiber/fiber/v2"
)

type AppRouterInterface interface {
	SetupUserRouter(app *fiber.App)
}

type appRouter struct{}

func NewAppRouter() AppRouterInterface {
	return &appRouter{}
}

func (*appRouter) SetupUserRouter(app *fiber.App) {
	userRoutes := app.Group("/api/user")

	userRoutes.Post("/signup", repositories.UserRepository.SignUpHandler)
	userRoutes.Post("/signin", repositories.UserRepository.SignInHandler)
	userRoutes.Post("/refresh", repositories.UserRepository.RefreshTokenHandler)
	userRoutes.Post("/logout", repositories.UserRepository.LogOutHandler)
	userRoutes.Get("/", repositories.UserRepository.AuthorizationMiddleware, repositories.UserRepository.GetUserHandler)
}

// This is the AppRouter instance
var Router AppRouterInterface = NewAppRouter()
