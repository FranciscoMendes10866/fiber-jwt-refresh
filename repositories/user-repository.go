package repositories

import (
	"go-refresh/entities"
	"go-refresh/internals"
	"go-refresh/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenRequestBody struct {
	RefreshToken string `json:"refreshToken"`
}

type UserRepositoryInterface interface {
	SignUpHandler(c *fiber.Ctx) error
	SignInHandler(c *fiber.Ctx) error
	RefreshTokenHandler(c *fiber.Ctx) error
	LogOutHandler(c *fiber.Ctx) error
	GetUserHandler(c *fiber.Ctx) error
	AuthorizationMiddleware(c *fiber.Ctx) error
}

type userRepository struct{}

func NewUserRepository() UserRepositoryInterface {
	return &userRepository{}
}

func (*userRepository) SignUpHandler(c *fiber.Ctx) error {
	body := new(entities.User)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var user entities.User
	internals.Database.Where("email = ?", body.Email).First(&user)

	if user.ID != 0 {
		return c.Status(409).SendString("User already exists")
	}

	hashedPassword, err := services.HashService.HashPassword(body.Password)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	body.Password = hashedPassword

	if err := internals.Database.Create(&body).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(201)
}

func (*userRepository) SignInHandler(c *fiber.Ctx) error {
	body := new(entities.User)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var user entities.User
	internals.Database.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		return c.Status(404).SendString("User not found")
	}

	match, err := services.HashService.VerifyPassword(body.Password, user.Password)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if !match {
		return c.Status(401).SendString("Invalid credentials")
	}

	accessToken, err := services.TokenService.GenerateToken(services.TokenClaims{
		UserID: user.ID,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	refreshToken, err := services.TokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := internals.Database.Create(&refreshToken).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"session": map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken.Token,
		},
		"email":  user.Email,
		"userId": user.ID,
	})
}

func (*userRepository) RefreshTokenHandler(c *fiber.Ctx) error {
	body := new(TokenRequestBody)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var refreshToken entities.RefreshToken
	internals.Database.Where("token = ?", body.RefreshToken).First(&refreshToken)

	if refreshToken.ID == 0 {
		return c.Status(404).SendString("Refresh token not found")
	}

	userId := refreshToken.UserID

	isExpired := services.TokenService.IsRefreshTokenExpired(refreshToken.ExpiresAt)
	if isExpired {
		return c.Status(403).SendString("Refresh token has expired")
	}

	if err := internals.Database.Delete(&refreshToken).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	newAccessToken, err := services.TokenService.GenerateToken(services.TokenClaims{
		UserID: userId,
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	newRefreshToken, err := services.TokenService.GenerateRefreshToken(userId)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := internals.Database.Create(&newRefreshToken).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"session": map[string]string{
			"accessToken":  newAccessToken,
			"refreshToken": newRefreshToken.Token,
		},
	})
}

func (*userRepository) LogOutHandler(c *fiber.Ctx) error {
	body := new(TokenRequestBody)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	var refreshToken entities.RefreshToken
	internals.Database.Where("token = ?", body.RefreshToken).First(&refreshToken)

	if refreshToken.ID == 0 {
		return c.Status(404).SendString("Refresh token not found")
	}

	if err := internals.Database.Delete(&refreshToken).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("Successfully logged out")
}

func (*userRepository) GetUserHandler(c *fiber.Ctx) error {
	userContext := services.UserContextService.GetUserContext(c)

	var user entities.User
	internals.Database.Where("id = ?", userContext.UserID).First(&user)

	if user.ID == 0 {
		return c.Status(404).SendString("User not found")
	}

	return c.JSON(user)
}

func (*userRepository) AuthorizationMiddleware(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")

	if authorization == "" {
		return c.SendStatus(403)
	}

	if authorization[:7] != "Bearer " {
		return c.SendStatus(403)
	}

	token := strings.Split(authorization, " ")[1]
	if token == "" {
		return c.SendStatus(401)
	}

	claims, err := services.TokenService.VerifyToken(token)
	if err != nil {
		return c.SendStatus(401)
	}

	services.UserContextService.SetUserContext(c, services.UserContext{
		UserID: claims.UserID,
	})

	return c.Next()
}

// This is the UserRepository instance
var UserRepository UserRepositoryInterface = NewUserRepository()
