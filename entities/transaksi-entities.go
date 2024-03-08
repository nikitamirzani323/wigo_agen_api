package entities

type Model_transaksi2D30S struct {
	Transaksi2D30S_id         string `json:"transaksi2D30s_id"`
	Transaksi2D30S_idcurr     string `json:"transaksi2D30s_idcurr"`
	Transaksi2D30S_date       string `json:"transaksi2D30s_date"`
	Transaksi2D30S_result     string `json:"transaksi2D30s_result"`
	Transaksi2D30S_totalbet   int    `json:"transaksi2D30s_totalbet"`
	Transaksi2D30S_totalwin   int    `json:"transaksi2D30s_totalwin"`
	Transaksi2D30S_winlose    int    `json:"transaksi2D30s_winlose"`
	Transaksi2D30S_status     string `json:"transaksi2D30s_status"`
	Transaksi2D30S_status_css string `json:"transaksi2D30s_status_css"`
	Transaksi2D30S_create     string `json:"transaksi2D30s_create"`
	Transaksi2D30S_update     string `json:"transaksi2D30s_update"`
}
type Model_transaksi2D30SPrediksi struct {
	Transaksi2D30Sprediksi_id         string `json:"transaksi2D30sprediksi_id"`
	Transaksi2D30Sprediksi_username   string `json:"transaksi2D30sprediksi_username"`
	Transaksi2D30Sprediksi_date       string `json:"transaksi2D30sprediksi_date"`
	Transaksi2D30Sprediksi_nomor      string `json:"transaksi2D30sprediksi_nomor"`
	Transaksi2D30Sprediksi_bet        int    `json:"transaksi2D30sprediksi_bet"`
	Transaksi2D30Sprediksi_win        int    `json:"transaksi2D30sprediksi_win"`
	Transaksi2D30Sprediksi_winlose    int    `json:"transaksi2D30sprediksi_winlose"`
	Transaksi2D30Sprediksi_status     string `json:"transaksi2D30sprediksi_status"`
	Transaksi2D30Sprediksi_status_css string `json:"transaksi2D30sprediksi_status_css"`
}
type Controller_transaksi2D30Ssave struct {
	Page                   string `json:"page" validate:"required"`
	Transaksi2D30S_invoice string `json:"transaksi2D30s_invoice" validate:"required"`
	Transaksi2D30S_result  string `json:"transaksi2D30s_result" validate:"required"`
}
type Controller_transaksi2D30S struct {
	Transaksi2D30S_invoice string `json:"transaksi2D30s_invoice"`
}
type Controller_prediksitransaksi2D30S struct {
	Transaksi2D30Sprediksi_invoice string `json:"transaksi2D30sprediksi_invoice"`
	Transaksi2D30Sprediksi_result  string `json:"transaksi2D30sprediksi_result"`
}
