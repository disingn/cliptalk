package sever

import (
	"douyinshibie/api"
	"github.com/gofiber/fiber/v2"
)

func VideoProcessing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "error parsing JSON",
				"error":   err.Error(),
			})
		}
		if len(data["url"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "error parsing JSON",
				"error":   "url is empty",
			})
		}
		if len(data["model"]) == 0 {
			data["model"] = "gemini"
		}
		videoIdOrLink := api.ProcessUserInput(data["url"])
		var videoId string
		if videoIdOrLink != "" {
			if api.IsNumeric(videoIdOrLink) {
				videoId = videoIdOrLink
			} else {
				videoId = api.ExtractVideoId(videoIdOrLink)
			}
		}
		if len(videoId) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": " 视频解释失败，请检查视频链接是否正确",
				"error":   "videoId is not found",
			})

		}
		finalUrl, title, err := api.GetVideoInfo(videoId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": " 解析失败，请检查视频链接是否正确",
				"error":   err.Error(),
			})
		}
		err, d := api.VideoSlice(finalUrl, data["model"])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "出现未知的错误，请重试",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "success",
			"title":    title,
			"finalUrl": finalUrl,
			"desc":     d,
		})
	}
}

// RemoveWatermark 去除水印
func RemoveWatermark() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "error parsing JSON",
				"error":   err.Error(),
			})
		}
		if len(data["url"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "error parsing JSON",
				"error":   "url is empty",
			})
		}
		videoIdOrLink := api.ProcessUserInput(data["url"])
		var videoId string
		if videoIdOrLink != "" {
			if api.IsNumeric(videoIdOrLink) {
				videoId = videoIdOrLink
			} else {
				videoId = api.ExtractVideoId(videoIdOrLink)
			}
		}
		if len(videoId) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": " 视频解释失败，请检查视频链接是否正确",
				"error":   "videoId is not found",
			})

		}
		finalUrl, title, err := api.GetVideoInfo(videoId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": " 解析失败，请检查视频链接是否正确",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "success",
			"finalUrl": finalUrl,
			"title":    title,
		})
	}
}
