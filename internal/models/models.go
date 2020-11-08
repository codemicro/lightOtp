package models

type CodesFile struct {
	Checksum string `json:"checksum"`
	Codes    string `json:"codes"`
}

type TOTPCode struct {
	Issuer      string `json:"issuer"`
	AccountName string `json:"accountName"`
	Digits      int    `json:"digits"`
	Secret      string `json:"secret"`
	Period      uint   `json:"period"`
}

type Settings struct {
	CodesLocation     string `json:"codesLocation"`
	DefaultCodeLength int    `json:"defaultCodeLength"`
}
