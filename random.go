package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/urfave/cli/v2"
)

// Implement a method for selecting a random movie (maybe pick genre)
func Random(context *cli.Context) error {
	if context.Args().Len() < 1 {
		try := "\"" + "go run . stats <username> [genre]" + "\""
		return cli.Exit("Error: No arguments provided. Try: " + try, 3)
	} else if context.Args().Len() > 3 {
		return cli.Exit("Error: Too many arguments. Usernames are a single argument with no spaces.", 4)
	}

	userName = context.Args().Get(0)

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		userName = "\"" + userName + "\""
		fmt.Println("Username", userName, "not found, re-run the program and try again with an existing username")
		os.Exit(10)
	})

	c.Visit("https://www.letterboxd.com/" + userName + "/films/")

	/*  
		If 2 inputs, then genre is one word
		If 3 inputs, then genre is two words
		Let randomGenre know so it can hadnle it
		Else, genre does not matter
	*/ 
	if context.Args().Len() == 2{
		return randomGenre(context, userName, 1)
	} else if context.Args().Len() == 3 {
		return randomGenre(context, userName, 2)
	} else {
		return random(context, userName)
	}
}

func random(context *cli.Context, userName string) error {
	c := colly.NewCollector()

	c.OnHTML("li.poster-container", func(h *colly.HTMLElement) {
		filmTitle := h.ChildAttr("img", "alt")
		filmPage := h.Request.AbsoluteURL("/film/" + h.ChildAttr("div", "data-film-slug"))

		fmt.Println("Title:", filmTitle)

		url := h.Request.AbsoluteURL(h.ChildAttr("div", "data-target-link"))
		filmLengthAndRating(url)

		fmt.Println("Page:", filmPage)
		
		os.Exit(0)
	})

	c.Visit("https://letterboxd.com/" + userName + "/watchlist/by/shuffle/")

	return nil
}

func randomGenre(context *cli.Context, userName string, genreLength int) error {
	c := colly.NewCollector()
	genre := ""

	if genreLength == 1 {
		genre = context.Args().Get(1)
	} else if genreLength == 2 {
		genre = context.Args().Get(1) + "-" + context.Args().Get(2)
	}

	c.OnHTML("li.poster-container", func(h *colly.HTMLElement) {
		filmTitle := h.ChildAttr("img", "alt")
		filmPage := h.Request.AbsoluteURL("/film/" + h.ChildAttr("div", "data-film-slug"))

		fmt.Println("Title:", filmTitle)

		url := h.Request.AbsoluteURL(h.ChildAttr("div", "data-target-link"))
		filmLengthAndRating(url)

		fmt.Println("Page:", filmPage)
		
		os.Exit(0)
	})

	c.OnError(func(r *colly.Response, err error) {
		genre = "\"" + genre + "\""
		fmt.Println("Genre", genre, "not found, re-run the program and try again with a proper genre")
		os.Exit(20)
	})

	c.Visit("https://letterboxd.com/" + userName + "/watchlist/genre/" + genre + "/by/shuffle/")

	return nil
}

// Method for scraping and printing film length and rating
func filmLengthAndRating(url string) {
	c := colly.NewCollector()

	c.OnHTML("meta[name='twitter:data2']", func(h *colly.HTMLElement) {
		letterboxdRating := h.Attr("content")[:4]

		fmt.Println("Letterboxd Rating:", letterboxdRating)
	})

	c.OnHTML("p.text-link.text-footer", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.Text)
		// Find index of "m" from "mins"
		firstM := strings.Index(text, "m")
		// Slice text to first "m", to get film length
		filmLengthString := strings.TrimSpace(text[0:firstM])
		// Convert film length from string to int and print length in hours and minutes
		// Otherwise, print "n/a"
		if (filmLengthString[0] != 'M') {
			filmLengthInt, err := strconv.Atoi(filmLengthString)

			if err != nil {
				fmt.Println("Error:", err)
			}

			fmt.Println("Length:", minutesToHoursMinutes(filmLengthInt))
		} else {
			fmt.Println("Length: n/a")
		}
	})

	c.Visit(url)
}