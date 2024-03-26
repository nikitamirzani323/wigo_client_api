package controllers

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/wigo_client_api/entities"
	"bitbucket.org/isbtotogroup/wigo_client_api/helpers"
	"bitbucket.org/isbtotogroup/wigo_client_api/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const invoice_client_redis = "CLIENT:LISTINVOICE"
const invoice_result_redis = "CLIENT:RESULT"
const listmoney_redis = "CLIENT:LISTMONEY"

type c_tai struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
	Time    string      `json:"time"`
}

func CheckToken(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.CheckToken)
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

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	result := false
	client_company := ""
	client_username := ""
	client_name := ""
	client_credit := 0
	switch client.Token {
	case "qC5YmBvXzabGp34jJlKvnC6wCrr3pLCwBzsLoSzl4k=":
		client_company = "NUKE"
		client_username = "developer"
		client_name = "developer"
		client_credit = 100000
		result = true
	case "12345BvXzabGp34jJlKvnC6wCrr3pLCwBzsLoSzl4k=":
		client_company = "NUKE"
		client_username = "developer12"
		client_name = "developer12"
		client_credit = 500000
		result = true
	case "12345BvXzabGp34jJlKvnC6wCrr3pLCwBzsL1234567":
		client_company = "NUKE"
		client_username = "developer55"
		client_name = "developer55"
		client_credit = 1000000
		result = true
	}

	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})

	} else {
		var obj_record c_tai
		var obj entities.Model_listbet
		var arraobj []entities.Model_listbet
		resultredis, flag := helpers.GetRedis("nuke:" + listmoney_redis)
		jsonredis := []byte(resultredis)
		status_RD, _ := jsonparser.GetInt(jsonredis, "status")
		message_RD, _ := jsonparser.GetString(jsonredis, "message")
		record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
		jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			money_bet, _ := jsonparser.GetInt(value, "money_bet")
			obj.Money_bet = int(money_bet)
			arraobj = append(arraobj, obj)
		})
		obj_record.Status = int(status_RD)
		obj_record.Message = message_RD
		obj_record.Record = arraobj

		//CONFIG
		fieldconfig_redis := "CONFIG_ALL_NUKE"
		resultredis_conf, flag_conf := helpers.GetRedis(fieldconfig_redis)
		jsonredis_conf := []byte(resultredis_conf)
		// currRD, _ := jsonparser.GetString(jsonredis_conf, "curr")
		// minbetRD, _ := jsonparser.GetInt(jsonredis_conf, "minbet")
		// maxbetRD, _ := jsonparser.GetInt(jsonredis_conf, "maxbet")
		win_angkaRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_angka")
		win_redblackRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_redblack")
		win_lineRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_line")
		status_redblacklineRD, _ := jsonparser.GetString(jsonredis_conf, "status_redblackline")
		sstatus_maintenanceRD, _ := jsonparser.GetString(jsonredis_conf, "status_maintenance")

		engine_win_angka := 0.0
		engine_win_redblack := 0.0
		engine_win_line := 0.0
		engine_status_game_redblackline := ""
		engine_status_maintenance := ""

		if flag_conf {
			fmt.Println("CONF CACHE")

			engine_win_angka = float64(win_angkaRD)
			engine_win_redblack = float64(win_redblackRD)
			engine_win_line = float64(win_lineRD)
			engine_status_game_redblackline = status_redblacklineRD
			engine_status_maintenance = sstatus_maintenanceRD
		} else {
			fmt.Println("CONF DATABASE")

			win_angkaDB, win_redblackDB, win_lineDB, status_redblacklineDB, status_maintenanceDB := models.GetInfo_CompanyConf("NUKE")

			engine_win_angka = win_angkaDB
			engine_win_redblack = win_redblackDB
			engine_win_line = win_lineDB
			engine_status_game_redblackline = status_redblacklineDB
			engine_status_maintenance = status_maintenanceDB
		}

		if !flag {
			result, err := models.Fetch_listbet("NUKE")
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": err.Error(),
					"record":  nil,
				})
			}
			helpers.SetRedis(listmoney_redis+"_NUKE", result, 60*time.Minute)
			fmt.Println("LISTBET DATABASE")
			return c.JSON(fiber.Map{
				"status":                          fiber.StatusOK,
				"client_company":                  client_company,
				"client_name":                     client_name,
				"client_username":                 client_username,
				"client_credit":                   client_credit,
				"client_listbet":                  result,
				"engine_multiplier_angka":         engine_win_angka,
				"engine_multiplier_redblack":      engine_win_redblack,
				"engine_multiplier_line":          engine_win_line,
				"engine_status_game_redblackline": engine_status_game_redblackline,
				"engine_status_maintenance":       engine_status_maintenance,
			})
		} else {
			fmt.Println("LISTBET CACHE")
			return c.JSON(fiber.Map{
				"status":                          fiber.StatusOK,
				"client_company":                  client_company,
				"client_name":                     client_name,
				"client_username":                 client_username,
				"client_credit":                   client_credit,
				"client_listbet":                  obj_record,
				"engine_multiplier_angka":         engine_win_angka,
				"engine_multiplier_redblack":      engine_win_redblack,
				"engine_multiplier_line":          engine_win_line,
				"engine_status_game_redblackline": engine_status_game_redblackline,
				"engine_status_maintenance":       engine_status_maintenance,
			})
		}

	}
}

func ListInvoiceclient(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_invoice)
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

	var obj entities.Model_invoiceclient
	var arraobj []entities.Model_invoiceclient
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(strings.ToLower(client.Invoice_company) + ":" + invoice_client_redis + "_" + strings.ToLower(client.Invoice_username))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		invoiceclient_id, _ := jsonparser.GetString(value, "invoiceclient_id")
		invoiceclient_date, _ := jsonparser.GetString(value, "invoiceclient_date")
		invoiceclient_result, _ := jsonparser.GetString(value, "invoiceclient_result")
		invoiceclient_username, _ := jsonparser.GetString(value, "invoiceclient_username")
		invoiceclient_nomor, _ := jsonparser.GetString(value, "invoiceclient_nomor")
		invoiceclient_tipebet, _ := jsonparser.GetString(value, "invoiceclient_tipebet")
		invoiceclient_bet, _ := jsonparser.GetInt(value, "invoiceclient_bet")
		invoiceclient_win, _ := jsonparser.GetInt(value, "invoiceclient_win")
		invoiceclient_multiplier, _ := jsonparser.GetFloat(value, "invoiceclient_multiplier")
		invoiceclient_status, _ := jsonparser.GetString(value, "invoiceclient_status")
		invoiceclient_status_css, _ := jsonparser.GetString(value, "invoiceclient_status_css")

		obj.Invoiceclient_id = invoiceclient_id
		obj.Invoiceclient_date = invoiceclient_date
		obj.Invoiceclient_result = invoiceclient_result
		obj.Invoiceclient_username = invoiceclient_username
		obj.Invoiceclient_nomor = invoiceclient_nomor
		obj.Invoiceclient_tipebet = invoiceclient_tipebet
		obj.Invoiceclient_bet = int(invoiceclient_bet)
		obj.Invoiceclient_win = int(invoiceclient_win)
		obj.Invoiceclient_multiplier = float64(invoiceclient_multiplier)
		obj.Invoiceclient_status = invoiceclient_status
		obj.Invoiceclient_status_css = invoiceclient_status_css
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_invoice_client(client.Invoice_company, client.Invoice_username)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(strings.ToLower(client.Invoice_company)+":"+invoice_client_redis+"_"+strings.ToLower(client.Invoice_username), result, 5*time.Minute)
		fmt.Printf("INVOICECLIENT DATABASE %s-%s\n", client.Invoice_company, client.Invoice_username)
		return c.JSON(result)
	} else {
		fmt.Printf("INVOICECLIENT CACHE %s-%s\n", client.Invoice_company, client.Invoice_username)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func ListResult(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_result)
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

	var obj entities.Model_result
	var arraobj []entities.Model_result
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(strings.ToLower(client.Invoice_company) + ":" + invoice_result_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		result_invoice, _ := jsonparser.GetString(value, "result_invoice")
		result_date, _ := jsonparser.GetString(value, "result_date")
		result_result, _ := jsonparser.GetString(value, "result_result")

		obj.Result_invoice = result_invoice
		obj.Result_date = result_date
		obj.Result_result = result_result
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_result(client.Invoice_company)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(strings.ToLower(client.Invoice_company)+":"+invoice_result_redis, result, 30*time.Minute)
		fmt.Printf("RESULT DATABASE %s\n", client.Invoice_company)
		return c.JSON(result)
	} else {
		fmt.Printf("RESULT CACHE %s\n", client.Invoice_company)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func TransaksidetailSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_transaksidetailsave)
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
	// user := c.Locals("jwt").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	// temp_decp := helpers.Decryption(name)
	// client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	//idcompany, idtransaksi, username, listdatabet string, total_bet int
	result, err := models.Save_transaksidetail(client.Transaksidetail_company,
		client.Transaksidetail_idtransaksi, client.Transaksidetail_username,
		client.Transaksidetail_listdatabet, client.Transaksidetail_totalbet)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_wigo(client.Transaksidetail_company, client.Transaksidetail_username)
	return c.JSON(result)
}

func _deleteredis_wigo(company, username string) {
	val_invoice := helpers.DeleteRedis(strings.ToLower(company) + ":" + invoice_client_redis + "_" + strings.ToLower(username))
	fmt.Printf("Redis Delete INVOICE CLIENT : %d - %s %s\n", val_invoice, company, username)

	val_result := helpers.DeleteRedis(strings.ToLower(company) + ":" + invoice_result_redis)
	fmt.Printf("Redis Delete RESULT : %d - %s %s\n", val_result, company, username)

}
