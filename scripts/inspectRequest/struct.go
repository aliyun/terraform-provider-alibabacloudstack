package main

type Parm struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	RealName string `json:"realname"`
}

type Api struct {
	FileName  string `json:"fileName"`
	Name      string `json:"name"`
	Product   string `json:"product"`
	Version   string `json:"version"`
	name      string // 此字段不被导出和编码
	Requests  []Parm `json:"requests"`
	Responses []Parm `json:"responses"`
	Type      string
}
