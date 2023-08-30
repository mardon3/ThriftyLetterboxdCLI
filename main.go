package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var statsCommand = &cli.Command{
    Name:    "stats",
    Usage:   "Display Letterboxd stats for a user",
    Action:  Stats,
}

var randomCommand = &cli.Command{
    Name:    "random",
    Usage:   "Pick a random film from a user's Letterboxd watchlist",
    Action:  Random,
}

func main() {
    app := &cli.App{
        Name:  "thriftyletterboxd",
        Usage: "Let's you see Letterboxd stats for free or pick a random film from watchlist",
        Commands: []*cli.Command{
            statsCommand,
            randomCommand,
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        panic(err)
    }
}
