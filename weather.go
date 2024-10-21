package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-sql-driver/mysql"
)

type Weather struct {
	High    string
	Low     string
	Weather string
	Channel string
	DaysOut int
}

const FORECAST_DAYS = 14

func scrapeWeather() [FORECAST_DAYS]Weather {
	res, _ := http.Get("https://weather.com/weather/tenday/l/7396101a49100eaa70080ede34244e1a9ae8e2a237da460c85e6472f1c8ed113")
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	numericReg, _ := regexp.Compile("[^0-9]")
	high, low := "", ""
	var weatherData [FORECAST_DAYS]Weather
	tempSelection := doc.Find("summary > div > div[data-testid='DetailsSummary'] div[data-testid='detailsTemperature'] span[data-testid='TemperatureValue']")
	weatherSelection := doc.Find("div[data-testid='wxIcon'] span")
	for daysOut := 0; daysOut < tempSelection.Length(); daysOut++ {
		tempHtmlNode := tempSelection.Eq(daysOut)
		temperatureText := tempHtmlNode.Text()
		temperatureNumeric := numericReg.ReplaceAllString(temperatureText, "")
		weatherHtmlNode := weatherSelection.Eq((daysOut + 1) / 2)
		weatherText := strings.ToLower(weatherHtmlNode.Text())
		// TODO: convert weatherText to one of: Rain, Clouds, Snow, Sun
		if strings.Contains(weatherText, "rain") {
			weatherText = "Rain"
		} else if strings.Contains(weatherText, "cloud") {
			weatherText = "Clouds"
		} else if strings.Contains(weatherText, "snow") {
			weatherText = "Snow"
		} else if strings.Contains(weatherText, "sun") {
			weatherText = "Sun"
		}
		if high == "" {
			high = temperatureNumeric
		} else {
			low = temperatureNumeric
			weatherData[daysOut/2] = Weather{High: high, Low: low, Weather: weatherText, Channel: "weather.com", DaysOut: daysOut/2 + 1}
			high, low = "", ""
		}
	}
	return weatherData
}

// DB persistence for now
func persistWeather(weatherData [FORECAST_DAYS]Weather) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDR"),
		DBName: os.Getenv("DB_NAME"),
	}
	fmt.Printf("User: %s, Pass: %s, Addr: %s, DBName: %s", cfg.User, cfg.Passwd, cfg.Addr, cfg.DBName)
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO weather (low, high, weather, channel, time, days_out) VALUES( ?, ?, ?, ?, now(), ?)")
	defer stmtIns.Close()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	for _, weather := range weatherData {
		stmtIns.Exec(weather.Low, weather.High, weather.Weather, weather.Channel, weather.DaysOut)
	}
}

func gatherAndPersist() {
	fmt.Printf("Scraping from weather.com")
	weatherData := scrapeWeather()
	persistWeather(weatherData)
}

func main() {
	fmt.Printf("Starting scraper...")
	gatherAndPersist()

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		gatherAndPersist()
	}
}
