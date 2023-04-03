package bean

type TlyLoginResp struct {
	Ret        string
	Time       int64
	Info       string
	Transfer   float32
	NodeNumber int16 `json:"node_number"`
	Node       []VmessConfig
}
