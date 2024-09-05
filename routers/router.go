package routers

import (
	"bitbucket.org/isbtotogroup/wigo_agen_api/controllers"
	"bitbucket.org/isbtotogroup/wigo_agen_api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(requestid.New())
	app.Use(etag.New())
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/"
		},
		Format: "${time} | ${status} | ${latency} | ${ips[0]} | ${method} | ${path} - ${queryParams} ${body}\n",
	}))
	app.Get("/ipaddress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      "data",
			"BASEURL":     c.BaseURL(),
			"HOSTNAME":    c.Hostname(),
			"IP":          c.IP(),
			"IPS":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})
	app.Get("/dashboard", monitor.New())

	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/checkdomain", controllers.Domaincheck)

	api := app.Group("/api", middleware.JWTProtected())
	api.Post("/valid", controllers.Home)
	api.Post("/alladmin", controllers.Adminhome)
	api.Post("/saveadmin", controllers.AdminSave)
	api.Post("/alladminrule", controllers.Adminrulehome)
	api.Post("/saveadminrule", controllers.AdminruleSave)
	api.Post("/transaksi2d30s", controllers.Transaksi2D30Shome)
	api.Post("/transaksi2d30ssummarydaily", controllers.Transaksi2D30Ssummarydaily)
	api.Post("/transaksi2d30sinfo", controllers.Transaksi2D30Sinfo)
	api.Post("/transaksi2d30sprediksi", controllers.Transaksi2D30Sprediksi)
	api.Post("/transaksi2d30sdetail", controllers.Transaksi2D30Sdetail)
	api.Post("/transaksi2d30ssave", controllers.Transaksi2D30SSave)
	api.Post("/conf", controllers.AgenConf)
	api.Post("/confsave", controllers.AgenConfSave)

	return app
}
