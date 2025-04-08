package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthProtected(db *gorm.DB, redis *redis.Client) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 1. 获取并验证Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			log.Warnf("empty authorization header")

			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}
		// 2. 解析Bearer Token
		// Bearer ajidosjdawsfqwoi23142
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf("invaild token parts")

			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		// 3. 解析Token
		tokenStr := tokenParts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			log.Warnf("invaild token")

			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		userId := uint(token.Claims.(jwt.MapClaims)["id"].(float64))
		role := token.Claims.(jwt.MapClaims)["role"].(string)

		// 4. 尝试从Redis获取用户会话
		session, err := utils.GetUserSession(redis, ctx.Context(), userId)
		if err == nil && len(session) > 0 {
			if session["token"] == tokenStr {
				// 设置用户信息到上下文
				ctx.Locals("userId", userId)
				ctx.Locals("userRole", role)
				return ctx.Next()
			}
		}

		// 5. 如果Redis中没有会话，尝试从数据库获取用户信息
		var user models.User
		if err := db.Model(&models.User{}).Where("id = ?", userId).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("user not found in the db")

			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		// 6. 将用户会话存入Redis
		err = utils.SetUserSession(redis, ctx.Context(), userId, tokenStr, role)
		if err != nil {
			log.Warnf("failed to set user session in redis: %v", err)
		}

		// 7. 设置会话过期时间
		key := fmt.Sprintf("user:%d:session", userId)
		utils.SetExpiration(redis, ctx.Context(), key, time.Hour*24)

		// 8. 检查是否是管理员路由
		if strings.Contains(ctx.Path(), "/statistics") && user.Role != models.Manager {
			return ctx.Status(fiber.StatusForbidden).JSON(&fiber.Map{
				"status":  "fail",
				"message": "需要管理员权限",
			})
		}

		// 9. 设置用户信息到上下文
		ctx.Locals("userId", userId)
		ctx.Locals("userRole", role)
		return ctx.Next()
	}
}
