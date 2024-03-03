package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/entities"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/helpers"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/models"
)

const invoice_client_redis = "CLIENT_LISTINVOICE"
const invoice_result_redis = "CLIENT_RESULT"

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

	// result, ruleadmin, err := models.Login_Model(client.Username, client.Password, client.Ipaddress)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	result := false
	if client.Token == "qC5YmBvXzabGp34jJlKvnC6wCrr3pLCwBzsLoSzl4k=" {
		result = true
	}

	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data Not Found",
			})

	} else {
		return c.JSON(fiber.Map{
			"status":            fiber.StatusOK,
			"client_company":    "ajuna",
			"client_name":       "developer",
			"client_username":   "developer212",
			"client_credit":     100000,
			"engine_multiplier": 5,
		})
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
	resultredis, flag := helpers.GetRedis(invoice_client_redis + "_" + strings.ToLower(client.Invoice_company) + "_" + strings.ToLower(client.Invoice_username))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		invoiceclient_id, _ := jsonparser.GetString(value, "invoiceclient_id")
		invoiceclient_date, _ := jsonparser.GetString(value, "invoiceclient_date")
		invoiceclient_result, _ := jsonparser.GetString(value, "invoiceclient_result")
		invoiceclient_username, _ := jsonparser.GetString(value, "invoiceclient_username")
		invoiceclient_nomor, _ := jsonparser.GetString(value, "invoiceclient_nomor")
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
		helpers.SetRedis(invoice_client_redis+"_"+strings.ToLower(client.Invoice_company)+"_"+strings.ToLower(client.Invoice_username), result, 5*time.Minute)
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
	resultredis, flag := helpers.GetRedis(invoice_result_redis + "_" + strings.ToLower(client.Invoice_company))
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
		helpers.SetRedis(invoice_result_redis+"_"+strings.ToLower(client.Invoice_company), result, 30*time.Minute)
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

	//idcompany, idtransaksi, username, nomor string, bet int, multiplier float64
	result, err := models.Save_transaksidetail(client.Transaksidetail_company,
		client.Transaksidetail_idtransaksi, client.Transaksidetail_username,
		client.Transaksidetail_nomor, client.Transaksidetail_bet, client.Transaksidetail_multiplier)

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
	val_invoice := helpers.DeleteRedis(invoice_client_redis + "_" + strings.ToLower(company) + "_" + strings.ToLower(username))
	fmt.Printf("Redis Delete INVOICE CLIENT : %d - %s %s\n", val_invoice, company, username)

	val_result := helpers.DeleteRedis(invoice_result_redis + "_" + strings.ToLower(company))
	fmt.Printf("Redis Delete RESULT : %d - %s %s\n", val_result, company, username)

}
