package entities

type Model_invoiceclient struct {
	Invoiceclient_id         string  `json:"invoiceclient_id"`
	Invoiceclient_date       string  `json:"invoiceclient_date"`
	Invoiceclient_result     string  `json:"invoiceclient_result"`
	Invoiceclient_username   string  `json:"invoiceclient_username"`
	Invoiceclient_nomor      string  `json:"invoiceclient_nomor"`
	Invoiceclient_bet        int     `json:"invoiceclient_bet"`
	Invoiceclient_win        int     `json:"invoiceclient_win"`
	Invoiceclient_multiplier float64 `json:"invoiceclient_multiplier"`
	Invoiceclient_status     string  `json:"invoiceclient_status"`
	Invoiceclient_status_css string  `json:"invoiceclient_status_css"`
}

type Model_result struct {
	Result_invoice string `json:"result_invoice"`
	Result_date    string `json:"result_date"`
	Result_result  string `json:"result_result"`
}
type Model_invoicedetail struct {
	Invoicedetail_id     string `json:"invoicedetail_id"`
	Invoicedetail_date   string `json:"invoicedetail_date"`
	Invoicedetail_round  int    `json:"invoicedetail_round"`
	Invoicedetail_bet    int    `json:"invoicedetail_bet"`
	Invoicedetail_win    int    `json:"invoicedetail_win"`
	Invoicedetail_status string `json:"invoicedetail_status"`
}

type Controller_invoice struct {
	Invoice_company  string `json:"invoice_company" validate:"required"`
	Invoice_username string `json:"invoice_username" validate:"required"`
}
type Controller_result struct {
	Invoice_company string `json:"invoice_company" validate:"required"`
}
type Controller_invoicedetail struct {
	Invoice_id      string `json:"invoice_id" validate:"required"`
	Invoice_company string `json:"invoice_company" validate:"required"`
}
type Controller_transaksisave struct {
	Transaksi_company string `json:"transaksi_company" validate:"required"`
	Transaksi_idcurr  string `json:"transaksi_idcurr" validate:"required"`
}

// idtransaksi, resulcard_win string, round_bet, bet, c_before, c_after, win, idpoin int
type Controller_transaksidetailsave struct {
	Transaksidetail_company     string  `json:"transaksidetail_company" validate:"required"`
	Transaksidetail_idtransaksi string  `json:"transaksidetail_idtransaksi" validate:"required"`
	Transaksidetail_username    string  `json:"transaksidetail_username" validate:"required"`
	Transaksidetail_nomor       string  `json:"transaksidetail_nomor" validate:"required"`
	Transaksidetail_bet         int     `json:"transaksidetail_bet"`
	Transaksidetail_multiplier  float64 `json:"transaksidetail_multiplier"`
}
