package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func  main2()  {
	app := cli.NewApp()
	app.Name = "ThrifyLetterboxd"
	app.Usage = "Let's you see Letterboxd stats for free or pick a random film from watchlist"
	app.Commands = []*cli.Command{
		{
			Name:        "stats",
			HelpName:    "stats",
			Action:      Stats,
			ArgsUsage:   "<username>",
			Usage:       "mandatory username <username>",
			Description: "Scrapes and displays Letterboxd stats of given username's watched films.",
		},
		{
			Name:        "random",
			HelpName:    "random",
			Action:      Random,
			ArgsUsage:   "<username> [genre]",
			Usage: 		 "mandatory username, optional genre <username> [genre] ",
			Description: "Picks random film from given username's Letterboxd watchlist.",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}