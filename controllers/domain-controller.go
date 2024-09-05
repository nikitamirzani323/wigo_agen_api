package controllers

import (
	"bitbucket.org/isbtotogroup/wigo_agen_api/entities"
	"bitbucket.org/isbtotogroup/wigo_agen_api/helpers"
	"bitbucket.org/isbtotogroup/wigo_agen_api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Domaincheck(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_domaincheck)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	result, err := models.Fetch_checkdomain(client.Domain, client.Tipe)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status":       fiber.StatusUnauthorized,
			"domainstatus": err.Error(),
		})
	}

	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":       fiber.StatusUnauthorized,
				"domainstatus": "Username or Password Not Found",
			})

	} else {
		return c.JSON(fiber.Map{
			"status":       fiber.StatusOK,
			"domainstatus": "ACTIVE",
		})

	}
}
