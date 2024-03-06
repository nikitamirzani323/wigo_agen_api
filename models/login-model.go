package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikitamirzani323/WIGO_AGEN_API/configs"
	"github.com/nikitamirzani323/WIGO_AGEN_API/db"
	"github.com/nikitamirzani323/WIGO_AGEN_API/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(username, password, ipaddress string) (bool, string, string, error) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	tglnow, _ := goment.New()
	var idcompanyDB, passwordDB, idadminDB string
	sql_select := `
			SELECT
			adminpassword, idcompany,idcompadminrule    
			FROM ` + configs.DB_tbl_mst_company_admin + ` 
			WHERE adminusername  = $1
			AND statuscompadmin = 'Y' 
		`

	fmt.Println(sql_select, username)
	row := con.QueryRowContext(ctx, sql_select, username)
	switch e := row.Scan(&passwordDB, &idcompanyDB, &idadminDB); e {
	case sql.ErrNoRows:
		return false, "", "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)

	if hashpass != passwordDB {
		return false, "", "", nil
	}

	if flag {
		sql_update := `
			UPDATE ` + configs.DB_tbl_mst_company_admin + ` 
			SET lastlogincompadmin=$1, ipaddresscompadmin=$2,  
			updatecompadmin=$3,  updatedatecompadmin=$4   
			WHERE adminusername  = $5 
			AND statuscompadmin = 'Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_company_admin, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), ipaddress,
			username, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username)

		if flag_update {
			flag = true
			fmt.Println(msg_update)
		} else {
			fmt.Println(msg_update)
		}
	}

	return true, idcompanyDB, idadminDB, nil
}
