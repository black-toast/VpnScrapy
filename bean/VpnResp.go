package bean

type TlyLoginResp struct {
	Ret        string
	Time       int64
	Info       string
	Transfer   float32
	NodeNumber int16 `json:"node_number"`
	Node       []VmessConfig
}

type SsrSubscribe struct {
	Remarks    string `json:"remarks"`     // VIP6专属 香港 QAT - Accelerate 01
	Server     string `json:"server"`      // qat01.virtual-strike.com
	ServerPort int32  `json:"server_port"` // 13870
	Method     string `json:"method"`      // rc4-md5
	Password   string `json:"password"`    // rjKqfL
}
