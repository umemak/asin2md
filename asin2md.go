package asin2md

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Get(asin string) (string, error) {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("https://www.amazon.co.jp/dp/%s/", asin))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	// #rpi-attribute-book_details-publisher > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span
	// //*[@id="rpi-attribute-book_details-publisher"]/div[3]/span
	doc.Find("#rpi-attribute-book_details-publisher > div.a-section.a-spacing-none.a-text-center.rpi-attribute-value > span").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
	return "", nil
}
