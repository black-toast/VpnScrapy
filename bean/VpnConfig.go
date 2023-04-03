package bean

type VmessConfig struct {
	Id          int16  // 1
	NodeServer  string `json:"node_server"` // "yd.zoommeeting.work"
	NodeName    string `json:"node_name"`   // \u81ea\u52a8|cn-auto|\u6b63\u5e38
	NodeMethod  string `json:"node_method"` // ws
	NodeStatus  string `json:"node_status"` // \u6b63\u5e38
	WsPath      string `json:"ws-path"`     // \/Tlyx3pBN1vez3NQudNkB
	Tls         string // true
	TlsHostName string `json:"tls-hostname"` // yd.zoommeeting.work
	AlterId     string // 0
	Recommend   string // 1
	Cipher      string // none
	Port        int32  // 8830
	Pass        string // b34daf14-7e12-5c28-a909-d138dc342a81
}
