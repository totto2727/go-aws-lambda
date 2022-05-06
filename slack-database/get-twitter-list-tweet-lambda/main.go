package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sivchari/gotwtr"
)

var (
	API             = "twitter"
	TWEET_URL       = "https://twitter.com/%s/status/%s"
	CHANNEL         = "la-priere"
	DATETIME_FORMAT = time.RFC3339
	MAX_RESULTS     = 20
)

type ResultItem struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}

type Result struct {
	Channel     string       `json:"channel"`
	Data        []ResultItem `json:"data"`
	API         string       `json:"api"`
	ResultCount int          `json:"result_count"`
	UpdateedAt  string       `json:"updated_at"`
}

func handler() (Result, error) {
	client := gotwtr.New(os.Getenv("BEARER"))
	// log.Println(os.Getenv("BEARER"))

	ts, err := client.LookUpListTweets(context.Background(), "1516932575876378625", &gotwtr.ListTweetsOption{
		MaxResults:  MAX_RESULTS,
		Expansions:  []gotwtr.Expansion{gotwtr.ExpansionAuthorID},
		TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldCreatedAt},
	})
	if err != nil {
		log.Panicln(err)
	}

	authorMap := map[string]*gotwtr.User{}
	for _, u := range ts.Includes.Users {
		authorMap[u.ID] = u
	}

	resultItems := []ResultItem{}
	for _, t := range ts.Tweets {
		author := authorMap[t.AuthorID]
		userName := author.UserName
		URL := fmt.Sprintf(TWEET_URL, userName, t.ID)
		// 2022-05-04T03:41:21.000Z
		createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(URL)
		resultItems = append(resultItems, ResultItem{
			UserName:  userName,
			AuthorID:  t.AuthorID,
			ID:        t.ID,
			CreatedAt: createdAt.Format(DATETIME_FORMAT),
			URL:       URL,
		})
	}

	return Result{
		API:         API,
		Channel:     CHANNEL,
		ResultCount: ts.Meta.ResultCount,
		UpdateedAt:  time.Now().UTC().Format(DATETIME_FORMAT),
		Data:        resultItems,
	}, nil
}

func main() {
	lambda.Start(handler)
}
