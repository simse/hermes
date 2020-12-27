package main

import (
	"fmt"
	"log"
	"os"

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
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
