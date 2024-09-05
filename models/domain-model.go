package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/wigo_agen_api/configs"
	"bitbucket.org/isbtotogroup/wigo_agen_api/db"
	"bitbucket.org/isbtotogroup/wigo_agen_api/entities"
	"bitbucket.org/isbtotogroup/wigo_agen_api/helpers"
)

const database_domain_local = configs.DB_tbl_mst_domain

func Fetch_checkdomain(domain, tipe string) (bool, error) {
	ctx := context.Background()
	con := db.CreateCon()
	flag := true
	var result entities.Model_checkdomain
	var nmdomain, tipedomain, statusdomain string
	field_redis := "DOMAINLIST:" + strings.Replace(domain, ":", "_", -1) + "-" + tipe

	_, flagRedis := helpers.GetRedis(field_redis)

	if !flagRedis {
		sql_select := `
			SELECT
			nmdomain, tipedomain, statusdomain  
			FROM ` + database_domain_local + ` 
			WHERE nmdomain = $1 
			AND tipedomain = $2 
			AND statusdomain = 'Y' 
		`

		row := con.QueryRowContext(ctx, sql_select, domain, tipe)
		switch e := row.Scan(&nmdomain, &tipedomain, &statusdomain); e {
		case sql.ErrNoRows:
			return false, errors.New("domain is not registered")
		case nil:
			flag = true
			result.Domain_name = nmdomain
			result.Domain_tipe = tipedomain
			result.Domain_status = statusdomain
			helpers.SetRedis(field_redis, result, 5*time.Hour)
		default:
			return false, errors.New("domain is not registered")
		}
	}

	return flag, nil
}
