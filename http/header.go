package http

import "net/http"

func SetCommonHeader(header http.Header) {
	header.Set("Accept", WithConstant().headerAccept)
	header.Set("Accept-Language", WithConstant().headerAcceptLanguage)
	header.Set("Cache-Control", WithConstant().headerCacheControl)
	header.Set("Sec-Fetch-Dest", WithConstant().headerSecFetchDest)
	header.Set("Sec-Fetch-Mode", WithConstant().headerSecFetchMode)
	header.Set("Sec-Fetch-Site", WithConstant().headerUaChrome)
	header.Set("Sec-Fetch-User", WithConstant().headerUaChrome)
	header.Set("Upgrade-Insecure-Requests", WithConstant().headerUaChrome)
	header.Set("User-Agent", WithConstant().headerUaChrome)
	header.Set("sec-ch-ua", WithConstant().headerUaChrome)
	header.Set("sec-ch-ua-mobile", WithConstant().headerUaChrome)
	header.Set("sec-ch-ua-platform", WithConstant().headerUaChrome)
}
