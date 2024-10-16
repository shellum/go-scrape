package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func scrapeWeather() {
	res, _ := http.Get("https://weather.com/weather/tenday/l/7396101a49100eaa70080ede34244e1a9ae8e2a237da460c85e6472f1c8ed113")
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	numericReg, _ := regexp.Compile("[^0-9]")
	high, low := "", ""
	doc.Find("summary > div > div[data-testid='DetailsSummary'] div[data-testid='detailsTemperature'] span[data-testid='TemperatureValue']").Each(func(i int, s *goquery.Selection) {
		temperatureText := s.Text()
		temperatureNumeric := numericReg.ReplaceAllString(temperatureText, "")
		if high == "" {
			high = temperatureNumeric
		} else {
			low = temperatureNumeric
			fmt.Printf("High: %s, Low: %s\n", high, low)
			high, low = "", ""
		}
	})
}

func main() {
	scrapeWeather()
}
