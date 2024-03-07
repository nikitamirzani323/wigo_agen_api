package models

import (
	"context"
	"database/sql"
	"fmt"
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

const database_admin_local = configs.DB_tbl_mst_company_admin
const database_rule_local = configs.DB_tbl_mst_company_adminrule

func Fetch_adminHome(idcompany string) (helpers.ResponseAdmin, error) {
	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var res helpers.ResponseAdmin
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idcompadmin, adminusername , nameadmin, idcompadminrule,
		statuscompadmin, 
		to_char(COALESCE(lastlogincompadmin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		to_char(COALESCE(createdatecompadmin,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		ipaddresscompadmin   
		FROM ` + database_admin_local + ` 
		WHERE idcompany=$1 
		ORDER BY lastlogincompadmin DESC 
	`

	row, err := con.QueryContext(ctx, sql_select, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcompadmin_db, idcompadminrule_db                                                       int
			adminusername_db, nameadmin_db                                                           string
			statuscompadmin_db, lastlogincompadmin_db, createdatecompadmin_db, ipaddresscompadmin_db string
		)

		err = row.Scan(
			&idcompadmin_db, &adminusername_db, &nameadmin_db, &idcompadminrule_db,
			&statuscompadmin_db, &lastlogincompadmin_db, &createdatecompadmin_db,
			&ipaddresscompadmin_db)

		helpers.ErrorCheck(err)
		status_css := configs.STATUS_CANCEL

		if statuscompadmin_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}
		if lastlogincompadmin_db == createdatecompadmin_db {
			lastlogincompadmin_db = ""
		}
		obj.Admin_id = idcompadmin_db
		obj.Admin_username = adminusername_db
		obj.Admin_nama = nameadmin_db
		obj.Admin_idrule = idcompadminrule_db
		obj.Admin_rule = _Get_adminrule(idcompadminrule_db, idcompany)
		obj.Admin_joindate = createdatecompadmin_db
		obj.Admin_lastlogin = lastlogincompadmin_db
		obj.Admin_lastIpaddress = ipaddresscompadmin_db
		obj.Admin_status = statuscompadmin_db
		obj.Admin_status_css = status_css
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objRule entities.Model_adminrule
	var arraobjRule []entities.Model_adminrule
	sql_listrule := `SELECT 
		idcompadminrule, nmruleadmin 	
		FROM ` + database_rule_local + ` 
		WHERE idcompany='` + idcompany + `' 
	`
	row_listrule, err_listrule := con.QueryContext(ctx, sql_listrule)

	helpers.ErrorCheck(err_listrule)
	for row_listrule.Next() {
		var (
			idcompadminrule_DB int
			nmruleadmin        string
		)

		err = row_listrule.Scan(&idcompadminrule_DB, &nmruleadmin)

		helpers.ErrorCheck(err)

		objRule.Adminrule_id = idcompadminrule_DB
		objRule.Adminrule_name = nmruleadmin
		arraobjRule = append(arraobjRule, objRule)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listrule = arraobjRule
	res.Time = time.Since(start).String()

	return res, nil
}

func Save_adminHome(admin, idcompany, username, password, nama, status, sData string, idrecord, idrule int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_admin_local, "username", username)
		if !flag {
			sql_insert := `
				insert into
				` + database_admin_local + ` (
					idcompadmin, idcompadminrule, idcompany,
					adminusername , adminpassword, nameadmin, statuscompadmin, lastlogincompadmin,
					createcompadmin, createdatecompadmin
				) values (
					$1, $2, $3, 
					$4, $5, $6, $7, $8, 
					$9, $10
				)
			`
			createdate := tglnow.Format("YYYY-MM-DD HH:mm:ss")
			hashpass := helpers.HashPasswordMD5(password)
			field_column := database_companyadmin_local + strings.ToLower(idcompany) + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			idrecord := tglnow.Format("YY") + strconv.Itoa(idrecord_counter)
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_admin_local, "INSERT",
				idrecord, idrule, idcompany,
				username, hashpass, nama, status, createdate,
				admin, createdate)

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + database_admin_local + `  
				SET nameadmin =$1, idcompadminrule=$2, statuscompadmin=$3,  
				updatecompadmin=$4, updatedatecompadmin=$5 
				WHERE idcompadmin=$6 
			`

			flag_update, msg_update := Exec_SQL(sql_update, database_admin_local, "UPDATE",
				nama, idrule, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_admin + `   
				SET nameadmin=$1, adminpassword=$2, idcompadminrule=$3, statuscompadmin=$4,  
				updatecompadmin=$5, updatedatecompadmin=$6 
				WHERE idcompadmin =$7 
			`
			flag_update, msg_update := Exec_SQL(sql_update2, database_admin_local, "UPDATE",
				nama, hashpass, idrule, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username)

			if flag_update {
				msg = "Succes"
			} else {
				fmt.Println(msg_update)
			}
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _Get_adminrule(idrule int, idcompany string) string {
	con := db.CreateCon()
	ctx := context.Background()
	nmruleadmin := ""
	sql_select := `SELECT
			nmruleadmin    
			FROM ` + database_companyadminrule_local + `  
			WHERE idcompadminrule='` + strconv.Itoa(idrule) + `'       
			AND idcompany='` + idcompany + `'       
		`

	row := con.QueryRowContext(ctx, sql_select)
	switch e := row.Scan(&nmruleadmin); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}

	return nmruleadmin
}
