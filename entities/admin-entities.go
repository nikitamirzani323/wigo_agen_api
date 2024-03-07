package entities

type Model_admin struct {
	Admin_id            int    `json:"admin_id"`
	Admin_username      string `json:"admin_username"`
	Admin_nama          string `json:"admin_nama"`
	Admin_idrule        int    `json:"admin_idrule"`
	Admin_rule          string `json:"admin_rule"`
	Admin_joindate      string `json:"admin_joindate"`
	Admin_lastlogin     string `json:"admin_lastlogin"`
	Admin_lastIpaddress string `json:"admin_lastipaddres"`
	Admin_status        string `json:"admin_status"`
	Admin_status_css    string `json:"admin_status_css"`
}
type Model_adminrule struct {
	Adminrule_id   int    `json:"adminrule_id"`
	Adminrule_name string `json:"adminrule_name"`
}
type Model_adminsave struct {
	Username string `json:"admin_username"`
	Nama     string `json:"admin_nama"`
	Rule     string `json:"admin_rule"`
	Status   string `json:"admin_status"`
	Create   string `json:"admin_create"`
	Update   string `json:"admin_update"`
}
type Controller_admindetail struct {
	Username string `json:"admin_username" validate:"required"`
}
type Controller_adminsave struct {
	Sdata          string `json:"sdata" validate:"required"`
	Page           string `json:"page" validate:"required"`
	Admin_id       int    `json:"admin_id"`
	Admin_username string `json:"admin_username" validate:"required"`
	Admin_password string `json:"admin_password"`
	Admin_nama     string `json:"admin_nama" validate:"required"`
	Admin_idrule   int    `json:"admin_idrule" validate:"required"`
	Admin_status   string `json:"admin_status"`
}

type Responseredis_adminhome struct {
	Admin_username     string `json:"admin_username"`
	Admin_nama         string `json:"admin_nama"`
	Admin_rule         string `json:"admin_rule"`
	Admin_joindate     string `json:"admin_joindate"`
	Admin_timezone     string `json:"admin_timezone"`
	Admin_lastlogin    string `json:"admin_lastlogin"`
	Admin_lastipaddres string `json:"admin_lastipaddres"`
	Admin_status       string `json:"admin_status"`
}
type Responseredis_adminrule struct {
	Adminrule_idrule string `json:"adminrule_idruleadmin"`
}
