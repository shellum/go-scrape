package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-sql-driver/mysql"
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

func persistWeather(channel string, low int, high int, weather string, daysOut int) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDR"),
		DBName: os.Getenv("DB_NAME"),
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO weather (low, high, weather, channel, time, days_out) VALUES( ?, ?, ?, ?, now(), ?)")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	stmtIns.Exec(low, high, weather, channel, daysOut)
	defer stmtIns.Close()
}

func main() {
	scrapeWeather()
	// Test persistence
	persistWeather("weather.com", 50, 70, "sunny", 0)
}
