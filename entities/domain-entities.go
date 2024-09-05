package entities

type Model_checkdomain struct {
	Domain_name   string `json:"domain_name"`
	Domain_tipe   string `json:"domain_tipe"`
	Domain_status string `json:"domain_status"`
}
type Controller_domaincheck struct {
	Domain string `json:"domain" validate:"required"`
	Tipe   string `json:"tipe" validate:"required"`
}
