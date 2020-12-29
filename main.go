package main

import (
	"fmt"
	"log"
	"os"

	"github.com/simse/hermes/cmd"
	"github.com/simse/hermes/internal/assets"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hermes",
		Usage: "a tool for deploying to and managing websites on S3+CloudFront",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Creates a hermes stack",
				Action: cmd.InitCommand,
			},
		},
	}

	// Set up assets box
	assets.InitBox()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
