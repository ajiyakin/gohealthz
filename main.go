package main

import (
	"log"

	sdk "github.com/gaia-pipeline/gosdk"
)

func sayHello(args sdk.Arguments) error {
	log.Println("Hello world pipeline from GoHealthz's pipeline")
	return nil
}

func main() {
	jobs := sdk.Jobs{
		sdk.Job{
			Handler:     sayHello,
			Title:       "Say Hello",
			Description: "This job prints greeting message",
		},
	}
	if err := sdk.Serve(jobs); err != nil {
		panic(err)
	}
}
