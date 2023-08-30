package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/urfave/cli/v2"
)

var userName string = ""
var filmsWithoutUserRatingCounter int = 0
var sumOfUserRatings float64 = 0
var amountOfFilmsWithUserRatings float64 = 0

type film struct {
	userRating float64
	pageURL string
	title string
	length int
	letterboxdRating float64
}

var films []film

func symboltoRating(symbol string) (float64, error) {
	symbolToRatingMap := map[string]float64{
		"½": 0.5,
		"★½": 1.5,
		"★★½": 2.5,
		"★★★½": 3.5,
		"★★★★½": 4.5,
		"★": 1.0,
		"★★": 2.0,
		"★★★": 3.0,
		"★★★★": 4.0,
		"★★★★★": 5.0,
	}

	symbol = strings.TrimSpace(symbol)

	if rating, ok := symbolToRatingMap[symbol]; ok {
		return rating, nil
	}

	return 0.0, fmt.Errorf("invalid symbol %s", symbol)
}

func sumOfSlice[T float64 | int](slice []film, field string) T {
	var sum T

	if field == "userRating" {
		for _, film := range slice {
			sum += T(film.userRating)
		}
	} else if field == "length" {
		for _, film := range slice {
			sum += T(film.length)
		}
	} else if field == "letterboxdRating" {
		for _, film := range slice {
			sum += T(film.letterboxdRating)
		}
	}
	
	return sum
}

func roundToDecimal(num float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(num * shift) / shift
}

func minutesToDaysHoursMinutes(totalMinutes int) string {
	days := totalMinutes / 1440
	hours := (totalMinutes % 1440) / 60
	minutes := totalMinutes % 60

	return fmt.Sprintf("%d days %02d hours %02d minutes", days, hours, minutes)
}

func minutesToHoursMinutes(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours == 1 && minutes == 1 {
		return fmt.Sprintf("%02d hour %02d minute", hours, minutes)
	} else if minutes == 1 {
		return fmt.Sprintf("%02d hours %02d minute", hours, minutes)
	} else if hours == 1{
		return fmt.Sprintf("%02d hour %02d minutes", hours, minutes)
	} else {
		return fmt.Sprintf("%02d hours %02d minutes", hours, minutes)
	}
}

func Stats(context *cli.Context) error { 
	if context.Args().Len() == 0 {
		try := "\"" + "go run . stats <username>" + "\""
		return cli.Exit("Error: No username provided. Try: " + try, 1)
	} else if context.Args().Len() > 1 {
		return cli.Exit("Error: Too many arguments. Usernames are a single argument with no spaces.", 2)
	}

	userName = context.Args().Get(0)

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		userName = "\"" + userName + "\""
		fmt.Println("Username", userName, "not found, re-run the program and try again with an existing username")
		os.Exit(10)
	})

	c.Visit("https://www.letterboxd.com/" + userName + "/films/")

	watchedFilmsPageData()
	individualFilmPageData()
	
	return nil
}

func watchedFilmsPageData() {
	fmt.Println("!!!DISCLAIMER: Films you've watched but haven't rated are excluded from all rating calculations and will be listed below!!!")

	c := colly.NewCollector()

	c.OnHTML("li.poster-container", func(h *colly.HTMLElement) {
		rating, err := symboltoRating(h.Text)
		
		if (err != nil) {
			fmt.Println("Excluded:", h.ChildAttr("img", "alt"))
			filmsWithoutUserRatingCounter++
		}

		filmPage := h.Request.AbsoluteURL("/film/" + h.ChildAttr("div", "data-film-slug"))

		films = append(films, film{
			userRating: rating,
			pageURL: filmPage,
			title: h.ChildAttr("img", "alt"),
		})
	})

	c.OnHTML("[class=next]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(nextPage)
	})

	c.Visit("https://www.letterboxd.com/" + userName + "/films/")

	sumOfUserRatings = sumOfSlice[float64](films, "userRating")
	amountOfFilmsWithUserRatings = float64(len(films) - filmsWithoutUserRatingCounter)
}

// Scrape data from actual film pages, and print out results (film lengths and Letterboxd ratings)
func individualFilmPageData() {
	var filmsWithNoLengthProvided []string
	var filmsWithNoLetterboxdRating []string
	filmCounter := 0

	c := colly.NewCollector()

	c.OnHTML("p.text-link.text-footer", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.Text)
		// Find index of "m" from "mins"
		firstM := strings.Index(text, "m")
		// Slice text to first "m", to get film length
		filmLengthString := strings.TrimSpace(text[0:firstM])
		// Convert film length from string to int if length is valid
		// Otherwise, append film title to filmsWithNoLengthProvided
		if (filmLengthString[0] != 'M') {
			filmLengthInt, err := strconv.Atoi(filmLengthString)

			if err != nil {
				fmt.Println("Error:", err)
			}

			films[filmCounter].length = filmLengthInt
		} else {
			filmsWithNoLengthProvided = append(filmsWithNoLengthProvided, films[filmCounter].title)
		}
	})

	c.OnHTML("meta[name='twitter:data2']", func(h *colly.HTMLElement) {
		// Covert string of form "X.XX out of 5" to float64
		letterboxdRatingString := h.Attr("content")
		letterboxdRating, err := strconv.ParseFloat(letterboxdRatingString[:4], 64)

		if err != nil {
			fmt.Println("Error:", err)
		}

		films[filmCounter].letterboxdRating = letterboxdRating
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode)
		fmt.Println("Error:", err)
	})

	for _, film := range films {
		err := c.Visit(film.pageURL)

		if err != nil {
			fmt.Println("Error during visit:", err)
		}

		filmCounter++
	}

	amountOfFilmsWithLength := len(films) - len(filmsWithNoLengthProvided)
	totalWatchtime := sumOfSlice[int](films, "length")
	averageWatchtime := totalWatchtime / amountOfFilmsWithLength

	fmt.Println("\nTotal watchtime:", minutesToDaysHoursMinutes(totalWatchtime))
	fmt.Println("Average watchtime:", minutesToHoursMinutes(averageWatchtime))

	// Check if film had a letterboxd rating
	for _, film := range films {
		if film.letterboxdRating == 0 && film.userRating != 0 {
			filmsWithNoLetterboxdRating = append(filmsWithNoLetterboxdRating, film.title)
			sumOfUserRatings -= film.userRating
			amountOfFilmsWithUserRatings--
		}
	}

	averageUserRating := sumOfUserRatings / amountOfFilmsWithUserRatings
	amountOfFilmsWithLetterboxdRating := float64(len(films) - len(filmsWithNoLetterboxdRating) - filmsWithoutUserRatingCounter)
	averageLetterboxdRating := sumOfSlice[float64](films, "letterboxdRating") / amountOfFilmsWithLetterboxdRating
	
	fmt.Println("Average user rating:", roundToDecimal(averageUserRating, 2))
	fmt.Println("Average Letterboxd rating of films watched:", roundToDecimal(averageLetterboxdRating, 2))

	if len(filmsWithNoLengthProvided) > 0 {
		fmt.Println()
		fmt.Println("Films with no length provided (Excluded from watchtime calculations):", filmsWithNoLengthProvided)
	}
	
	if len(filmsWithNoLetterboxdRating) > 0 {
		fmt.Println()
		fmt.Println("Films with no Letterboxd rating (Excluded from average letterboxd rating calculations):", filmsWithNoLetterboxdRating)
	}	
}