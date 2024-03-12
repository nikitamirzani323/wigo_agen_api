package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/WIGO_AGEN_API/entities"
	"github.com/nikitamirzani323/WIGO_AGEN_API/helpers"
	"github.com/nikitamirzani323/WIGO_AGEN_API/models"
)

const Fieldtransaksi2d30s_home_redis = "LISTINVOICE_2D30S_AGEN"

func Transaksi2D30Shome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_transaksi2D30S)
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
	_, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	fieldredis := ""
	if client.Transaksi2D30S_invoice != "" {
		fieldredis = Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(client_company) + "_" + client.Transaksi2D30S_invoice
	} else {
		if client.Transaksi2D30S_search != "" {
			fieldredis = Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(client_company) + "_" + strconv.Itoa(client.Transaksi2D30S_page) + "_" + client.Transaksi2D30S_search
			val_pattern := helpers.DeleteRedis(fieldredis)
			fmt.Printf("Redis Delete INVOICE : %d", val_pattern)
		} else {
			fieldredis = Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(client_company) + "_" + strconv.Itoa(client.Transaksi2D30S_page) + "_" + client.Transaksi2D30S_search
		}
	}

	var obj entities.Model_transaksi2D30S
	var arraobj []entities.Model_transaksi2D30S
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(fieldredis)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalbet_RD, _ := jsonparser.GetInt(jsonredis, "totalbet")
	totalwin_RD, _ := jsonparser.GetInt(jsonredis, "totalwin")
	winlose_agen_RD, _ := jsonparser.GetInt(jsonredis, "winlose_agen")
	winlose_member_RD, _ := jsonparser.GetInt(jsonredis, "winlose_member")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	periode_RD, _ := jsonparser.GetString(jsonredis, "periode")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		transaksi2D30s_id, _ := jsonparser.GetString(value, "transaksi2D30s_id")
		transaksi2D30s_idcurr, _ := jsonparser.GetString(value, "transaksi2D30s_idcurr")
		transaksi2D30s_date, _ := jsonparser.GetString(value, "transaksi2D30s_date")
		transaksi2D30s_result, _ := jsonparser.GetString(value, "transaksi2D30s_result")
		transaksi2D30s_totalmember, _ := jsonparser.GetInt(value, "transaksi2D30s_totalmember")
		transaksi2D30s_totalbet, _ := jsonparser.GetInt(value, "transaksi2D30s_totalbet")
		transaksi2D30s_totalwin, _ := jsonparser.GetInt(value, "transaksi2D30s_totalwin")
		transaksi2D30s_winlose, _ := jsonparser.GetInt(value, "transaksi2D30s_winlose")
		transaksi2D30s_status, _ := jsonparser.GetString(value, "transaksi2D30s_status")
		transaksi2D30s_status_css, _ := jsonparser.GetString(value, "transaksi2D30s_status_css")
		transaksi2D30s_create, _ := jsonparser.GetString(value, "transaksi2D30s_create")
		transaksi2D30s_update, _ := jsonparser.GetString(value, "transaksi2D30s_update")

		obj.Transaksi2D30S_id = transaksi2D30s_id
		obj.Transaksi2D30S_idcurr = transaksi2D30s_idcurr
		obj.Transaksi2D30S_date = transaksi2D30s_date
		obj.Transaksi2D30S_result = transaksi2D30s_result
		obj.Transaksi2D30S_totalmember = int(transaksi2D30s_totalmember)
		obj.Transaksi2D30S_totalbet = int(transaksi2D30s_totalbet)
		obj.Transaksi2D30S_totalwin = int(transaksi2D30s_totalwin)
		obj.Transaksi2D30S_winlose = int(transaksi2D30s_winlose)
		obj.Transaksi2D30S_status = transaksi2D30s_status
		obj.Transaksi2D30S_status_css = transaksi2D30s_status_css
		obj.Transaksi2D30S_create = transaksi2D30s_create
		obj.Transaksi2D30S_update = transaksi2D30s_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		//idcompany, idinvoice, search string, page int
		result, err := models.Fetch_transaksi2D30SHome(client_company, client.Transaksi2D30S_invoice, client.Transaksi2D30S_search, client.Transaksi2D30S_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(fieldredis, result, 60*time.Minute)
		fmt.Println("TRANSAKSI 2D30S DATABASE")
		return c.JSON(result)
	} else {
		fmt.Println("TRANSAKSI 2D30S CACHE")
		return c.JSON(fiber.Map{
			"status":         fiber.StatusOK,
			"message":        "Success",
			"record":         arraobj,
			"periode":        periode_RD,
			"perpage":        perpage_RD,
			"totalrecord":    totalrecord_RD,
			"totalbet":       totalbet_RD,
			"totalwin":       totalwin_RD,
			"winlose_agen":   winlose_agen_RD,
			"winlose_member": winlose_member_RD,
			"time":           time.Since(render_page).String(),
		})
	}
}
func Transaksi2D30Sdetail(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_transaksidetail2D30S)
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
	_, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	fieldredis := ""
	if client.Transaksidetail2D30S_invoice != "" {
		fieldredis = Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(client_company) + "_DETAIL_" + client.Transaksidetail2D30S_invoice + "_" + client.Transaksidetail2D30S_status
	} else {
		fieldredis = Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(client_company) + "_DETAIL_" + client.Transaksidetail2D30S_status
	}

	var obj entities.Model_transaksi2D30Sdetail
	var arraobj []entities.Model_transaksi2D30Sdetail
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(fieldredis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		transaksi2D30sdetail_id, _ := jsonparser.GetString(value, "transaksi2D30sdetail_id")
		transaksi2D30sdetail_date, _ := jsonparser.GetString(value, "transaksi2D30sdetail_date")
		transaksi2D30sdetail_ipaddress, _ := jsonparser.GetString(value, "transaksi2D30sdetail_ipaddress")
		transaksi2D30sdetail_device, _ := jsonparser.GetString(value, "transaksi2D30sdetail_device")
		transaksi2D30sdetail_browser, _ := jsonparser.GetString(value, "transaksi2D30sdetail_browser")
		transaksi2D30sdetail_username, _ := jsonparser.GetString(value, "transaksi2D30sdetail_username")
		transaksi2D30sdetail_tipebet, _ := jsonparser.GetString(value, "transaksi2D30sdetail_tipebet")
		transaksi2D30sdetail_nomor, _ := jsonparser.GetString(value, "transaksi2D30sdetail_nomor")
		transaksi2D30sdetail_bet, _ := jsonparser.GetInt(value, "transaksi2D30sdetail_bet")
		transaksi2D30sdetail_win, _ := jsonparser.GetInt(value, "transaksi2D30sdetail_win")
		transaksi2D30sdetail_multiplier, _ := jsonparser.GetFloat(value, "transaksi2D30sdetail_multiplier")
		transaksi2D30sdetail_status, _ := jsonparser.GetString(value, "transaksi2D30sdetail_status")
		transaksi2D30sdetail_status_css, _ := jsonparser.GetString(value, "transaksi2D30sdetail_status_css")
		transaksi2D30sdetail_create, _ := jsonparser.GetString(value, "transaksi2D30sdetail_create")
		transaksi2D30sdetail_update, _ := jsonparser.GetString(value, "transaksi2D30sdetail_update")

		obj.Transaksi2D30Sdetail_id = transaksi2D30sdetail_id
		obj.Transaksi2D30Sdetail_date = transaksi2D30sdetail_date
		obj.Transaksi2D30Sdetail_ipaddress = transaksi2D30sdetail_ipaddress
		obj.Transaksi2D30Sdetail_browser = transaksi2D30sdetail_browser
		obj.Transaksi2D30Sdetail_device = transaksi2D30sdetail_device
		obj.Transaksi2D30Sdetail_username = transaksi2D30sdetail_username
		obj.Transaksi2D30Sdetail_tipebet = transaksi2D30sdetail_tipebet
		obj.Transaksi2D30Sdetail_nomor = transaksi2D30sdetail_nomor
		obj.Transaksi2D30Sdetail_bet = int(transaksi2D30sdetail_bet)
		obj.Transaksi2D30Sdetail_win = int(transaksi2D30sdetail_win)
		obj.Transaksi2D30Sdetail_multiplier = float64(transaksi2D30sdetail_multiplier)
		obj.Transaksi2D30Sdetail_status = transaksi2D30sdetail_status
		obj.Transaksi2D30Sdetail_status_css = transaksi2D30sdetail_status_css
		obj.Transaksi2D30Sdetail_create = transaksi2D30sdetail_create
		obj.Transaksi2D30Sdetail_update = transaksi2D30sdetail_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		//idcompany, idtransaksi, status string
		result, err := models.Fetch_transaksi2D30SDetail(client_company, client.Transaksidetail2D30S_invoice, client.Transaksidetail2D30S_status)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(fieldredis, result, 60*time.Minute)
		fmt.Println("TRANSAKSI DETAIL 2D30S DATABASE")
		return c.JSON(result)
	} else {
		fmt.Println("TRANSAKSI DETAIL 2D30S CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Transaksi2D30Sprediksi(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_prediksitransaksi2D30S)
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
	_, client_company, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_transaksi2D30SPrediksi
	var arraobj []entities.Model_transaksi2D30SPrediksi
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldtransaksi2d30s_home_redis + "_PREDIKSI_" + strings.ToLower(client_company) + "_" + client.Transaksi2D30Sprediksi_invoice + "_" + client.Transaksi2D30Sprediksi_result)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	totalmember_RD, _ := jsonparser.GetInt(jsonredis, "totalmember")
	totalbet_RD, _ := jsonparser.GetInt(jsonredis, "totalbet")
	totalwin_RD, _ := jsonparser.GetInt(jsonredis, "totalwin")
	winlose_RD, _ := jsonparser.GetInt(jsonredis, "winlose")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		transaksi2D30sprediksi_id, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_id")
		transaksi2D30sprediksi_username, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_username")
		transaksi2D30sprediksi_date, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_date")
		transaksi2D30sprediksi_nomor, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_nomor")
		transaksi2D30sprediksi_bet, _ := jsonparser.GetInt(value, "transaksi2D30sprediksi_bet")
		transaksi2D30sprediksi_win, _ := jsonparser.GetInt(value, "transaksi2D30sprediksi_win")
		transaksi2D30sprediksi_winlose, _ := jsonparser.GetInt(value, "transaksi2D30sprediksi_winlose")
		transaksi2D30sprediksi_status, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_status")
		transaksi2D30sprediksi_status_css, _ := jsonparser.GetString(value, "transaksi2D30sprediksi_status_css")

		obj.Transaksi2D30Sprediksi_id = transaksi2D30sprediksi_id
		obj.Transaksi2D30Sprediksi_date = transaksi2D30sprediksi_date
		obj.Transaksi2D30Sprediksi_username = transaksi2D30sprediksi_username
		obj.Transaksi2D30Sprediksi_nomor = transaksi2D30sprediksi_nomor
		obj.Transaksi2D30Sprediksi_bet = int(transaksi2D30sprediksi_bet)
		obj.Transaksi2D30Sprediksi_win = int(transaksi2D30sprediksi_win)
		obj.Transaksi2D30Sprediksi_winlose = int(transaksi2D30sprediksi_winlose)
		obj.Transaksi2D30Sprediksi_status = transaksi2D30sprediksi_status
		obj.Transaksi2D30Sprediksi_status_css = transaksi2D30sprediksi_status_css
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_transaksi2D30SPrediksi(client_company, client.Transaksi2D30Sprediksi_invoice, client.Transaksi2D30Sprediksi_result)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldtransaksi2d30s_home_redis+"_PREDIKSI_"+strings.ToLower(client_company)+"_"+client.Transaksi2D30Sprediksi_invoice+"_"+client.Transaksi2D30Sprediksi_result, result, 10*time.Minute)
		fmt.Println("TRANSAKSI 2D30S PREDIKSI DATABASE")
		return c.JSON(result)
	} else {
		fmt.Println("TRANSAKSI 2D30S PREDIKSI CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"totalmember": totalmember_RD,
			"totalbet":    totalbet_RD,
			"totalwin":    totalwin_RD,
			"winlose":     winlose_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}

func Transaksi2D30SSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_transaksi2D30Ssave)
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

	// admin, idrecord, idcompany, result string
	result, err := models.Save_updateresult2D30S(
		client_admin,
		client.Transaksi2D30S_invoice, client_company, client.Transaksi2D30S_result)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_transaksi(client.Transaksi2D30S_invoice, client_company)
	return c.JSON(result)
}

func _deleteredis_transaksi(idinvoice, idcompany string) {
	val_transaksi2d30s := helpers.DeleteRedis(Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(idcompany))
	fmt.Printf("Redis Delete AGEN TRANSAKSI2D30S INVOICE : %d", val_transaksi2d30s)
	val_transaksi2d30s2 := helpers.DeleteRedis(Fieldtransaksi2d30s_home_redis + "_" + strings.ToLower(idcompany) + "_" + idinvoice)
	fmt.Printf("Redis Delete AGEN TRANSAKSI2D30S INVOICE : %d", val_transaksi2d30s2)
}
