package sever

import (
	"douyinshibie/api"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strings"
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
		var finalUrl, title string
		if strings.Contains(data["url"], "tiktok.com") {
			f, t, err := api.GetTikTokInfo(data["url"])
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": " 视频链接解析失败，请检查视频链接是否正确",
					"error":   err.Error(),
				})
			}
			finalUrl = f
			title = t
		} else if strings.Contains(data["url"], "douyin.com") {
			f, t, err := api.GetDouYinInfo(data["url"])
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": " 视频链接解析失败，请检查视频链接是否正确",
					"error":   err.Error(),
				})
			}
			finalUrl = f
			title = t
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "未知的视频链接，请检查视频链接是否正确",
				"error":   "url is not tiktok or douyin",
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
			"content":  d,
		})
	}
}

// isAllowedVideoExtension 检查文件扩展名是否为允许的视频扩展名
func isAllowedVideoExtension(ext string) bool {
	allowedExtensions := []string{".mp4", ".avi", ".mov", ".wmv", ".mkv"}
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func VideoFileProcessing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "文件上传失败，请重试",
				"error":   err.Error(),
			})
		}
		if !isAllowedVideoExtension(filepath.Ext(file.Filename)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "文件格式错误，请上传mp4、avi、mov、wmv、mkv格式的视频",
				"error":   "Invalid file extension: " + filepath.Ext(file.Filename),
			})
		}

		model := c.FormValue("model")
		if len(model) == 0 {
			model = "gemini"
		}
		fileStream, err := file.Open()
		if err != nil {
			return err
		}
		defer fileStream.Close()
		err, s := api.VideoFileSlice(fileStream, model)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": " 视频解析出现错误，请重试",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"content": s,
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
		var finalUrl, title string
		if strings.Contains(data["url"], "tiktok.com") {
			f, t, err := api.GetTikTokInfo(data["url"])
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": " 视频链接解析失败，请检查视频链接是否正确",
					"error":   err.Error(),
				})
			}
			finalUrl = f
			title = t
		} else if strings.Contains(data["url"], "douyin.com") {
			f, t, err := api.GetDouYinInfo(data["url"])
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": " 视频链接解析失败，请检查视频链接是否正确",
					"error":   err.Error(),
				})
			}
			finalUrl = f
			title = t
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "未知的视频链接，请检查视频链接是否正确",
				"error":   "url is not tiktok or douyin",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "success",
			"finalUrl": finalUrl,
			"title":    title,
		})
	}
}
