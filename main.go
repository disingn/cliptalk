package main

import (
	"douyinshibie/cfg"
	"douyinshibie/sever"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	cfg.ConfigInit()
}

func main() {
	r := fiber.New(fiber.Config{
		BodyLimit: cfg.Config.Sever.MaxFileSize * 1024 * 1024,
	})
	r.Use(cors.New(cors.ConfigDefault))
	r.Use(logger.New(logger.ConfigDefault))
	r.Static("/", "./web/cliptalk/dist")
	r.Post("/video", sever.VideoProcessing())
	r.Post("/remove", sever.RemoveWatermark())
	r.Post("/video-file", sever.VideoFileProcessing())
	r.Listen(":" + cfg.Config.Sever.Port)
}
