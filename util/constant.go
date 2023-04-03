package util

type Constant struct {
	vmessFormat string
	tlyEmail    string
	tlyPassword string
	aesKey      []byte
	aesIv       []byte
}

var constant *Constant

func init() {
	constant = new(Constant)
	constant.vmessFormat = "{\"v\": \"2\",\"ps\": \"%s\",\"add\": \"%s\",\"port\": \"%d\",\"id\": \"%s\",\"aid\": \"0\",\"scy\": \"auto\",\"net\": \"%s\",\"type\": \"%s\",\"host\": \"\",\"path\": \"%s\",\"tls\": \"%s\",\"sni\": \"\",\"alpn\": \"\"}"
	constant.tlyEmail = "1479405751@qq.com"
	constant.tlyPassword = "1479405751yy"
	constant.aesKey = []byte("tlynet923456789k")
	constant.aesIv = []byte("9987654321fedcsu")
}

func WithConstant() *Constant {
	return constant
}

func (constant *Constant) GetVmessFormat() string {
	return constant.vmessFormat
}

func (constant *Constant) GetTlyEmail() string {
	return constant.tlyEmail
}

func (constant *Constant) GetTlyPassword() string {
	return constant.tlyPassword
}

func (constant *Constant) GetAesKey() []byte {
	return constant.aesKey
}

func (constant *Constant) GetAesIv() []byte {
	return constant.aesIv
}
