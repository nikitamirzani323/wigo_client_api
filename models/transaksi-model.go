package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/wigo_client_api/configs"
	"bitbucket.org/isbtotogroup/wigo_client_api/db"
	"bitbucket.org/isbtotogroup/wigo_client_api/entities"
	"bitbucket.org/isbtotogroup/wigo_client_api/helpers"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_listbet(idcompany string) (helpers.Response, error) {
	var obj entities.Model_listbet
	var arraobj []entities.Model_listbet
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "compmoney  "
	sql_select += "FROM " + configs.DB_tbl_mst_company_money + " "
	sql_select += "WHERE idcompany ='" + idcompany + "' "
	sql_select += "ORDER BY compmoney ASC  "

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			compmoney_db int
		)

		err = row.Scan(&compmoney_db)

		helpers.ErrorCheck(err)

		obj.Money_bet = compmoney_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_invoice_client(idcompany, username string) (helpers.Response, error) {
	var obj entities.Model_invoiceclient
	var arraobj []entities.Model_invoiceclient
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, _, _, view_invoice_client := Get_mappingdatabase(idcompany)

	tglnow, _ := goment.New()
	tglskrg := tglnow.Format("YYYY-MM-DD HH:mm:ss")
	tglbefore := tglnow.Add(-31, "days").Format("YYYY-MM-DD HH:mm:ss")

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi, to_char(COALESCE(tgltransaksi,now()), 'YYYY-MM-DD HH24:MI:SS') as tgltransaksi,  "
	sql_select += "username_client, resultwigo, nomor, tipebet,  "
	sql_select += "bet, win, multiplier, status_transaksidetail    "
	sql_select += "FROM " + view_invoice_client + " "
	sql_select += "WHERE tgltransaksi >='" + tglbefore + "' "
	sql_select += "AND tgltransaksi <='" + tglskrg + "' "
	sql_select += "AND username_client ='" + username + "' "
	sql_select += "ORDER BY tgltransaksi DESC  LIMIT 60 "

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksi_db, tgltransaksi_db                                    string
			username_client_db, resultwigo_db, nomor_db, tipebet_db, status_db string
			bet_db, win_db                                                     int
			multiplier_db                                                      float64
		)

		err = row.Scan(&idtransaksi_db, &tgltransaksi_db,
			&username_client_db, &resultwigo_db, &nomor_db, &tipebet_db,
			&bet_db, &win_db, &multiplier_db, &status_db)

		helpers.ErrorCheck(err)
		status_css := configs.STATUS_CANCEL
		if status_db == "WIN" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Invoiceclient_id = idtransaksi_db
		obj.Invoiceclient_date = tgltransaksi_db
		obj.Invoiceclient_result = resultwigo_db
		obj.Invoiceclient_username = username_client_db
		obj.Invoiceclient_nomor = nomor_db
		obj.Invoiceclient_tipebet = tipebet_db
		obj.Invoiceclient_bet = bet_db
		obj.Invoiceclient_win = win_db
		obj.Invoiceclient_multiplier = multiplier_db
		obj.Invoiceclient_status = status_db
		obj.Invoiceclient_status_css = status_css
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_result(idcompany string) (helpers.Response, error) {
	var obj entities.Model_result
	var arraobj []entities.Model_result
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, tbl_trx_transaksi, _, _ := Get_mappingdatabase(idcompany)

	tglnow, _ := goment.New()
	// tglskrg := tglnow.Format("YYYY-MM-DD HH:mm:ss")
	// tglbefore := tglnow.Add(-31, "days").Format("YYYY-MM-DD HH:mm:ss")
	tglstart := tglnow.Format("YYYY-MM-DD") + " " + " 00:00:00"
	tgltutup := tglnow.Format("YYYY-MM-DD") + " " + " 23:59:50"
	fmt.Println("tgl skrg :" + tglstart)
	fmt.Println("tgl before :" + tgltutup)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi , to_char(COALESCE(createdate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS') as datetransaksi,  "
	sql_select += "resultwigo  "
	sql_select += "FROM " + tbl_trx_transaksi + " "
	sql_select += "WHERE createdate_transaksi >='" + tglstart + "' "
	sql_select += "AND createdate_transaksi <='" + tgltutup + "' "
	sql_select += "AND resultwigo!='' "
	sql_select += "ORDER BY createdate_transaksi DESC LIMIT 100"

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksi_db, datetransaksi_db, resultwigo_db string
		)

		err = row.Scan(&idtransaksi_db, &datetransaksi_db, &resultwigo_db)

		helpers.ErrorCheck(err)

		obj.Result_invoice = idtransaksi_db
		obj.Result_date = datetransaksi_db
		obj.Result_result = resultwigo_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_transaksidetail(idcompany, idtransaksi, username, listdatabet string, total_bet int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	status := "CLOSED"

	dayendmonth := helpers.GetEndRangeDate(tglnow.Format("MM"))
	tglstart := tglnow.Format("YYYY-MM-") + "01 00:00:00"
	tglend := tglnow.Format("YYYY-MM-") + dayendmonth + " 23:59:59"
	fmt.Println("tgl start :" + tglstart)
	fmt.Println("tgl end :" + tglend)

	_, tbl_trx_transaksi, tbl_trx_transaksidetail, _ := Get_mappingdatabase(idcompany)

	status = _GetInfo_Transaksi(tbl_trx_transaksi, idtransaksi)
	if status == "OPEN" {
		type Invoicemonth struct {
			Totalbet int `json:"totalbet"`
			Totalwin int `json:"totalwin"`
		}
		type Invoicedetail struct {
			Listbet        interface{} `json:"listbet"`
			Summary        interface{} `json:"summary"`
			Totaltransaksi int         `json:"totaltransaksi"`
		}
		type Invoicedetaillistbet struct {
			Client_id         string  `json:"client_id"`
			Client_username   string  `json:"client_username"`
			Client_tipebet    string  `json:"client_tipebet"`
			Client_nomor      string  `json:"client_nomor"`
			Client_bet        int     `json:"client_bet"`
			Client_multiplier float32 `json:"client_multiplier"`
			Client_status     string  `json:"client_status"`
		}
		type Invoicesumarynomor struct {
			Nomor        string `json:"nomor"`
			Totalinvoice int    `json:"totalinvoice"`
			Totalbet     int    `json:"totalbet"`
			Totalwin     int    `json:"totalwin"`
		}

		totalbet := 0
		var objinvoicemonth Invoicemonth
		var objinvoice_parent Invoicedetail
		// var arraobjinvoice_parent []Invoicedetail
		var objinvoice_listdetail Invoicedetaillistbet
		var arraobjinvoice_listdetail []Invoicedetaillistbet
		var objinvoice_sumary Invoicesumarynomor
		var arraobjinvoice_sumary []Invoicesumarynomor
		json := []byte(listdatabet)
		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			ipaddress, _ := jsonparser.GetString(value, "ipaddress")
			mobile, _ := jsonparser.GetString(value, "mobile")
			device, _ := jsonparser.GetString(value, "device")
			tipebet, _ := jsonparser.GetString(value, "tipebet")
			nomor, _ := jsonparser.GetString(value, "nomor")
			bet, _ := jsonparser.GetInt(value, "bet")
			multiplier, _ := jsonparser.GetFloat(value, "multiplier")

			sql_insert := `
				insert into
				` + tbl_trx_transaksidetail + ` (
					idtransaksidetail, idtransaksi , username_client, ipaddress_client, browser_client, device_client, tipebet, nomor, 
					bet, multiplier, status_transaksidetail, 
					create_transaksidetail, createdate_transaksidetail  
				) values (
					$1, $2, $3, $4, $5, $6, $7, $8, 
					$9, $10, $11,   
					$12, $13      
				)
			`

			field_column := tbl_trx_transaksidetail + tglnow.Format("YYYY") + tglnow.Format("MM")
			idrecord_counter := Get_counter(field_column)
			idrecrod_value := tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)
			flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_transaksidetail, "INSERT",
				idrecrod_value, idtransaksi, username, ipaddress, mobile, device, tipebet, nomor,
				bet, multiplier, "RUNNING",
				"SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Success"

				sql_update := `
					UPDATE 
					` + tbl_trx_transaksi + `  
					SET total_bet=$1, total_member=$2, 
					update_transaksi=$3, updatedate_transaksi=$4           
					WHERE idtransaksi=$5          
				`

				flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_transaksi, "UPDATE",
					_GetTotalBet_Transaksi(tbl_trx_transaksidetail, idtransaksi),
					_GetTotalMember_Transaksi(tbl_trx_transaksidetail, idtransaksi),
					"SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"), idtransaksi)

				if flag_update {
					msg = "Success"
				} else {
					fmt.Println(msg_update)
				}
			} else {
				fmt.Println(msg_insert)
			}
			totalbet = totalbet + int(bet)

			//LIST DETAIL
			objinvoice_listdetail.Client_id = idrecrod_value
			objinvoice_listdetail.Client_username = username
			objinvoice_listdetail.Client_tipebet = tipebet
			objinvoice_listdetail.Client_nomor = nomor
			objinvoice_listdetail.Client_bet = int(bet)
			objinvoice_listdetail.Client_multiplier = float32(multiplier)
			objinvoice_listdetail.Client_status = "RUNNING"
			arraobjinvoice_listdetail = append(arraobjinvoice_listdetail, objinvoice_listdetail)

			//SUMARY
			win := int(bet) + int(float32(bet)*float32(multiplier))
			objinvoice_sumary.Nomor = nomor
			objinvoice_sumary.Totalbet = int(bet)
			objinvoice_sumary.Totalwin = int(win)
			objinvoice_sumary.Totalinvoice = 1
			arraobjinvoice_sumary = append(arraobjinvoice_sumary, objinvoice_sumary)

		})
		tglstart_redis := tglnow.Format("YYYYMM") + "01000000"
		tglend_redis := tglnow.Format("YYYYMM") + dayendmonth + "235959"

		keyredis_invoicemonth := strings.ToLower(idcompany) + "_game_12d_" + tglstart_redis + tglend_redis
		keyredis := strings.ToLower(idcompany) + "_game_12d_" + idtransaksi
		resultRD_invoice, flag_invoice := helpers.GetRedis(keyredis)
		resultRD_invoicemonth, flag_invoicemonth := helpers.GetRedis(keyredis_invoicemonth)
		if !flag_invoice {
			fmt.Println("INVOICE DATABASE")

			objinvoice_parent.Listbet = arraobjinvoice_listdetail
			objinvoice_parent.Summary = arraobjinvoice_sumary
			objinvoice_parent.Totaltransaksi = totalbet

			helpers.SetRedis(keyredis, objinvoice_parent, 60*time.Minute)
		} else {
			fmt.Println("INVOICE CACHE")

			var objinvoice_parent_RD Invoicedetail
			var objinvoice_listdetail_RD Invoicedetaillistbet
			var arraobjinvoice_listdetail_RD []Invoicedetaillistbet
			var objinvoice_sumary_RD Invoicesumarynomor
			var arraobjinvoice_sumary_RD []Invoicesumarynomor

			jsonredis := []byte(resultRD_invoice)
			totaltransaksi_RD, _ := jsonparser.GetInt(jsonredis, "totaltransaksi")
			recordlistbet_RD, _, _, _ := jsonparser.Get(jsonredis, "listbet")
			recordsummary_RD, _, _, _ := jsonparser.Get(jsonredis, "summary")
			jsonparser.ArrayEach(recordsummary_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				nomorRD, _ := jsonparser.GetString(value, "nomor")
				totalinvoiceRD, _ := jsonparser.GetInt(value, "totalinvoice")
				totalbetRD, _ := jsonparser.GetInt(value, "totalbet")
				totalwinRD, _ := jsonparser.GetInt(value, "totalwin")

				objinvoice_sumary_RD.Nomor = nomorRD
				objinvoice_sumary_RD.Totalinvoice = int(totalinvoiceRD)
				objinvoice_sumary_RD.Totalbet = int(totalbetRD)
				objinvoice_sumary_RD.Totalwin = int(totalwinRD)
				arraobjinvoice_sumary_RD = append(arraobjinvoice_sumary_RD, objinvoice_sumary_RD)

			})
			jsonparser.ArrayEach(recordlistbet_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				client_id, _ := jsonparser.GetString(value, "client_id")
				client_username, _ := jsonparser.GetString(value, "client_username")
				client_tipebet, _ := jsonparser.GetString(value, "client_tipebet")
				client_nomor, _ := jsonparser.GetString(value, "client_nomor")
				client_bet, _ := jsonparser.GetInt(value, "client_bet")
				client_multiplier, _ := jsonparser.GetFloat(value, "client_multiplier")
				client_status, _ := jsonparser.GetString(value, "client_status")

				objinvoice_listdetail_RD.Client_id = client_id
				objinvoice_listdetail_RD.Client_username = client_username
				objinvoice_listdetail_RD.Client_tipebet = client_tipebet
				objinvoice_listdetail_RD.Client_nomor = client_nomor
				objinvoice_listdetail_RD.Client_bet = int(client_bet)
				objinvoice_listdetail_RD.Client_multiplier = float32(client_multiplier)
				objinvoice_listdetail_RD.Client_status = client_status
				arraobjinvoice_listdetail_RD = append(arraobjinvoice_listdetail_RD, objinvoice_listdetail_RD)

			})

			for i := 0; i < len(arraobjinvoice_listdetail); i++ { // data diatas listdetailbet
				client_id := arraobjinvoice_listdetail[i].Client_id
				client_username := arraobjinvoice_listdetail[i].Client_username
				client_tipebet := arraobjinvoice_listdetail[i].Client_tipebet
				client_nomor := arraobjinvoice_listdetail[i].Client_nomor
				client_bet := arraobjinvoice_listdetail[i].Client_bet
				client_multiplier := arraobjinvoice_listdetail[i].Client_multiplier
				client_status := arraobjinvoice_listdetail[i].Client_status

				objinvoice_listdetail_RD.Client_id = client_id
				objinvoice_listdetail_RD.Client_username = client_username
				objinvoice_listdetail_RD.Client_tipebet = client_tipebet
				objinvoice_listdetail_RD.Client_nomor = client_nomor
				objinvoice_listdetail_RD.Client_bet = client_bet
				objinvoice_listdetail_RD.Client_multiplier = client_multiplier
				objinvoice_listdetail_RD.Client_status = client_status
				arraobjinvoice_listdetail_RD = append(arraobjinvoice_listdetail_RD, objinvoice_listdetail_RD)
			}

			for i := 0; i < len(arraobjinvoice_sumary); i++ { // data diatas
				nomor_loop := arraobjinvoice_sumary[i].Nomor
				invoice_loop := arraobjinvoice_sumary[i].Totalinvoice
				bet_loop := arraobjinvoice_sumary[i].Totalbet
				win_loop := arraobjinvoice_sumary[i].Totalwin
				flag_insert := true

				jsonparser.ArrayEach(recordsummary_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					nomorRD, _ := jsonparser.GetString(value, "nomor")
					totalinvoiceRD, _ := jsonparser.GetInt(value, "totalinvoice")
					totalbetRD, _ := jsonparser.GetInt(value, "totalbet")
					totalwinRD, _ := jsonparser.GetInt(value, "totalwin")

					totalinvoice_temp := 0
					totalbet_temp := 0
					totalwin_temp := 0

					if nomor_loop == nomorRD {
						totalinvoice_temp = int(totalinvoiceRD) + invoice_loop
						totalbet_temp = bet_loop + int(totalbetRD)
						totalwin_temp = win_loop + int(totalwinRD)
						for j := 0; j < len(arraobjinvoice_sumary_RD); j++ {
							if arraobjinvoice_sumary_RD[j].Nomor == nomor_loop {
								arraobjinvoice_sumary_RD[j].Totalinvoice = totalinvoice_temp
								arraobjinvoice_sumary_RD[j].Totalbet = totalbet_temp
								arraobjinvoice_sumary_RD[j].Totalwin = totalwin_temp
							}
						}
						flag_insert = false
					}
				})
				if flag_insert {
					objinvoice_sumary_RD.Nomor = nomor_loop
					objinvoice_sumary_RD.Totalinvoice = int(invoice_loop)
					objinvoice_sumary_RD.Totalbet = int(bet_loop)
					objinvoice_sumary_RD.Totalwin = int(win_loop)
					arraobjinvoice_sumary_RD = append(arraobjinvoice_sumary_RD, objinvoice_sumary_RD)
				}
			}
			totalbetnew := totalbet + int(totaltransaksi_RD)

			objinvoice_parent_RD.Listbet = arraobjinvoice_listdetail_RD
			objinvoice_parent_RD.Summary = arraobjinvoice_sumary_RD
			objinvoice_parent_RD.Totaltransaksi = totalbetnew

			helpers.SetRedis(keyredis, objinvoice_parent_RD, 60*time.Minute)
		}

		if !flag_invoicemonth {
			fmt.Println("INVOICE MONTH DATABASE")
			totalbet_DB, totalwin_DB := _GetTotalBet_ByDate(tbl_trx_transaksi, tglstart, tglend)

			objinvoicemonth.Totalbet = totalbet_DB
			objinvoicemonth.Totalwin = totalwin_DB

			helpers.SetRedis(keyredis_invoicemonth, objinvoicemonth, 0)
		} else {
			fmt.Println("INVOICE MONTH CACHE")
			jsonredis := []byte(resultRD_invoicemonth)
			totalbet_RD, _ := jsonparser.GetInt(jsonredis, "totalbet")
			totalwin_RD, _ := jsonparser.GetInt(jsonredis, "totalwin")

			totalbetnew_month := totalbet + int(totalbet_RD)

			objinvoicemonth.Totalbet = totalbetnew_month
			objinvoicemonth.Totalwin = int(totalwin_RD)

			helpers.SetRedis(keyredis_invoicemonth, objinvoicemonth, 0)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _GetTotalBet_ByDate(table, startdate, enddate string) (int, int) {
	con := db.CreateCon()
	ctx := context.Background()
	totalbet := 0
	totalwin := 0
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "COALESCE(SUM(total_bet),0) AS totalbet,  COALESCE(sum(total_win),0) as totalwin "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE createdate_transaksi >='" + startdate + "'   "
	sql_select += "AND createdate_transaksi <='" + enddate + "'   "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&totalbet, &totalwin); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return totalbet, totalwin
}
func _GetTotalBet_Transaksi(table, idtransaksi string) int {
	con := db.CreateCon()
	ctx := context.Background()
	total_bet := 0
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "SUM(bet) AS total_bet "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE idtransaksi='" + idtransaksi + "'   "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&total_bet); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return total_bet
}
func _GetTotalMember_Transaksi(table, idtransaksi string) int {
	con := db.CreateCon()
	ctx := context.Background()
	total_member := 0
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "count(distinct(username_client)) as totalmember "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE idtransaksi='" + idtransaksi + "'   "
	sql_select += "group by username_client "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&total_member); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return total_member
}
func _GetInfo_Transaksi(table, idtransaksi string) string {
	con := db.CreateCon()
	ctx := context.Background()
	status := "CLOSED"

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "status_transaksi "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE idtransaksi='" + idtransaksi + "'   "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&status); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return status
}
func GetInfo_CompanyConf(idcompany string) (float64, float64, float64, string, string) {
	con := db.CreateCon()
	ctx := context.Background()
	win_angka := 0.0
	win_redblack := 0.0
	win_line := 0.0
	status_redblack := "N"
	status_maintenance := "N"

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "conf_2digit_30_win, conf_2digit_30_win_redblack,conf_2digit_30_win_line,  "
	sql_select += "conf_2digit_30_status_redblack_line,conf_2digit_30_maintenance  "
	sql_select += "FROM " + configs.DB_tbl_mst_company_config + " "
	sql_select += "WHERE idcompany='" + idcompany + "'   "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&win_angka, &win_redblack, &win_line, &status_redblack, &status_maintenance); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return win_angka, win_redblack, win_line, status_redblack, status_maintenance
}
