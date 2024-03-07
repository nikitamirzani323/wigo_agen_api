package models

import (
	"context"
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

const database_adminrule_local = configs.DB_tbl_mst_company_adminrule

func Fetch_adminruleHome(idcompany string) (helpers.Response, error) {
	var obj entities.Model_adminruleall
	var arraobj []entities.Model_adminruleall
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcompadminrule, nmruleadmin, ruleadmin 
			FROM ` + database_adminrule_local + ` 
			ORDER BY nmruleadmin ASC  
		`

	row, err := con.QueryContext(ctx, sql_select)

	var no int = 0
	helpers.ErrorCheck(err)
	for row.Next() {
		no += 1
		var (
			idcompadminrule_db           int
			nmruleadmin_db, ruleadmin_db string
		)

		err = row.Scan(&idcompadminrule_db, &nmruleadmin_db, &ruleadmin_db)

		helpers.ErrorCheck(err)

		obj.Adminrule_id = idcompadminrule_db
		obj.Adminrule_name = nmruleadmin_db
		obj.Adminrule_rule = ruleadmin_db
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
func Save_adminrule(admin, idcompany, name, rule, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + database_adminrule_local + ` (
					idcompadminrule,idcompany,nmruleadmin, 
					createcompadminrule,createdatecompadminrule 
				) values (
					$1,$2,$3
					$4,$5
				) 
			`
		field_column := database_companyadminrule_local + strings.ToLower(idcompany) + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		idrecord := tglnow.Format("YY") + strconv.Itoa(idrecord_counter)
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_adminrule_local, "INSERT",
			idrecord, idcompany, name,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_adminrule_local + `   
				SET nmruleadmin=$1, ruleadmin=$2, 
				updatecompadminrule=$3, updatedatecompadminrule=$3
				WHERE idadmin=$5 
			`
		flag_update, msg_update := Exec_SQL(sql_update, database_adminrule_local, "UPDATE",
			name, rule,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
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
