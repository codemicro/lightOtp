package models

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
