package entities

type Model_transaksi2D30S struct {
	Transaksi2D30S_id          string `json:"transaksi2D30s_id"`
	Transaksi2D30S_idcurr      string `json:"transaksi2D30s_idcurr"`
	Transaksi2D30S_date        string `json:"transaksi2D30s_date"`
	Transaksi2D30S_result      string `json:"transaksi2D30s_result"`
	Transaksi2D30S_totalmember int    `json:"transaksi2D30s_totalmember"`
	Transaksi2D30S_totalbet    int    `json:"transaksi2D30s_totalbet"`
	Transaksi2D30S_totalwin    int    `json:"transaksi2D30s_totalwin"`
	Transaksi2D30S_winlose     int    `json:"transaksi2D30s_winlose"`
	Transaksi2D30S_status      string `json:"transaksi2D30s_status"`
	Transaksi2D30S_status_css  string `json:"transaksi2D30s_status_css"`
	Transaksi2D30S_create      string `json:"transaksi2D30s_create"`
	Transaksi2D30S_update      string `json:"transaksi2D30s_update"`
}
type Model_transaksi2D30S_daily struct {
	Transaksi2D30Ssummarydaily_id       string `json:"transaksi2D30ssummarydaily_id"`
	Transaksi2D30Ssummarydaily_periode  string `json:"transaksi2D30ssummarydaily_periode"`
	Transaksi2D30Ssummarydaily_totalbet int    `json:"transaksi2D30ssummarydaily_totalbet"`
	Transaksi2D30Ssummarydaily_totalwin int    `json:"transaksi2D30ssummarydaily_totalwin"`
	Transaksi2D30Ssummarydaily_winlose  int    `json:"transaksi2D30ssummarydaily_winlose"`
	Transaksi2D30Ssummarydaily_create   string `json:"transaksi2D30ssummarydaily_create"`
	Transaksi2D30Ssummarydaily_update   string `json:"transaksi2D30ssummarydaily_update"`
}
type Model_transaksi2D30SInfoInvoice struct {
	Transaksi2D30Sinfo_id          string `json:"transaksi2D30sinfo_id"`
	Transaksi2D30Sinfo_result      string `json:"transaksi2D30sinfo_result"`
	Transaksi2D30Sinfo_totalmember int    `json:"transaksi2D30sinfo_totalmember"`
	Transaksi2D30Sinfo_totalbet    int    `json:"transaksi2D30sinfo_totalbet"`
	Transaksi2D30Sinfo_totalwin    int    `json:"transaksi2D30sinfo_totalwin"`
	Transaksi2D30Sinfo_winlose     int    `json:"transaksi2D30sinfo_winlose"`
	Transaksi2D30Sinfo_status      string `json:"transaksi2D30sinfo_status"`
}
type Model_transaksi2D30Ssummary struct {
	Transaksi2D30Ssummary_nomor        string `json:"transaksi2D30ssummary_nomor"`
	Transaksi2D30Ssummary_totalinvoice int    `json:"transaksi2D30ssummary_totalinvoice"`
	Transaksi2D30Ssummary_totalbet     int    `json:"transaksi2D30ssummary_totalbet"`
	Transaksi2D30Ssummary_totalwin     int    `json:"transaksi2D30ssummary_totalwin"`
}
type Model_transaksi2D30Sdetail struct {
	Transaksi2D30Sdetail_id         string  `json:"transaksi2D30sdetail_id"`
	Transaksi2D30Sdetail_date       string  `json:"transaksi2D30sdetail_date"`
	Transaksi2D30Sdetail_ipaddress  string  `json:"transaksi2D30sdetail_ipaddress"`
	Transaksi2D30Sdetail_device     string  `json:"transaksi2D30sdetail_device"`
	Transaksi2D30Sdetail_browser    string  `json:"transaksi2D30sdetail_browser"`
	Transaksi2D30Sdetail_username   string  `json:"transaksi2D30sdetail_username"`
	Transaksi2D30Sdetail_tipebet    string  `json:"transaksi2D30sdetail_tipebet"`
	Transaksi2D30Sdetail_nomor      string  `json:"transaksi2D30sdetail_nomor"`
	Transaksi2D30Sdetail_bet        int     `json:"transaksi2D30sdetail_bet"`
	Transaksi2D30Sdetail_win        int     `json:"transaksi2D30sdetail_win"`
	Transaksi2D30Sdetail_multiplier float64 `json:"transaksi2D30sdetail_multiplier"`
	Transaksi2D30Sdetail_status     string  `json:"transaksi2D30sdetail_status"`
	Transaksi2D30Sdetail_status_css string  `json:"transaksi2D30sdetail_status_css"`
	Transaksi2D30Sdetail_create     string  `json:"transaksi2D30sdetail_create"`
	Transaksi2D30Sdetail_update     string  `json:"transaksi2D30sdetail_update"`
}
type Model_transaksi2D30SPrediksi struct {
	Transaksi2D30Sprediksi_id         string  `json:"transaksi2D30sprediksi_id"`
	Transaksi2D30Sprediksi_username   string  `json:"transaksi2D30sprediksi_username"`
	Transaksi2D30Sprediksi_date       string  `json:"transaksi2D30sprediksi_date"`
	Transaksi2D30Sprediksi_nomor      string  `json:"transaksi2D30sprediksi_nomor"`
	Transaksi2D30Sprediksi_bet        int     `json:"transaksi2D30sprediksi_bet"`
	Transaksi2D30Sprediksi_multiplier float64 `json:"transaksi2D30sprediksi_multiplier"`
	Transaksi2D30Sprediksi_win        int     `json:"transaksi2D30sprediksi_win"`
	Transaksi2D30Sprediksi_winlose    int     `json:"transaksi2D30sprediksi_winlose"`
	Transaksi2D30Sprediksi_status     string  `json:"transaksi2D30sprediksi_status"`
	Transaksi2D30Sprediksi_status_css string  `json:"transaksi2D30sprediksi_status_css"`
}
type Model_agenconf struct {
	Agenconf_2digit_30_time        int     `json:"agenconf_2digit_30_time"`
	Agenconf_2digit_30_winangka    float64 `json:"agenconf_2digit_30_winangka"`
	Agenconf_2digit_30_winredblack float64 `json:"agenconf_2digit_30_winredblack"`
	Agenconf_2digit_30_winline     float64 `json:"agenconf_2digit_30_winline"`
	Agenconf_2digit_30_winzona     float64 `json:"agenconf_2digit_30_winzona"`
	Agenconf_2digit_30_winjackpot  float64 `json:"agenconf_2digit_30_winjackpot"`
	Agenconf_2digit_30_operator    string  `json:"agenconf_2digit_30_operator"`
}
type Controller_transaksi2D30Ssave struct {
	Page                   string `json:"page" validate:"required"`
	Transaksi2D30S_invoice string `json:"transaksi2D30s_invoice" validate:"required"`
	Transaksi2D30S_result  string `json:"transaksi2D30s_result" validate:"required"`
}
type Controller_agenconfsave struct {
	Page                        string `json:"page" validate:"required"`
	Agenconf_2digit_30_operator string `json:"agenconf_2digit_30_operator" validate:"required"`
}
type Controller_transaksi2D30S struct {
	Transaksi2D30S_search  string `json:"transaksi2D30s_search"`
	Transaksi2D30S_page    int    `json:"transaksi2D30s_page"`
	Transaksi2D30S_invoice string `json:"transaksi2D30s_invoice"`
}
type Controller_transaksi2D30SSummaryDaily struct {
	Transaksi2D30Ssummarydaily_search string `json:"transaksi2D30ssummarydaily_search"`
	Transaksi2D30Ssummarydaily_page   int    `json:"transaksi2D30ssummarydaily_page"`
}
type Controller_transaksi2D30Sinfo struct {
	Transaksi2D30S_invoice string `json:"transaksi2D30s_invoice"`
}
type Controller_prediksitransaksi2D30S struct {
	Transaksi2D30Sprediksi_invoice string `json:"transaksi2D30sprediksi_invoice"`
	Transaksi2D30Sprediksi_result  string `json:"transaksi2D30sprediksi_result"`
}
type Controller_transaksidetail2D30S struct {
	Transaksidetail2D30S_invoice string `json:"transaksidetail2D30s_invoice"`
	Transaksidetail2D30S_status  string `json:"transaksidetail2D30s_status"`
}
