package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, m := range sqsEvent.Records {
		log.Println(m.Body)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
