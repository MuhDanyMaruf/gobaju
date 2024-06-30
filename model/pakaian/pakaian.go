package pakaian

type Pakaian struct {
	PakaianId    int    `json:"id_pakaian"`
	NamaPakaian  string `json:"namapakaian"`
	Deskripsi    string `json:"deskripsi"`
	UkuranBaju   string `json:"ukuran"`
	StokBaju     int    `json:"stok"`
	HargaPakaian int    `json:"hargapakaian"`
}
