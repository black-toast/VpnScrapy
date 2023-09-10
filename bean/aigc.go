package bean

type EsheepResp struct {
	Code int        `json:"code"`
	Data EsheepData `json:"data"`
}

type EsheepData struct {
	Images []EsheepImage `json:"images"`
}

type EsheepImage struct {
	Url string `json:"url"`
}
