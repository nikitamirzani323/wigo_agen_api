package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nikitamirzani323/WIGO_AGEN_API/controllers"
	"github.com/nikitamirzani323/WIGO_AGEN_API/middleware"
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
	app.Post("/api/transaksi2d30sprediksi", middleware.JWTProtected(), controllers.Transaksi2D30Sprediksi)
	app.Post("/api/transaksi2d30ssave", middleware.JWTProtected(), controllers.Transaksi2D30SSave)

	app.Post("/api/curr", middleware.JWTProtected(), controllers.Currhome)
	app.Post("/api/currsave", middleware.JWTProtected(), controllers.CurrSave)
	app.Post("/api/company", middleware.JWTProtected(), controllers.Companyhome)
	app.Post("/api/companyadmin", middleware.JWTProtected(), controllers.Companyadminhome)
	app.Post("/api/companyadminrule", middleware.JWTProtected(), controllers.Companyadminrulehome)
	app.Post("/api/companymoney", middleware.JWTProtected(), controllers.Companymoneyhome)
	app.Post("/api/companyconf", middleware.JWTProtected(), controllers.Companyconfhome)
	app.Post("/api/companysave", middleware.JWTProtected(), controllers.CompanySave)
	app.Post("/api/companyadminsave", middleware.JWTProtected(), controllers.CompanyadminSave)
	app.Post("/api/companyadminrulesave", middleware.JWTProtected(), controllers.CompanyadminruleSave)
	app.Post("/api/companymoneysave", middleware.JWTProtected(), controllers.CompanymoneySave)
	app.Post("/api/companymoneydelete", middleware.JWTProtected(), controllers.CompanymoneyDelete)
	app.Post("/api/companyconfsave", middleware.JWTProtected(), controllers.CompanyconfSave)

	return app
}
