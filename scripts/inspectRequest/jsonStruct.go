package main

type Parameter struct {
	Name        string `json:"@name"`
	Type        string `json:"@type"`
	TagName     string `json:"@tagName"`
	Required    string `json:"@required"`
	TagPosition string `json:"@tagPosition"`
}
type ParameterGroup struct {
	Index      string      `json:"@index"`
	Parameters []Parameter `json:"Parameter"`
}
type Parameters struct {
	Parameter []Parameter `json:"Parameter"`
}

type Parameters2 struct {
	ParameterGroups []ParameterGroup `json:"ParameterGroup"`
}

type ResultMapping1 struct {
	Member Member `json:"Member"`
}

type ResultMapping2 struct {
	Member []Member `json:"Member"`
}

type Member struct {
	Name    string `json:"@name"`
	Type    string `json:"@type"`
	TagName string `json:"@tagName"`
}

type Result1 struct {
	Uuid          string         `json:"uuid"`
	Registrant    interface{}    `json:"registrant"`
	Managers      []interface{}  `json:"managers"`
	Versionid     int            `json:"versionid"`
	DataFrom      string         `json:"data_from"`
	FormatedData  string         `json:"formated_data"`
	Name          string         `json:"@name"`
	Type          string         `json:"@type"`
	Status        string         `json:"@status"`
	Product       string         `json:"@product"`
	Version       string         `json:"@version"`
	AuthType      string         `json:"@authType"`
	Visibility    string         `json:"@visibility"`
	ResponseLog   string         `json:"@responseLog"`
	TagPosition   string         `json:"@tagPosition"`
	Parameters    Parameters     `json:"Parameters"`
	ResultMapping ResultMapping1 `json:"ResultMapping"`
	ControlPolicy interface{}    `json:"@controlPolicy"`
	IsolationType interface{}    `json:"@isolationType"`
	ParameterType string         `json:"@parameterType"`
	ProductExact  string         `json:"product"`
}

type Result2 struct {
	Uuid          string         `json:"uuid"`
	Registrant    interface{}    `json:"registrant"`
	Managers      []interface{}  `json:"managers"`
	Versionid     int            `json:"versionid"`
	DataFrom      string         `json:"data_from"`
	FormatedData  string         `json:"formated_data"`
	Name          string         `json:"@name"`
	Type          string         `json:"@type"`
	Status        string         `json:"@status"`
	Product       string         `json:"@product"`
	Version       string         `json:"@version"`
	AuthType      string         `json:"@authType"`
	Visibility    string         `json:"@visibility"`
	ResponseLog   string         `json:"@responseLog"`
	TagPosition   string         `json:"@tagPosition"`
	Parameters    Parameters     `json:"Parameters"`
	ResultMapping ResultMapping2 `json:"ResultMapping"`
	ControlPolicy interface{}    `json:"@controlPolicy"`
	IsolationType interface{}    `json:"@isolationType"`
	ParameterType string         `json:"@parameterType"`
	ProductExact  string         `json:"product"`
}

type Result3 struct {
	Uuid          string        `json:"uuid"`
	Registrant    interface{}   `json:"registrant"`
	Managers      []interface{} `json:"managers"`
	Versionid     int           `json:"versionid"`
	DataFrom      string        `json:"data_from"`
	FormatedData  string        `json:"formated_data"`
	Name          string        `json:"@name"`
	Type          string        `json:"@type"`
	Status        string        `json:"@status"`
	Product       string        `json:"@product"`
	Version       string        `json:"@version"`
	AuthType      string        `json:"@authType"`
	Visibility    string        `json:"@visibility"`
	ResponseLog   string        `json:"@responseLog"`
	TagPosition   string        `json:"@tagPosition"`
	Parameters    Parameters2   `json:"Parameters"`
	ControlPolicy interface{}   `json:"@controlPolicy"`
	IsolationType interface{}   `json:"@isolationType"`
	ParameterType string        `json:"@parameterType"`
	ProductExact  string        `json:"product"`
}

type ApiResponse1 struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Page     int         `json:"page"`
	Results  []Result1   `json:"results"`
}

type ApiResponse2 struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Page     int         `json:"page"`
	Results  []Result2   `json:"results"`
}

type ApiResponse3 struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Page     int         `json:"page"`
	Results  []Result3   `json:"results"`
}
