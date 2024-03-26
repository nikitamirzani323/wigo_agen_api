package controllers

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/wigo_agen_api/entities"
	"bitbucket.org/isbtotogroup/wigo_agen_api/helpers"
	"bitbucket.org/isbtotogroup/wigo_agen_api/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Fieldadminrule_home_redis = "AGEN:LISTRULE"

func Adminrulehome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_adminruleall
	var arraobj []entities.Model_adminruleall
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(strings.ToLower(client_company) + ":" + Fieldadminrule_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		adminrule_id, _ := jsonparser.GetInt(value, "adminrule_id")
		adminrule_name, _ := jsonparser.GetString(value, "adminrule_name")
		adminrule_rule, _ := jsonparser.GetString(value, "adminrule_rule")

		obj.Adminrule_id = int(adminrule_id)
		obj.Adminrule_name = adminrule_name
		obj.Adminrule_rule = adminrule_rule
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_adminruleHome(client_company)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(strings.ToLower(client_company)+":"+Fieldadminrule_home_redis, result, 60*time.Minute)
		fmt.Println("AGEN ADMIN RULE DATABASE")
		return c.JSON(result)
	} else {
		fmt.Println("AGEN ADMIN RULE CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func AdminruleSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminrulesave)
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
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	//admin, idcompany, name, rule, sData string, idrecord int
	result, err := models.Save_adminrule(client_admin, client_company,
		client.Adminrule_name, client.Adminrule_rule, client.Sdata, client.Adminrule_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_adminrule(client_company)
	return c.JSON(result)
}

func _deleteredis_adminrule(idcompany string) {
	val_master := helpers.DeleteRedis(strings.ToLower(idcompany) + ":" + Fieldadminrule_home_redis)
	fmt.Printf("Redis Delete AGEN ADMIN RULE : %d", val_master)

}
