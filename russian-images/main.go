package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Image struct {
	Link string
	Name string
}

func main() {
	log.Println("Collector initialization")

	collector := colly.NewCollector()

	log.Println("Song file creation")
	file, err := os.Create("./russian-images.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	log.Println("CSV writer initialize")
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header in cvs file
	err = writer.Write([]string{"Name", "Link"})
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	images := make([]Image, 0, 5)

	// log on response
	collector.OnResponse(func(r *colly.Response) {
		log.Println("Response received", r.StatusCode)
	})

	// a := regexp.MustCompile(`"`)
	// scraping
	collector.OnXML(`//*[@id="asset-nesterovich1-95413"]/div/div[2]/div/div[2]/text()`, func(x *colly.XMLElement) {
		image := Image{
			Link: strings.ReplaceAll(x.Text, `"`, ""),
		}
		images = append(images, image)
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		log.Println("Scrapping", r.URL.String())
	})

	collector.Visit("https://nesterovich1.livejournal.com/95413.html")
	collector.Wait()

	collector2 := colly.NewCollector()
	cpt := 0
	collector2.OnHTML("div.asset-body > div[class] > img", func(e *colly.HTMLElement) {
		if cpt < len(images)-1 {
			images[cpt].Name = e.Attr("src")
			err := writer.Write([]string{images[cpt].Name, images[cpt].Link})
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}
			cpt++
		}
	})
	collector2.Visit("https://nesterovich1.livejournal.com/95413.html")

	// can do something with images array
}
