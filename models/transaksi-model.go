package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/WIGO_AGEN_API/configs"
	"github.com/nikitamirzani323/WIGO_AGEN_API/db"
	"github.com/nikitamirzani323/WIGO_AGEN_API/entities"
	"github.com/nikitamirzani323/WIGO_AGEN_API/helpers"
	"github.com/nleeper/goment"
)

func Fetch_transaksi2D30SHome(idcompany, idinvoice string) (helpers.ResponseTransaksi2D30S, error) {
	var obj entities.Model_transaksi2D30S
	var arraobj []entities.Model_transaksi2D30S
	var res helpers.ResponseTransaksi2D30S
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()
	tglnow, _ := goment.New()
	month := tglnow.Format("MMM")
	day := helpers.GetEndRangeDate(strings.ToUpper(month))
	startdate := tglnow.Format("YYYY-MM") + "-01 00:00:00"
	enddate := tglnow.Format("YYYY-MM-") + day + " 23:59:00"
	periode := strings.ToUpper(month) + "-" + tglnow.Format("YYYY")

	_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi,idcurr,datetransaksi,  "
	sql_select += "resultwigo,total_bet,total_win,  "
	sql_select += "status_transaksi, "
	sql_select += "create_transaksi, to_char(COALESCE(createdate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "update_transaksi, to_char(COALESCE(updatedate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select += "FROM " + tbl_trx_transaksi + "   "
	sql_select += "WHERE LOWER(idcompany)='" + strings.ToLower(idcompany) + "' "
	if idinvoice != "" {
		sql_select += "AND idtransaksi='" + idinvoice + "' "
	} else {
		sql_select += "AND createdate_transaksi >='" + startdate + "' "
		sql_select += "AND createdate_transaksi <='" + enddate + "' "
	}
	sql_select += "ORDER BY datetransaksi DESC    "
	log.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksi_db, idcurr_db, datetransaksi_db, resultwigo_db, status_transaksi_db            string
			total_bet_db, total_win_db                                                                 int
			create_transaksi_db, createdate_transaksi_db, update_transaksi_db, updatedate_transaksi_db string
		)

		err = row.Scan(&idtransaksi_db, &idcurr_db, &datetransaksi_db,
			&resultwigo_db, &total_bet_db, &total_win_db, &status_transaksi_db,
			&create_transaksi_db, &createdate_transaksi_db, &update_transaksi_db, &updatedate_transaksi_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_transaksi_db != "" {
			create = create_transaksi_db + ", " + createdate_transaksi_db
		}
		if update_transaksi_db != "" {
			update = update_transaksi_db + ", " + updatedate_transaksi_db
		}
		if status_transaksi_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Transaksi2D30S_id = idtransaksi_db
		obj.Transaksi2D30S_idcurr = idcurr_db
		obj.Transaksi2D30S_date = datetransaksi_db
		obj.Transaksi2D30S_result = resultwigo_db
		obj.Transaksi2D30S_totalbet = total_bet_db
		obj.Transaksi2D30S_totalwin = total_win_db
		obj.Transaksi2D30S_winlose = total_bet_db - total_win_db
		obj.Transaksi2D30S_status = status_transaksi_db
		obj.Transaksi2D30S_status_css = status_css
		obj.Transaksi2D30S_create = create
		obj.Transaksi2D30S_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Periode = periode
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_transaksi2D30SPrediksi(idcompany, idinvoice, result string) (helpers.ResponseTransaksi2D30SPrediksi, error) {
	var obj entities.Model_transaksi2D30SPrediksi
	var arraobj []entities.Model_transaksi2D30SPrediksi
	var res helpers.ResponseTransaksi2D30SPrediksi
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, _, tbl_trx_transaksidetail := Get_mappingdatabase(idcompany)

	total_bet := 0
	total_win := 0
	winlose := 0
	sql_select_detail := `SELECT 
					idtransaksidetail , nomor, bet, multiplier, username_client,
					to_char(COALESCE(createdate_transaksidetail,now()), 'YYYY-MM-DD HH24:MI:SS') 
					FROM ` + tbl_trx_transaksidetail + `  
					WHERE status_transaksidetail='RUNNING'  
					AND idtransaksi='` + idinvoice + `'  `

	row, err := con.QueryContext(ctx, sql_select_detail)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			bet_db                                                                            int
			multiplier_db                                                                     float64
			idtransaksidetail_db, nomor_db, username_client_db, createdate_transaksidetail_db string
		)

		err = row.Scan(&idtransaksidetail_db, &nomor_db, &bet_db,
			&multiplier_db, &username_client_db, &createdate_transaksidetail_db)
		helpers.ErrorCheck(err)

		total_bet = total_bet + bet_db

		status_client := _rumuswigo2D30S(nomor_db, result)

		win := 0
		if status_client == "WIN" {
			win = bet_db + int(float64(bet_db)*multiplier_db)
			total_win = total_win + win

			status_css := configs.STATUS_CANCEL
			if status_client == "WIN" {
				status_css = configs.STATUS_COMPLETE
			}
			obj.Transaksi2D30Sprediksi_id = idtransaksidetail_db
			obj.Transaksi2D30Sprediksi_date = createdate_transaksidetail_db
			obj.Transaksi2D30Sprediksi_username = username_client_db
			obj.Transaksi2D30Sprediksi_nomor = nomor_db
			obj.Transaksi2D30Sprediksi_bet = bet_db
			obj.Transaksi2D30Sprediksi_win = win
			obj.Transaksi2D30Sprediksi_winlose = bet_db - win
			obj.Transaksi2D30Sprediksi_status = status_client
			obj.Transaksi2D30Sprediksi_status_css = status_css
			arraobj = append(arraobj, obj)
		}
		total_bet = total_bet + bet_db

	}
	defer row.Close()
	winlose = total_bet - total_win

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.TotalMember = 0
	res.Totalbet = total_bet
	res.Totalwin = total_win
	res.Winlose = winlose
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_updateresult2D30S(admin, idrecord, idcompany, result string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	const invoice_client_redis = "CLIENT_LISTINVOICE"
	const invoice_result_redis = "CLIENT_RESULT"

	_, tbl_trx_transaksi, tbl_trx_transaksidetail := Get_mappingdatabase(idcompany)

	result_db, status_db := _GetInvoiceInfo(idrecord, idcompany)
	if result_db == "" && status_db == "OPEN" {
		sql_update := `
			UPDATE 
			` + tbl_trx_transaksi + `  
				SET resultwigo=$1, status_transaksi=$2, 
				update_transaksi=$3, updatedate_transaksi=$4           
				WHERE idtransaksi=$5        
		`

		flag_update, msg_update := Exec_SQL(sql_update, tbl_trx_transaksi, "UPDATE",
			result, "CLOSED",
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"

			con := db.CreateCon()
			ctx := context.Background()
			flag_detail := false
			sql_select_detail := `SELECT 
					idtransaksidetail , nomor, bet, multiplier, username_client 
					FROM ` + tbl_trx_transaksidetail + `  
					WHERE status_transaksidetail='RUNNING'  
					AND idtransaksi='` + idrecord + `'  `

			row, err := con.QueryContext(ctx, sql_select_detail)
			helpers.ErrorCheck(err)
			for row.Next() {
				var (
					bet_db                                             int
					multiplier_db                                      float64
					idtransaksidetail_db, nomor_db, username_client_db string
				)

				err = row.Scan(&idtransaksidetail_db, &nomor_db, &bet_db, &multiplier_db, &username_client_db)
				helpers.ErrorCheck(err)

				status_client := _rumuswigo2D30S(nomor_db, result)
				win := 0
				if status_client == "WIN" {
					win = bet_db + int(float64(bet_db)*multiplier_db)
				}

				// UPDATE STATUS DETAIL
				sql_update_detail := `
					UPDATE 
					` + tbl_trx_transaksidetail + `  
					SET status_transaksidetail=$1, win=$2, 
					update_transaksidetail=$3, updatedate_transaksidetail=$4           
					WHERE idtransaksidetail=$5          
				`
				flag_update_detail, msg_update_detail := Exec_SQL(sql_update_detail, tbl_trx_transaksidetail, "UPDATE",
					status_client, win,
					"SYSTEM", tglnow.Format("YYYY-MM-DD HH:mm:ss"), idtransaksidetail_db)

				if !flag_update_detail {
					fmt.Println(msg_update_detail)
				}
				flag_detail = true

				key_redis_invoice_client := invoice_client_redis + "_" + strings.ToLower(idcompany) + "_" + strings.ToLower(username_client_db)
				val_invoice_client := helpers.DeleteRedis(key_redis_invoice_client)
				fmt.Println("")
				fmt.Printf("Redis Delete INVOICE : %d - %s \r", val_invoice_client, key_redis_invoice_client)
				fmt.Println("")
			}
			defer row.Close()
			if flag_detail {
				// UPDATE PARENT
				total_bet, total_win := _GetTotalBetWin_Transaksi(tbl_trx_transaksidetail, idrecord)
				sql_update_parent := `
					UPDATE 
					` + tbl_trx_transaksi + `  
					SET total_bet=$1, total_win=$2, 
					update_transaksi=$3, updatedate_transaksi=$4           
					WHERE idtransaksi=$5       
				`
				flag_update_parent, msg_update_parent := Exec_SQL(sql_update_parent, tbl_trx_transaksi, "UPDATE",
					total_bet, total_win,
					admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

				if !flag_update_parent {
					fmt.Println(msg_update_parent)

				}
			}
			key_redis_result := invoice_result_redis + "_" + strings.ToLower(idcompany)
			val_result := helpers.DeleteRedis(key_redis_result)
			fmt.Println("")
			fmt.Printf("Redis Delete RESULT : %d - %s \r", val_result, key_redis_result)
			fmt.Println("")
			idcurr := _GetCompanyInfo(idcompany)
			invoice := _Generate_incoive(strings.ToLower(idcompany), idcurr)

			fieldconfig_redis := "TIMER_" + strings.ToLower(idcompany)
			type Configure struct {
				Time    int    `json:"time"`
				Invoice string `json:"invoice"`
			}
			var obj Configure
			obj.Time = _GetCompanyConfInfo(idcompany)
			obj.Invoice = invoice
			helpers.SetRedis(fieldconfig_redis, obj, 5*time.Minute)
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _Generate_incoive(idcompany, idcurr string) string {
	tglnow, _ := goment.New()
	id_invoice := _GetInvoice(idcompany)
	if id_invoice == "" {
		_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

		sql_insert := `
			insert into
			` + tbl_trx_transaksi + ` (
				idtransaksi , idcurr, idcompany, datetransaksi, status_transaksi, 
				create_transaksi, createdate_transaksi 
			) values (
				$1, $2, $3, $4, $5,  
				$6, $7  
			)
		`

		field_column := tbl_trx_transaksi + tglnow.Format("YYYY") + tglnow.Format("MM") + tglnow.Format("DD")
		idrecord_counter := Get_counter(field_column)
		idrecrodparent_value := tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)
		date_transaksi := tglnow.Format("YYYY-MM-DD HH:mm:ss")

		flag_insert, msg_insert := Exec_SQL(sql_insert, tbl_trx_transaksi, "INSERT",
			idrecrodparent_value, idcurr, idcompany, date_transaksi, "OPEN",
			"SYSTEM", date_transaksi)

		if flag_insert {
			id_invoice = idrecrodparent_value

		} else {
			fmt.Println(msg_insert)
		}
	}

	return id_invoice
}
func _GetInvoice(idcompany string) string {
	con := db.CreateCon()
	ctx := context.Background()

	_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

	idtransaksi := ""

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi "
	sql_select += "FROM " + tbl_trx_transaksi + " "
	sql_select += "WHERE resultwigo='' "
	sql_select += "AND status_transaksi='OPEN' "
	sql_select += "ORDER BY idtransaksi DESC LIMIT 1"

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&idtransaksi); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return idtransaksi
}
func _GetInvoiceInfo(invoice, idcompany string) (string, string) {
	con := db.CreateCon()
	ctx := context.Background()

	_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

	result := ""
	status := ""

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "resultwigo, status_transaksi "
	sql_select += "FROM " + tbl_trx_transaksi + " "
	sql_select += "WHERE idtransaksi='" + invoice + "' "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&result, &status); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return result, status
}
func _GetCompanyInfo(idcompany string) string {
	con := db.CreateCon()
	ctx := context.Background()

	idcurr := ""

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idcurr "
	sql_select += "FROM " + configs.DB_tbl_mst_company + " "
	sql_select += "WHERE idcompany='" + idcompany + "' "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&idcurr); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return idcurr
}
func _GetCompanyConfInfo(idcompany string) int {
	con := db.CreateCon()
	ctx := context.Background()

	conf_2digit_30_time := 0

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "conf_2digit_30_time "
	sql_select += "FROM " + configs.DB_tbl_mst_company_config + " "
	sql_select += "WHERE idcompany='" + idcompany + "' "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&conf_2digit_30_time); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return conf_2digit_30_time
}
func _GetTotalBetWin_Transaksi(table, idtransaksi string) (int, int) {
	con := db.CreateCon()
	ctx := context.Background()
	total_bet := 0
	total_win := 0
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "SUM(bet) AS total_bet, SUM(win) AS total_win  "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE idtransaksi='" + idtransaksi + "'   "
	sql_select += "AND status_transaksidetail !='RUNNING'   "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&total_bet, &total_win); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return total_bet, total_win
}
func _rumuswigo2D30S(nomorclient, nomorkeluaran string) string {
	result := "LOSE"
	if nomorclient == nomorkeluaran {
		result = "WIN"
	}
	return result
}