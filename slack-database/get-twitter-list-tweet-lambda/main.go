package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sivchari/gotwtr"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	client := gotwtr.New(os.Getenv("BEARER"))
	log.Println(os.Getenv("BEARER"))
	// look up multiple tweets
	ts, err := client.SearchRecentTweets(context.Background(), "藍月なくる")
	if err != nil {
		panic(err)
	}
	for _, t := range ts.Tweets {
		log.Println(t)
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintln(ts),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
