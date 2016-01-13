package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocumentFromReader(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	html, err := doc.Html()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(html)
}
