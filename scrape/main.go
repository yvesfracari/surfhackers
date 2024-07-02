package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(

		colly.AllowedDomains("surfmappers-lowres-prod-vi.s3.us-east-1.amazonaws.com", "graphql.aws2.surfmappers.com", "www.surfmappers.com", "*.surfmappers.com", "surfmappers.com"),
	)

	foundItems := make(map[string]map[string]string)

	// // On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")

	// 	// Print link
	// 	// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	// 	// Visit link found on page
	// 	// Only those links are visited which are in AllowedDomains
	// 	c.Visit(e.Request.AbsoluteURL(link))
	// })

	// c.OnHTML("a > div[class='album-link']", func(e *colly.HTMLElement) {
	// 	fmt.Printf("Album link found: %q -> %s\n", e.Text, e.Attr("href"))
	// })

	c.OnHTML("a > div[class='session-cover']", func(e *colly.HTMLElement) {
		sessionUrl, _ := e.DOM.Parent().Attr("href")

		foundItems[sessionUrl] = make(map[string]string)
		foundItems[sessionUrl]["sessionUrl"] = sessionUrl

		sessionCoverUrl := e.Attr("style")
		backgroundImagePattern := regexp.MustCompile(`background-image:\s*url\(([^)]+)\)`)
		matches := backgroundImagePattern.FindStringSubmatch(sessionCoverUrl)
		if len(matches) > 1 {
			url := matches[1]
			foundItems[sessionUrl]["sessionCoverUrl"] = url
		}
	})

	c.OnHTML("div[class='gallery-item gallery-item-button']", func(e *colly.HTMLElement) {
		lowResImage := e.Attr("style")

		fmt.Printf("Low res image found: %q -> %s\n", e.Text, lowResImage)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.surfmappers.com/surfguidesc/sessions")

	jsonString, err := json.Marshal(foundItems)

	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile("data.json", jsonString, os.ModePerm)

	fmt.Println(string(jsonString))
}