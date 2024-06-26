package routers

import (
	"bitbucket.org/isbtotogroup/wigo_agen_api/controllers"
	"bitbucket.org/isbtotogroup/wigo_agen_api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		// c.Set("Content-Security-Policy", "frame-ancestors 'none'")
		// c.Set("X-XSS-Protection", "1; mode=block")
		// c.Set("X-Content-Type-Options", "nosniff")
		// c.Set("X-Download-Options", "noopen")
		// c.Set("Strict-Transport-Security", "max-age=5184000")
		// c.Set("X-Frame-Options", "SAMEORIGIN")
		// c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
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
	app.Post("/api/valid", middleware.JWTProtected(), controllers.Home)
	app.Post("/api/alladmin", middleware.JWTProtected(), controllers.Adminhome)
	app.Post("/api/saveadmin", middleware.JWTProtected(), controllers.AdminSave)
	app.Post("/api/alladminrule", middleware.JWTProtected(), controllers.Adminrulehome)
	app.Post("/api/saveadminrule", middleware.JWTProtected(), controllers.AdminruleSave)
	app.Post("/api/transaksi2d30s", middleware.JWTProtected(), controllers.Transaksi2D30Shome)
	app.Post("/api/transaksi2d30ssummarydaily", middleware.JWTProtected(), controllers.Transaksi2D30Ssummarydaily)
	app.Post("/api/transaksi2d30sinfo", middleware.JWTProtected(), controllers.Transaksi2D30Sinfo)
	app.Post("/api/transaksi2d30sprediksi", middleware.JWTProtected(), controllers.Transaksi2D30Sprediksi)
	app.Post("/api/transaksi2d30sdetail", middleware.JWTProtected(), controllers.Transaksi2D30Sdetail)
	app.Post("/api/transaksi2d30ssave", middleware.JWTProtected(), controllers.Transaksi2D30SSave)
	app.Post("/api/conf", middleware.JWTProtected(), controllers.AgenConf)
	app.Post("/api/confsave", middleware.JWTProtected(), controllers.AgenConfSave)

	return app
}
