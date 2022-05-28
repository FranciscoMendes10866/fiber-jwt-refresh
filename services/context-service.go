package services

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type UserContext struct {
	UserID uint
}

const userContextKey = "user"

type UserContextServiceInterface interface {
	SetUserContext(c *fiber.Ctx, data UserContext)
	GetUserContext(c *fiber.Ctx) UserContext
}

type userContextService struct{}

func NewUserContextService() UserContextServiceInterface {
	return &userContextService{}
}

func (*userContextService) SetUserContext(c *fiber.Ctx, data UserContext) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, userContextKey, data)
	c.SetUserContext(ctx)
}

func (*userContextService) GetUserContext(c *fiber.Ctx) UserContext {
	userCtx := c.UserContext()
	return userCtx.Value(userContextKey).(UserContext)
}

// This is the UserContextService instance
var UserContextService UserContextServiceInterface = NewUserContextService()
