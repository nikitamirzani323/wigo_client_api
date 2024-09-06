package controllers

import (
	"fmt"
	"log"
	"strconv"
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
const fieldlogin_redis = "CLIENT_LOGIN"

type c_tai struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
	Time    string      `json:"time"`
}
type C_InfoLogin struct {
	Client_company  string `json:"client_company"`
	Client_username string `json:"client_username"`
	Client_credit   int    `json:"client_credit"`
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
		client_credit = 1000000
		result = true
	case "12345BvXzabGp34jJlKvnC6wCrr3pLCwBzsLoSzl4k=":
		client_company = "NUKE"
		client_username = "developer12"
		client_name = "developer12"
		client_credit = 5000000
		result = true
	case "12345BvXzabGp34jJlKvnC6wCrr3pLCwBzsL1234567":
		client_company = "NUKE"
		client_username = "developer55"
		client_name = "developer55"
		client_credit = 10000000
		result = true
	}

	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})

	} else {
		logininfo_redis, flag_loginfo := helpers.GetRedis(fieldlogin_redis + "_" + client.Token)
		jsonredis_info := []byte(logininfo_redis)
		client_companyRD, _ := jsonparser.GetString(jsonredis_info, "client_company")
		client_usernameRD, _ := jsonparser.GetString(jsonredis_info, "client_username")
		client_creditRD, _ := jsonparser.GetInt(jsonredis_info, "client_credit")
		fmt.Println("Data Redis : " + client_companyRD + " - " + client_usernameRD)
		if flag_loginfo {
			client_credit = int(client_creditRD)
		} else {
			var objlogin_record C_InfoLogin
			objlogin_record.Client_company = client_company
			objlogin_record.Client_username = client_username
			objlogin_record.Client_credit = client_credit
			helpers.SetRedis(fieldlogin_redis+"_"+client.Token, objlogin_record, 1440*time.Minute)
		}

		var obj_record c_tai
		var obj entities.Model_listbet
		var arraobj []entities.Model_listbet
		resultredis, flag := helpers.GetRedis(strings.ToLower(client_company) + ":" + listmoney_redis)
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

		fieldconfig_redis := strings.ToLower(client_company) + ":12D30S:CONFIG"
		resultredis_conf, flag_conf := helpers.GetRedis(fieldconfig_redis)
		jsonredis_conf := []byte(resultredis_conf)
		// currRD, _ := jsonparser.GetString(jsonredis_conf, "curr")
		// minbetRD, _ := jsonparser.GetInt(jsonredis_conf, "minbet")
		// maxbetRD, _ := jsonparser.GetInt(jsonredis_conf, "maxbet")
		win_angkaRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_angka")
		win_redblackRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_redblack")
		win_lineRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_line")
		win_zonaRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_zona")
		win_jackpotRD, _ := jsonparser.GetFloat(jsonredis_conf, "win_jackpot")
		status_redblacklineRD, _ := jsonparser.GetString(jsonredis_conf, "status_redblackline")
		sstatus_maintenanceRD, _ := jsonparser.GetString(jsonredis_conf, "status_maintenance")

		engine_win_angka := 0.0
		engine_win_redblack := 0.0
		engine_win_line := 0.0
		engine_win_zona := 0.0
		engine_win_jackpot := 0.0
		engine_status_game_redblackline := ""
		engine_status_maintenance := ""

		if flag_conf {
			fmt.Println("CONF CACHE")
			log.Println(win_zonaRD)
			log.Println(win_jackpotRD)
			engine_win_angka = float64(win_angkaRD)
			engine_win_redblack = float64(win_redblackRD)
			engine_win_line = float64(win_lineRD)
			engine_win_zona = float64(win_zonaRD)
			engine_win_jackpot = float64(win_jackpotRD)
			engine_status_game_redblackline = status_redblacklineRD
			engine_status_maintenance = sstatus_maintenanceRD
		} else {
			fmt.Println("CONF DATABASE")

			win_angkaDB, win_redblackDB, win_lineDB, win_zonaDB, win_jackpotDB, status_redblacklineDB, status_maintenanceDB := models.GetInfo_CompanyConf("NUKE")

			engine_win_angka = win_angkaDB
			engine_win_redblack = win_redblackDB
			engine_win_line = win_lineDB
			engine_win_zona = win_zonaDB
			engine_win_jackpot = win_jackpotDB
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
				"engine_multiplier_zona":          engine_win_zona,
				"engine_multiplier_jackpot":       engine_win_jackpot,
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
				"engine_multiplier_zona":          engine_win_zona,
				"engine_multiplier_jackpot":       engine_win_jackpot,
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

	logininfo_redis, flag_loginfo := helpers.GetRedis(fieldlogin_redis + "_" + client.Client_token)
	jsonredis_info := []byte(logininfo_redis)
	client_companyRD, _ := jsonparser.GetString(jsonredis_info, "client_company")
	client_usernameRD, _ := jsonparser.GetString(jsonredis_info, "client_username")
	// client_creditRD, _ := jsonparser.GetInt(jsonredis, "client_credit")
	// fmt.Println("Data Redis : " + client_companyRD + " - " + client_usernameRD)
	if flag_loginfo {
		var obj entities.Model_invoiceclient
		var arraobj []entities.Model_invoiceclient
		render_page := time.Now()
		resultredis, flag := helpers.GetRedis(strings.ToLower(client_companyRD) + ":" + invoice_client_redis + "_" + strings.ToLower(client_usernameRD))
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
			result, err := models.Fetch_invoice_client(strings.ToLower(client_companyRD), strings.ToLower(client_usernameRD))
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": err.Error(),
					"record":  nil,
				})
			}
			helpers.SetRedis(strings.ToLower(client_companyRD)+":"+invoice_client_redis+"_"+strings.ToLower(client_usernameRD), result, 5*time.Minute)
			fmt.Printf("INVOICECLIENT DATABASE %s-%s\n", client_companyRD, client_usernameRD)
			return c.JSON(result)
		} else {
			fmt.Printf("INVOICECLIENT CACHE %s-%s\n", client_companyRD, client_usernameRD)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusOK,
				"message": "Success",
				"record":  arraobj,
				"time":    time.Since(render_page).String(),
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})
	}

}
func ListResult(c *fiber.Ctx) error {
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

	logininfo_redis, flag_loginfo := helpers.GetRedis(fieldlogin_redis + "_" + client.Client_token)
	jsonredis_info := []byte(logininfo_redis)
	client_companyRD, _ := jsonparser.GetString(jsonredis_info, "client_company")
	// client_usernameRD, _ := jsonparser.GetString(jsonredis_info, "client_username")
	// client_creditRD, _ := jsonparser.GetInt(jsonredis, "client_credit")
	// fmt.Println("Data Redis : " + client_companyRD + " - " + client_usernameRD)

	if flag_loginfo {
		var obj entities.Model_result
		var arraobj []entities.Model_result
		render_page := time.Now()
		resultredis, flag := helpers.GetRedis(strings.ToLower(client_companyRD) + ":" + invoice_result_redis)
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
			result, err := models.Fetch_result(strings.ToLower(client_companyRD))
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": err.Error(),
					"record":  nil,
				})
			}
			helpers.SetRedis(strings.ToLower(client_companyRD)+":"+invoice_result_redis, result, 30*time.Minute)
			fmt.Printf("RESULT DATABASE %s\n", client_companyRD)
			return c.JSON(result)
		} else {
			fmt.Printf("RESULT CACHE %s\n", client_companyRD)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusOK,
				"message": "Success",
				"record":  arraobj,
				"time":    time.Since(render_page).String(),
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
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

	logininfo_redis, flag_loginfo := helpers.GetRedis(fieldlogin_redis + "_" + client.Client_token)
	jsonredis_info := []byte(logininfo_redis)
	client_companyRD, _ := jsonparser.GetString(jsonredis_info, "client_company")
	client_usernameRD, _ := jsonparser.GetString(jsonredis_info, "client_username")
	client_creditRD, _ := jsonparser.GetInt(jsonredis_info, "client_credit")
	// fmt.Println("Data Redis : " + client_companyRD + " - " + client_usernameRD)

	if flag_loginfo {
		if int(client_creditRD) > client.Transaksidetail_totalbet {
			var objlogin_record C_InfoLogin
			objlogin_record.Client_company = client_companyRD
			objlogin_record.Client_username = client_usernameRD
			objlogin_record.Client_credit = int(client_creditRD) - client.Transaksidetail_totalbet
			helpers.SetRedis(fieldlogin_redis+"_"+client.Client_token, objlogin_record, 1440*time.Minute)

			//BETROUND + GENERATE INVOICE
			round_redis := "round_detail_" + strings.ToLower(client_companyRD) + ":" + client.Transaksidetail_idtransaksi
			_, flaground_redis := helpers.GetRedis(round_redis)
			if !flaground_redis {
				result := models.Get_counter(round_redis)
				fmt.Println("init redis betround")
				helpers.SetRedis(round_redis, result-1, 24*time.Hour)
			}

			resultround := helpers.IncrPipeRedis(round_redis, "2", 24*time.Hour)
			roundBet, _ := strconv.Atoi(resultround)
			invoiceplayer := client.Transaksidetail_idtransaksi + strconv.Itoa(roundBet)
			fmt.Println(client_usernameRD, ": bet: get generatedinvoice: ", invoiceplayer)

			go models.Update_counter(round_redis, roundBet)

			//idcompany, idtransaksi, username, playerinvoice, listdatabet string, total_bet, betround int
			result, err := models.Save_transaksidetail(strings.ToLower(client_companyRD),
				client.Transaksidetail_idtransaksi, strings.ToLower(client_usernameRD), invoiceplayer,
				client.Transaksidetail_listdatabet, client.Transaksidetail_totalbet, roundBet)

			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": err.Error(),
					"record":  nil,
				})
			}
			_deleteredis_wigo(client_companyRD, client_usernameRD)
			return c.JSON(result)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": "Credit Tidak Cukup",
				})
		}

	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})
	}
}
func Balance(c *fiber.Ctx) error {
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
	render_page := time.Now()
	logininfo_redis, flag_loginfo := helpers.GetRedis(fieldlogin_redis + "_" + client.Client_token)
	jsonredis_info := []byte(logininfo_redis)
	// client_companyRD, _ := jsonparser.GetString(jsonredis_info, "client_company")
	// client_usernameRD, _ := jsonparser.GetString(jsonredis_info, "client_username")
	client_creditRD, _ := jsonparser.GetInt(jsonredis_info, "client_credit")
	// fmt.Println("Data Redis : " + client_companyRD + " - " + client_usernameRD)

	if flag_loginfo {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"credit":  client_creditRD,
			"time":    time.Since(render_page).String(),
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})
	}
}
func _deleteredis_wigo(company, username string) {
	val_invoice := helpers.DeleteRedis(strings.ToLower(company) + ":" + invoice_client_redis + "_" + strings.ToLower(username))
	fmt.Printf("Redis Delete INVOICE CLIENT : %d - %s %s\n", val_invoice, company, username)

	val_result := helpers.DeleteRedis(strings.ToLower(company) + ":" + invoice_result_redis)
	fmt.Printf("Redis Delete RESULT : %d - %s %s\n", val_result, company, username)

}
