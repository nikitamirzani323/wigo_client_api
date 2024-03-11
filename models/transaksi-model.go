package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/configs"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/db"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/entities"
	"github.com/nikitamirzani323/BTANGKAS_CLIENT_API/helpers"
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
	fmt.Println("tgl skrg :" + tglskrg)
	fmt.Println("tgl before :" + tglbefore)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi, to_char(COALESCE(tgltransaksi,now()), 'YYYY-MM-DD HH24:MI:SS') as tgltransaksi,  "
	sql_select += "username_client, resultwigo, nomor,   "
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
			idtransaksi_db, tgltransaksi_db                        string
			username_client_db, resultwigo_db, nomor_db, status_db string
			bet_db, win_db                                         int
			multiplier_db                                          float64
		)

		err = row.Scan(&idtransaksi_db, &tgltransaksi_db,
			&username_client_db, &resultwigo_db, &nomor_db,
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
func Save_transaksi(idcompany, idcurr string) (helpers.Responsetransaksi, error) {
	var res helpers.Responsetransaksi
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	_, tbl_trx_transaksi, _, _ := Get_mappingdatabase(idcompany)
	sql_insert := `
			insert into
			` + tbl_trx_transaksi + ` (
				idtransaksi , idcurr, idcompany, datetransaksi,
				create_transaksi, createdate_transaksi 
			) values (
				$1, $2, $3, $4, 
				$5, $6 
			)
		`

	field_column := tbl_trx_transaksi + tglnow.Format("YYYY") + tglnow.Format("MM")
	idrecord_counter := Get_counter(field_column)
	idrecrodparent_value := strings.ToUpper(idcompany) + "-" + tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)
	date_transaksi := tglnow.Format("YYYY-MM-DD HH:mm:ss")

	flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_transaksi, "INSERT",
		idrecrodparent_value, idcurr, idcompany, date_transaksi,
		"SYSTEM", date_transaksi)

	if flag_insert {
		msg = "Succes"

	} else {
		fmt.Println(msg_insert)
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Idtransaksi = idrecrodparent_value
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Save_transaksidetail(idcompany, idtransaksi, username, listdatabet string, total_bet int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	status := "CLOSED"

	_, tbl_trx_transaksi, tbl_trx_transaksidetail, _ := Get_mappingdatabase(idcompany)

	status = _GetInfo_Transaksi(tbl_trx_transaksi, idtransaksi)
	if status == "OPEN" {
		json := []byte(listdatabet)
		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			tipebet, _ := jsonparser.GetString(value, "tipebet")
			nomor, _ := jsonparser.GetString(value, "nomor")
			bet, _ := jsonparser.GetInt(value, "bet")
			multiplier, _ := jsonparser.GetFloat(value, "multiplier")

			sql_insert := `
				insert into
				` + tbl_trx_transaksidetail + ` (
					idtransaksidetail, idtransaksi , username_client, tipebet, nomor, 
					bet, multiplier, status_transaksidetail, 
					create_transaksidetail, createdate_transaksidetail  
				) values (
					$1, $2, $3, $4, $5,
					$6, $7, $8,   
					$9, $10      
				)
			`

			field_column := tbl_trx_transaksidetail + tglnow.Format("YYYY") + tglnow.Format("MM")
			idrecord_counter := Get_counter(field_column)
			idrecrod_value := tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)
			flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_transaksidetail, "INSERT",
				idrecrod_value, idtransaksi, username, tipebet, nomor,
				bet, multiplier, "RUNNING",
				"SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"

				sql_update := `
				UPDATE 
				` + tbl_trx_transaksi + `  
				SET total_bet=$1,
				update_transaksi=$2, updatedate_transaksi=$3          
				WHERE idtransaksi=$4         
			`

				flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_transaksi, "UPDATE",
					_GetTotalBet_Transaksi(tbl_trx_transaksidetail, idtransaksi),
					"SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"), idtransaksi)

				if flag_update {
					msg = "Succes"
				} else {
					fmt.Println(msg_update)
				}
			} else {
				fmt.Println(msg_insert)
			}
		})

	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
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
