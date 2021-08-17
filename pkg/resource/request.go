package resource

type Request struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	QName string `json:"qName"`
	QType QType  `json:"qType"`
}
