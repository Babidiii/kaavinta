package main

import (
	"encoding/csv"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.Println("Collector initialization")
	collector := colly.NewCollector()

	log.Println("Song file creation")
	file, err := os.Create("./songs.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	log.Println("CSV writer initialize")
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header in cvs file
	err = writer.Write([]string{"Source", "Date", "Link"})
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	// log on response
	collector.OnResponse(func(r *colly.Response) {
		log.Println("Response received", r.StatusCode)
	})

	// scraping
	collector.OnHTML(".elli-content", func(e *colly.HTMLElement) {
		title := e.ChildText("a[href].elco-anchor")
		date := e.ChildText("span.elco-date")
		link := e.ChildText("div.elli-annotation-content > a[href]")

		// remove bracket ()
		if date != "" {
			date = date[1 : len(date)-1]
		}

		// replace comam by hyphen
		title = strings.Replace(title, ",", "-", -1)

		log.Println(title, date, link)
		err := writer.Write([]string{title, date, link})
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		log.Println("Scrapping", r.URL.String())
	})

	// link to scrap
	for p := 1; p < 16; p++ {
		collector.Visit("https://www.senscritique.com/liste/VGM_Anthologie_des_meilleures_musiques_de_jeux_video_retro_l/250682/page-" + strconv.Itoa(p))
	}
}
