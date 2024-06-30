package pembeli

type Pembeli struct {
	PembeliId       int    `json:"id_pembeli"`
	Nama            string `json:"nama"`
	Email           string `json:"email"`
	JumlahPembelian int    `json:"jumlahpembelian"`
}
