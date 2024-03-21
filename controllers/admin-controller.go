package controllers

import (
	"log"
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

const Fieldadmin_home_redis = "LISTADMIN_AGEN"

func Adminhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var obj_listruleadmin entities.Model_adminrule
	var arraobj_listruleadmin []entities.Model_adminrule
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldadmin_home_redis + "_" + strings.ToLower(client_company))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listruleadmin_RD, _, _, _ := jsonparser.Get(jsonredis, "listruleadmin")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		admin_id, _ := jsonparser.GetInt(value, "admin_id")
		admin_username, _ := jsonparser.GetString(value, "admin_username")
		admin_nama, _ := jsonparser.GetString(value, "admin_nama")
		admin_idrule, _ := jsonparser.GetInt(value, "admin_idrule")
		admin_rule, _ := jsonparser.GetString(value, "admin_rule")
		admin_joindate, _ := jsonparser.GetString(value, "admin_joindate")
		admin_lastlogin, _ := jsonparser.GetString(value, "admin_lastlogin")
		admin_lastipaddres, _ := jsonparser.GetString(value, "admin_lastipaddres")
		admin_status, _ := jsonparser.GetString(value, "admin_status")
		admin_status_css, _ := jsonparser.GetString(value, "admin_status_css")

		obj.Admin_id = int(admin_id)
		obj.Admin_username = admin_username
		obj.Admin_idrule = int(admin_idrule)
		obj.Admin_rule = admin_rule
		obj.Admin_nama = admin_nama
		obj.Admin_joindate = admin_joindate
		obj.Admin_lastlogin = admin_lastlogin
		obj.Admin_lastIpaddress = admin_lastipaddres
		obj.Admin_status = admin_status
		obj.Admin_status_css = admin_status_css
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listruleadmin_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		adminrule_id, _ := jsonparser.GetInt(value, "adminrule_id")
		adminrule_name, _ := jsonparser.GetString(value, "adminrule_name")

		obj_listruleadmin.Adminrule_id = int(adminrule_id)
		obj_listruleadmin.Adminrule_name = adminrule_name
		arraobj_listruleadmin = append(arraobj_listruleadmin, obj_listruleadmin)
	})
	if !flag {
		result, err := models.Fetch_adminHome(client_company)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldadmin_home_redis+"_"+strings.ToLower(client_company), result, 60*time.Minute)
		log.Println("ADMIN DATABASE")
		return c.JSON(result)
	} else {
		log.Println("ADMIN CACHE")
		return c.JSON(fiber.Map{
			"status":        fiber.StatusOK,
			"message":       "Success",
			"record":        arraobj,
			"listruleadmin": arraobj_listruleadmin,
			"time":          time.Since(render_page).String(),
		})
	}
}

func AdminSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminsave)
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

	//admin, idcompany, username, password, nama, status, sData string, idrecord, idrule int
	result, err := models.Save_adminHome(
		client_admin, client_company,
		client.Admin_username, client.Admin_password, client.Admin_nama,
		client.Admin_status, client.Sdata, client.Admin_id, client.Admin_idrule)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_admin(client_company)
	return c.JSON(result)
}
func _deleteredis_admin(idcompany string) {
	val_master := helpers.DeleteRedis(Fieldadmin_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("Redis Delete AGEN ADMIN : %d", val_master)

}
