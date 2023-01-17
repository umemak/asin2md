package asin2md

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	asin                     string
	productTitle             string
	ebook_pages              string
	handwritten_sticky_notes string
	language                 string
	publisher                string
	publication_date         string
	file_size                string
	page_flip                string
	word_wise                string
	advanced_type_setting    string
	authors                  []string
	image                    string
	buy                      string
	from                     string
	to                       string
}

func (i Info) String() string {
	ret := "---\n"
	ret += i.Tags()
	ret += "aliases: [\"" + i.productTitle + "\"]\n"
	ret += "productTitle: " + i.productTitle + "\n"
	ret += "author: " + strings.Join(i.authors, ";") + "\n"
	ret += "imageURL: " + i.image + "\n"
	ret += "ebookPages: " + i.ebook_pages + "\n"
	ret += "handwrittenStickyNotes: " + i.handwritten_sticky_notes + "\n"
	ret += "language: " + i.language + "\n"
	ret += "publisher: " + i.publisher + "\n"
	ret += "publicationDate: " + i.publication_date + "\n"
	ret += "fileSize: " + i.file_size + "\n"
	ret += "pageFlip: " + i.page_flip + "\n"
	ret += "wordWise: " + i.word_wise + "\n"
	ret += "advancedTypeSetting: " + i.advanced_type_setting + "\n"
	ret += "\n---\n"
	ret += "# " + i.productTitle + "\n"
	ret += "## Metadata\n"
	ret += i.Authors()
	ret += "* ASIN: " + i.asin + "\n"
	ret += "* Reference: https://www.amazon.co.jp/dp/" + i.asin + "\n"
	ret += "* [Kindle link](kindle://book?action=open&asin=" + i.asin + ")\n"
	ret += "* ![](" + i.image + ")\n"
	ret += "* 出版社: [[" + i.publisher + "]]\n"
	ret += "* 出版日: " + i.publication_date + "\n"
	ret += "* 購入日:: " + i.buy + "\n"
	ret += "* 開始日:: " + i.from + "\n"
	ret += "* 完了日:: " + i.to + "\n"
	// ret += "\n#unlimited\n"
	return ret
}

func (i Info) Tags() string {
	tags := []string{}
	for _, v := range i.authors {
		vv := strings.ReplaceAll(v, " ", "")
		tags = append(tags, fmt.Sprintf(`"書籍/著者/%s"`, strings.ReplaceAll(vv, "　", "")))
	}
	if i.publisher != "" {
		tags = append(tags, fmt.Sprintf(`"書籍/出版社/%s"`, strings.ReplaceAll(i.publisher, " ", "")))
	}
	tags = append(tags, fmt.Sprintf(`"書籍/発売日/%s"`, strings.ReplaceAll(i.publication_date, "-", "/")))
	return "tags: [" + strings.Join(tags, ", ") + "]\n"
}

func (i Info) Authors() string {
	authors := []string{}
	for _, v := range i.authors {
		authors = append(authors, fmt.Sprintf(`[[%s]]`, v))
	}
	return "* Author: " + strings.Join(authors, "、") + "\n"
}

func Get(asin, buy, from, to string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://www.amazon.co.jp/dp/%s/", asin))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("%s: status code error: %d %s", asin, res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	info := Info{asin: asin, buy: buy, from: from, to: to}
	// タイトル
	info.productTitle = findText(doc, "#productTitle")
	// 画像
	info.image = findSrc(doc, "#ebooksImgBlkFront")
	// 本の長さ
	info.ebook_pages = findText(doc, "#rpi-attribute-book_details-ebook_pages > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span > a > span")
	// 付箋メモ
	info.handwritten_sticky_notes = findText(doc, "#rpi-attribute-book_details-handwritten_sticky_notes > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span > a > span")
	// 言語
	info.language = findText(doc, "#rpi-attribute-language > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span")
	// 出版社
	info.publisher = findText(doc, "#rpi-attribute-book_details-publisher > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span")
	// 発売日
	info.publication_date = dateConvert(findText(doc, "#rpi-attribute-book_details-publication_date > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span"))
	// ファイルサイズ
	info.file_size = findText(doc, "#rpi-attribute-book_details-file_size > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span")
	// Page Flip
	info.page_flip = findText(doc, "#rpi-attribute-book_details-page_flip > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span > a > span")
	// Word Wise
	info.word_wise = findText(doc, "#rpi-attribute-book_details-word_wise > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span > a > span")
	// タイプセッティングの改善
	info.advanced_type_setting = findText(doc, "#rpi-attribute-book_details-advanced_type_setting > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span > a > span")
	// 著者
	info.authors = findAuthors(doc)
	return info.String(), nil
}

func findText(doc *goquery.Document, selector string) string {
	ret := ""
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		ret = strings.TrimSpace(s.Text())
	})
	return ret
}

func findSrc(doc *goquery.Document, selector string) string {
	ret := ""
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			ret = strings.TrimSpace(src)
		}
	})
	return ret
}

func findAuthors(doc *goquery.Document) []string {
	res := []string{}
	doc.Find("#bylineInfo > span.author").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a.contributorNameID").Text()
		if name == "" {
			name = s.Find("a.a-link-normal").Text()
		}
		res = append(res, strings.TrimSpace(name))
	})
	return res
}

func dateConvert(src string) string {
	items := strings.Split(src, "/")
	y, err := strconv.ParseInt(items[0], 10, 64)
	if err != nil {
		return src
	}
	m, err := strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		return src
	}
	d, err := strconv.ParseInt(items[2], 10, 64)
	if err != nil {
		return src
	}
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}
