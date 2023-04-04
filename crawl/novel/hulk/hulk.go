package hulk

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

type HulkScrapy struct {
}

func (scrapy HulkScrapy) CreateParseDoc(content string) *html.Node {
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		fmt.Println("parse doc fail")
		return nil
	}
	return doc
}

func (scrapy HulkScrapy) ParseTitle(doc *html.Node) string {
	node := htmlquery.FindOne(doc, xpath.NovelTitle)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseAuthor(doc *html.Node) string {
	node := htmlquery.FindOne(doc, xpath.NovelAuthor)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseDesc(doc *html.Node) string {
	node := htmlquery.FindOne(doc, xpath.NovelDesc)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}
