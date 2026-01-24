package api

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sealchat/model"
)

const (
	userPreferenceKeyMaxLen    = 64
	userPreferenceValueMaxSize = 4096
)

type UserPreferenceDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func UserPreferencesGet(c *fiber.Ctx) error {
	user := getCurUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "未登录"})
	}
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "缺少 key"})
	}
	if len(key) > userPreferenceKeyMaxLen {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "key 过长"})
	}
	record, err := model.UserPreferenceGet(user.ID, key)
	if err != nil {
		return wrapError(c, err, "获取偏好失败")
	}
	if record == nil {
		return c.JSON(fiber.Map{"key": key, "value": "", "exists": false})
	}
	return c.JSON(fiber.Map{"key": record.PrefKey, "value": record.PrefValue, "exists": true})
}

func UserPreferencesUpsert(c *fiber.Ctx) error {
	user := getCurUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "未登录"})
	}
	var body UserPreferenceDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "请求参数错误"})
	}
	key := strings.TrimSpace(body.Key)
	if key == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "缺少 key"})
	}
	if len(key) > userPreferenceKeyMaxLen {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "key 过长"})
	}
	if len([]byte(body.Value)) > userPreferenceValueMaxSize {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "value 过长"})
	}
	record, err := model.UserPreferenceUpsert(user.ID, key, body.Value)
	if err != nil {
		return wrapError(c, err, "保存偏好失败")
	}
	return c.JSON(fiber.Map{"key": record.PrefKey, "value": record.PrefValue})
}
