package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/wigo_agen_api/configs"
	"bitbucket.org/isbtotogroup/wigo_agen_api/db"
	"bitbucket.org/isbtotogroup/wigo_agen_api/entities"
	"bitbucket.org/isbtotogroup/wigo_agen_api/helpers"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_transaksi2D30SHome(idcompany, idinvoice, search string, page int) (helpers.ResponseTransaksi2D30S, error) {
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

	perpage := configs.PAGING_PAGE
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idtransaksi) as totalpurchase  "
	sql_selectcount += "FROM " + tbl_trx_transaksi + "  "
	sql_selectcount += "WHERE LOWER(idcompany)='" + strings.ToLower(idcompany) + "' "
	if search != "" {
		sql_selectcount += "WHERE idtransaksi LIKE '%" + strings.ToLower(search) + "%' "
	} else {
		if idinvoice != "" {
			sql_selectcount += "AND idtransaksi='" + idinvoice + "' "
		} else {
			sql_selectcount += "AND createdate_transaksi >='" + startdate + "' "
			sql_selectcount += "AND createdate_transaksi <='" + enddate + "' "
		}
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi,idcurr,to_char(COALESCE(createdate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS'),  "
	sql_select += "resultwigo,total_member,total_bet,total_win,  "
	sql_select += "status_transaksi, "
	sql_select += "create_transaksi, to_char(COALESCE(createdate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "update_transaksi, to_char(COALESCE(updatedate_transaksi,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select += "FROM " + tbl_trx_transaksi + "   "
	sql_select += "WHERE LOWER(idcompany)='" + strings.ToLower(idcompany) + "' "
	if search != "" {
		sql_select += "WHERE idtransaksi LIKE '%" + strings.ToLower(search) + "%' "
	} else {
		if idinvoice != "" {
			sql_select += "AND idtransaksi='" + idinvoice + "' "
		} else {
			sql_select += "AND createdate_transaksi >='" + startdate + "' "
			sql_select += "AND createdate_transaksi <='" + enddate + "' "
		}
	}
	sql_select += "ORDER BY createdate_transaksi DESC   OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksi_db, idcurr_db, createdate_transaksi, resultwigo_db, status_transaksi_db        string
			total_member_db, total_bet_db, total_win_db                                                int
			create_transaksi_db, createdate_transaksi_db, update_transaksi_db, updatedate_transaksi_db string
		)

		err = row.Scan(&idtransaksi_db, &idcurr_db, &createdate_transaksi,
			&resultwigo_db, &total_member_db, &total_bet_db, &total_win_db, &status_transaksi_db,
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
		obj.Transaksi2D30S_date = createdate_transaksi
		obj.Transaksi2D30S_result = resultwigo_db
		obj.Transaksi2D30S_totalmember = total_member_db
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

	//FORMAT : nuke:12D30S:invoicemonth_2024040100000020240430235959
	dayendmonth := helpers.GetEndRangeDateTwo(tglnow.Format("MM"))
	tglstart_redis := tglnow.Format("YYYYMM") + "01000000"
	tglend_redis := tglnow.Format("YYYYMM") + dayendmonth + "235959"
	keyinvoicemonth_redis := strings.ToLower(idcompany) + ":12D30S:invoicemonth_" + tglstart_redis + tglend_redis
	log.Println(keyinvoicemonth_redis)
	invoicemonth_redis, flag_invoicemonth := helpers.GetRedis(keyinvoicemonth_redis)

	total_bet := 0
	total_win := 0
	var winlose_agen float64 = 0
	var winlose_member float64 = 0
	if idinvoice == "" {
		if !flag_invoicemonth {
			fmt.Println("INVOICE MONTHLY DATABASE")
			total_bet, total_win = _GetTotalBetWinByDate_Transaksi(tbl_trx_transaksi, startdate, enddate)
		} else {
			fmt.Println("INVOICE MONTHLY CACHE")
			jsonredis := []byte(invoicemonth_redis)
			totalbet_RD, _ := jsonparser.GetInt(jsonredis, "totalbet")
			totalwin_RD, _ := jsonparser.GetInt(jsonredis, "totalwin")
			total_bet = int(totalbet_RD)
			total_win = int(totalwin_RD)
		}

		winlose_agen = float64(total_bet - total_win)
		if winlose_agen < 0 {
			winlose_member = math.Abs(winlose_agen)
		} else {
			winlose_member = -winlose_agen
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Periode = periode
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.TotalBet = total_bet
	res.TotalWin = total_win
	res.Winlose_agen = int(winlose_agen)
	res.Winlose_member = int(winlose_member)
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_transaksi2D30SInfo(idcompany, idinvoice string) (helpers.ResponseTransaksi2D30SInfo, error) {
	var obj entities.Model_transaksi2D30SInfoInvoice
	var arraobj []entities.Model_transaksi2D30SInfoInvoice
	var res helpers.ResponseTransaksi2D30SInfo
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, tbl_trx_transaksi, tbl_trx_transaksidetail := Get_mappingdatabase(idcompany)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksi,  "
	sql_select += "resultwigo,total_member,total_bet,total_win,status_transaksi "
	sql_select += "FROM " + tbl_trx_transaksi + "   "
	sql_select += "WHERE LOWER(idcompany)='" + strings.ToLower(idcompany) + "' "
	sql_select += "AND idtransaksi='" + idinvoice + "' "

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksi_db, resultwigo_db, status_transaksi_db string
			total_member_db, total_bet_db, total_win_db        int
		)

		err = row.Scan(&idtransaksi_db, &resultwigo_db, &total_member_db, &total_bet_db, &total_win_db, &status_transaksi_db)

		helpers.ErrorCheck(err)

		obj.Transaksi2D30Sinfo_id = idtransaksi_db
		obj.Transaksi2D30Sinfo_result = resultwigo_db
		obj.Transaksi2D30Sinfo_totalmember = total_member_db
		obj.Transaksi2D30Sinfo_totalbet = total_bet_db
		obj.Transaksi2D30Sinfo_totalwin = total_win_db
		obj.Transaksi2D30Sinfo_winlose = total_bet_db - total_win_db
		obj.Transaksi2D30Sinfo_status = status_transaksi_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objdetail entities.Model_transaksi2D30Ssummary
	var arraobjdetail []entities.Model_transaksi2D30Ssummary

	sql_selectdetail := ""
	sql_selectdetail += "SELECT "
	sql_selectdetail += "nomor, count(idtransaksidetail) as totalinvoice,  "
	sql_selectdetail += "sum(bet) as totalbet, sum(bet+(bet*multiplier)) as totalwin  "
	sql_selectdetail += "FROM " + tbl_trx_transaksidetail + "   "
	sql_selectdetail += "WHERE idtransaksi='" + idinvoice + "' "
	sql_selectdetail += "group by nomor  "
	sql_selectdetail += "order by totalwin desc  "

	row_detail, err_detail := con.QueryContext(ctx, sql_selectdetail)
	helpers.ErrorCheck(err_detail)
	for row_detail.Next() {
		var (
			nomor_db                                  string
			totalinvoice_db, totalbet_db, totalwin_db int
		)

		err = row_detail.Scan(&nomor_db, &totalinvoice_db, &totalbet_db, &totalwin_db)

		helpers.ErrorCheck(err)

		objdetail.Transaksi2D30Ssummary_nomor = nomor_db
		objdetail.Transaksi2D30Ssummary_totalinvoice = totalinvoice_db
		objdetail.Transaksi2D30Ssummary_totalbet = totalbet_db
		objdetail.Transaksi2D30Ssummary_totalwin = totalwin_db
		arraobjdetail = append(arraobjdetail, objdetail)
		msg = "Success"
	}
	defer row_detail.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Summary = arraobjdetail
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_transaksi2D30SDetail(idcompany, idtransaksi, status string) (helpers.Response, error) {
	var obj entities.Model_transaksi2D30Sdetail
	var arraobj []entities.Model_transaksi2D30Sdetail
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	_, _, tbl_trx_transaksidetail := Get_mappingdatabase(idcompany)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idtransaksidetail,  "
	sql_select += "username_client,ipaddress_client,browser_client,device_client,  "
	sql_select += "nomor,tipebet,bet,win,multiplier,  "
	sql_select += "status_transaksidetail, "
	sql_select += "create_transaksidetail, to_char(COALESCE(createdate_transaksidetail,now()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "update_transaksidetail, to_char(COALESCE(updatedate_transaksidetail,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select += "FROM " + tbl_trx_transaksidetail + "   "
	if idtransaksi != "" {
		sql_select += "WHERE idtransaksi='" + idtransaksi + "' "
		sql_select += "AND status_transaksidetail='" + status + "' "
	} else {
		sql_select += "WHERE status_transaksidetail='" + status + "' "
	}

	sql_select += "ORDER BY createdate_transaksidetail DESC    "

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtransaksidetail_db, nomor_db, tipebet_db, status_transaksidetail_db                                              string
			username_client_db, ipaddress_client_db, browser_client_db, device_client_db                                       string
			bet_db, win_db                                                                                                     int
			multiplier_db                                                                                                      float64
			create_transaksidetail_db, createdate_transaksidetail_db, update_transaksidetail_db, updatedate_transaksidetail_db string
		)

		err = row.Scan(&idtransaksidetail_db, &username_client_db, &ipaddress_client_db, &browser_client_db, &device_client_db,
			&nomor_db, &tipebet_db, &bet_db, &win_db, &multiplier_db, &status_transaksidetail_db,
			&create_transaksidetail_db, &createdate_transaksidetail_db, &update_transaksidetail_db, &updatedate_transaksidetail_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if create_transaksidetail_db != "" {
			create = create_transaksidetail_db + ", " + createdate_transaksidetail_db
		}
		if update_transaksidetail_db != "" {
			update = update_transaksidetail_db + ", " + updatedate_transaksidetail_db
		}
		if status_transaksidetail_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Transaksi2D30Sdetail_id = idtransaksidetail_db
		obj.Transaksi2D30Sdetail_date = createdate_transaksidetail_db
		obj.Transaksi2D30Sdetail_ipaddress = ipaddress_client_db
		obj.Transaksi2D30Sdetail_browser = browser_client_db
		obj.Transaksi2D30Sdetail_device = device_client_db
		obj.Transaksi2D30Sdetail_username = username_client_db
		obj.Transaksi2D30Sdetail_tipebet = tipebet_db
		obj.Transaksi2D30Sdetail_nomor = nomor_db
		obj.Transaksi2D30Sdetail_bet = bet_db
		obj.Transaksi2D30Sdetail_win = win_db
		obj.Transaksi2D30Sdetail_multiplier = multiplier_db
		obj.Transaksi2D30Sdetail_status = status_transaksidetail_db
		obj.Transaksi2D30Sdetail_status_css = status_css
		obj.Transaksi2D30Sdetail_create = create
		obj.Transaksi2D30Sdetail_update = update
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
	var winlose_agen float64 = 0
	var winlose_member float64 = 0
	sql_select_detail := `SELECT 
					idtransaksidetail , nomor, tipebet, bet, multiplier, username_client,
					to_char(COALESCE(createdate_transaksidetail,now()), 'YYYY-MM-DD HH24:MI:SS') 
					FROM ` + tbl_trx_transaksidetail + `  
					WHERE status_transaksidetail='RUNNING'  
					AND idtransaksi='` + idinvoice + `'  
					`

	row, err := con.QueryContext(ctx, sql_select_detail)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			bet_db                                                                                        int
			multiplier_db                                                                                 float64
			idtransaksidetail_db, nomor_db, tipebet_db, username_client_db, createdate_transaksidetail_db string
		)

		err = row.Scan(&idtransaksidetail_db, &nomor_db, &tipebet_db, &bet_db,
			&multiplier_db, &username_client_db, &createdate_transaksidetail_db)
		helpers.ErrorCheck(err)
		total_bet = total_bet + bet_db
		status_client := _rumuswigo2D30S(tipebet_db, nomor_db, result)

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
			obj.Transaksi2D30Sprediksi_multiplier = multiplier_db
			obj.Transaksi2D30Sprediksi_win = win
			obj.Transaksi2D30Sprediksi_winlose = bet_db - win
			obj.Transaksi2D30Sprediksi_status = status_client
			obj.Transaksi2D30Sprediksi_status_css = status_css
			arraobj = append(arraobj, obj)
		}

	}
	defer row.Close()

	winlose_agen = float64(total_bet - total_win)
	if winlose_agen < 0 {
		winlose_member = math.Abs(winlose_agen)
	} else {
		winlose_member = -winlose_agen
	}
	_, _, total_member := _GetInvoiceInfo(idinvoice, idcompany)

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.TotalMember = total_member
	res.Totalbet = total_bet
	res.Totalwin = total_win
	res.Winlose_agen = int(winlose_agen)
	res.Winlose_member = int(winlose_member)
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_Agenconf(idcompany string) (helpers.Response, error) {
	var obj entities.Model_agenconf
	var arraobj []entities.Model_agenconf
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "conf_2digit_30_time, "
	sql_select += "conf_2digit_30_win,conf_2digit_30_win_redblack,conf_2digit_30_win_line, conf_2digit_30_win_zona, conf_2digit_30_win_jackpot,  "
	sql_select += "conf_2digit_30_operator  "
	sql_select += "FROM " + configs.DB_tbl_mst_company_config + " "
	sql_select += "WHERE idcompany ='" + idcompany + "' "
	fmt.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			conf_2digit_30_time_db                                                                                                                       int
			conf_2digit_30_win_db, conf_2digit_30_win_redblack_db, conf_2digit_30_win_line_db, conf_2digit_30_win_zona_db, conf_2digit_30_win_jackpot_db float64
			conf_2digit_30_operator_db                                                                                                                   string
		)

		err = row.Scan(&conf_2digit_30_time_db, &conf_2digit_30_win_db,
			&conf_2digit_30_win_redblack_db, &conf_2digit_30_win_line_db, &conf_2digit_30_win_zona_db, &conf_2digit_30_win_jackpot_db,
			&conf_2digit_30_operator_db)

		helpers.ErrorCheck(err)

		obj.Agenconf_2digit_30_time = conf_2digit_30_time_db
		obj.Agenconf_2digit_30_winangka = conf_2digit_30_win_db
		obj.Agenconf_2digit_30_winredblack = conf_2digit_30_win_redblack_db
		obj.Agenconf_2digit_30_winline = conf_2digit_30_win_line_db
		obj.Agenconf_2digit_30_winzona = conf_2digit_30_win_zona_db
		obj.Agenconf_2digit_30_winjackpot = conf_2digit_30_win_jackpot_db
		obj.Agenconf_2digit_30_operator = conf_2digit_30_operator_db
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
func Save_updateresult2D30S(admin, idrecord, idcompany, result string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

	result_db, status_db, _ := _GetInvoiceInfo(idrecord, idcompany)
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

			data_send_savetransaksi := ""
			fmt.Println("Send Data ke engine save transaksi")
			data_send_savetransaksi = idrecord + "|" + result + "|" + idcompany
			fmt.Printf("%s\n", data_send_savetransaksi)
			_senddata_enginesave(data_send_savetransaksi, idcompany)

			idcurr := _GetCompanyInfo(idcompany)
			invoice := _Generate_incoive(strings.ToLower(idcompany), idcurr)

			fieldconfig_redis := strings.ToLower(idcompany) + ":12D30S:TIMER"
			type Configure struct {
				Time    int    `json:"time"`
				Invoice string `json:"invoice"`
				Result  string `json:"result"`
			}
			var obj Configure
			obj.Time = _GetCompanyConfInfo(idcompany)
			obj.Invoice = invoice
			obj.Result = result
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
func Save_Agenconf(admin, idcompany, operator_2D30 string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_company_config + `  
				SET conf_2digit_30_operator=$1, 
				updateconf=$2, updatedateconf=$3      
				WHERE idcompany=$4   
			`

	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_company_config, "UPDATE",
		operator_2D30,
		admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany)

	if flag_update {
		msg = "Succes"
	} else {
		fmt.Println(msg_update)
	}

	//==DELETE REDIS TIMER
	val_timer := helpers.DeleteRedis("CONFIG" + "_" + strings.ToLower(idcompany))
	fmt.Printf("Redis Delete BACKEND TIMER CONFIG : %d\n", val_timer)

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
func _GetInvoiceInfo(invoice, idcompany string) (string, string, int) {
	con := db.CreateCon()
	ctx := context.Background()

	_, tbl_trx_transaksi, _ := Get_mappingdatabase(idcompany)

	result := ""
	status := ""
	total_member := 0

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "resultwigo, status_transaksi, total_member "
	sql_select += "FROM " + tbl_trx_transaksi + " "
	sql_select += "WHERE idtransaksi='" + invoice + "' "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&result, &status, &total_member); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return result, status, total_member
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

func _GetTotalBetWinByDate_Transaksi(table, startdate, enddate string) (int, int) {
	con := db.CreateCon()
	ctx := context.Background()
	total_bet := 0
	total_win := 0
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "SUM(total_bet) AS total_bet, SUM(total_win) AS total_win  "
	sql_select += "FROM " + table + " "
	sql_select += "WHERE createdate_transaksi >='" + startdate + "' "
	sql_select += "AND createdate_transaksi <='" + enddate + "' "

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&total_bet, &total_win); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return total_bet, total_win
}

// RUMUS
func _rumuswigo2D30S(tipebet, nomorclient, nomorkeluaran string) string {
	result := "LOSE"

	result_redblack, result_gangen, result_besarkecil, result_line := _nomorresult(nomorkeluaran)
	switch tipebet {
	case "ANGKA":
		if nomorclient == nomorkeluaran {
			result = "WIN"
		}
	case "REDBLACK":

		if nomorclient == result_redblack {
			result = "WIN"
		}
		if nomorclient == result_gangen {
			result = "WIN"
		}
		if nomorclient == result_besarkecil {
			result = "WIN"
		}
	case "LINE":
		if nomorclient == result_line {
			result = "WIN"
		}
	}

	return result
}
func _nomorresult(nomoresult string) (string, string, string, string) {
	type nomor_result_data struct {
		nomor_id         string
		nomor_value      string
		nomor_flag       bool
		nomor_css        string
		nomor_gangen     string
		nomor_besarkecil string
		nomor_line       string
		nomor_redblack   string
		nomor_zona       string
	}

	var cards = []nomor_result_data{
		{nomor_id: "01", nomor_value: "01", nomor_zona: "ZONA_A", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GANJIL", nomor_besarkecil: "KECIL", nomor_line: "LINE1", nomor_redblack: "RED"},
		{nomor_id: "04", nomor_value: "04", nomor_zona: "ZONA_A", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GENAP", nomor_besarkecil: "KECIL", nomor_line: "LINE2", nomor_redblack: "BLACK"},
		{nomor_id: "07", nomor_value: "07", nomor_zona: "ZONA_A", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GANJIL", nomor_besarkecil: "BESAR", nomor_line: "LINE3", nomor_redblack: "RED"},
		{nomor_id: "10", nomor_value: "10", nomor_zona: "ZONA_A", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GENAP", nomor_besarkecil: "BESAR", nomor_line: "LINE4", nomor_redblack: "BLACK"},
		{nomor_id: "02", nomor_value: "02", nomor_zona: "ZONA_B", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GENAP", nomor_besarkecil: "KECIL", nomor_line: "LINE1", nomor_redblack: "RED"},
		{nomor_id: "05", nomor_value: "05", nomor_zona: "ZONA_B", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GANJIL", nomor_besarkecil: "KECIL", nomor_line: "LINE2", nomor_redblack: "BLACK"},
		{nomor_id: "08", nomor_value: "08", nomor_zona: "ZONA_B", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GENAP", nomor_besarkecil: "BESAR", nomor_line: "LINE3", nomor_redblack: "RED"},
		{nomor_id: "11", nomor_value: "11", nomor_zona: "ZONA_B", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GANJIL", nomor_besarkecil: "BESAR", nomor_line: "LINE4", nomor_redblack: "BLACK"},
		{nomor_id: "03", nomor_value: "03", nomor_zona: "ZONA_C", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GANJIL", nomor_besarkecil: "KECIL", nomor_line: "LINE2", nomor_redblack: "RED"},
		{nomor_id: "06", nomor_value: "06", nomor_zona: "ZONA_C", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GENAP", nomor_besarkecil: "KECIL", nomor_line: "LINE2", nomor_redblack: "BLACK"},
		{nomor_id: "09", nomor_value: "09", nomor_zona: "ZONA_C", nomor_flag: false, nomor_css: "btn btn-error", nomor_gangen: "GANJIL", nomor_besarkecil: "BESAR", nomor_line: "LINE3", nomor_redblack: "RED"},
		{nomor_id: "12", nomor_value: "12", nomor_zona: "ZONA_C", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GENAP", nomor_besarkecil: "BESAR", nomor_line: "LINE4", nomor_redblack: "BLACK"},
		{nomor_id: "JP", nomor_value: "JP", nomor_flag: false, nomor_css: "btn", nomor_gangen: "GANJIL", nomor_besarkecil: "BESAR", nomor_line: "LINE3", nomor_redblack: "RED"}}

	result_redblack := ""
	result_gangen := ""
	result_besarkecil := ""
	result_line := ""
	for i := 0; i < len(cards); i++ {
		if cards[i].nomor_id == nomoresult {
			result_redblack = cards[i].nomor_redblack
			result_gangen = cards[i].nomor_gangen
			result_besarkecil = cards[i].nomor_besarkecil
			result_line = cards[i].nomor_line
		}
	}
	return result_redblack, result_gangen, result_besarkecil, result_line
}

func _senddata_enginesave(data, company string) {
	key := "payload" + "_enginesave_" + strings.ToLower(company)
	helpers.SetPublish(key, data)
}
