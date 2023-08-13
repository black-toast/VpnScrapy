package http

type Constant struct {
	headerUaOkhttp                string
	headerUaChrome                string
	headerAccept                  string
	headerAcceptLanguage          string
	headerCacheControl            string
	headerSecFetchDest            string
	headerSecFetchMode            string
	headerSecFetchSite            string
	headerSecFetchUser            string
	headerUpgradeInsecureRequests string
	headerSecChUa                 string
	headerSecChUaMobile           string
	headerSecChUaPlatform         string
}

var constant *Constant

func init() {
	constant = new(Constant)
	constant.headerUaOkhttp = "okhttp/3.2.0"
	constant.headerUaChrome = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	constant.headerAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	constant.headerAcceptLanguage = "zh-CN,zh;q=0.9,en;q=0.8"
	constant.headerCacheControl = "no-cache"
	constant.headerSecFetchDest = "document"
	constant.headerSecFetchMode = "navigate"
	constant.headerSecFetchSite = "none"
	constant.headerSecFetchUser = "?1"
	constant.headerUpgradeInsecureRequests = "1"
	constant.headerSecChUa = "\"Not_A Brand\";v=\"99\", \"Google Chrome\";v=\"109\", \"Chromium\";v=\"109\""
	constant.headerSecChUaMobile = "?0"
	constant.headerSecChUaPlatform = "\"Windows\""
}

func WithConstant() *Constant {
	return constant
}
