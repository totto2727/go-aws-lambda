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
	TWEET_URL       = "https://twitter.com/%s/status/%s"
	DATETIME_FORMAT = time.RFC3339
)

type MyEvent struct {
	TwitterList string `json:"twitter_list"`
	UpdatedAt   string `json:"updated_at"`
	MaxResults  int    `json:"max_results"`
}

type ResultItem struct {
	ID        string `json:"id"`
	UserID    string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}

type Result struct {
	Data         []ResultItem `json:"data"`
	MessageCount int          `json:"message_count"`
	UpdatedAt    string       `json:"updated_at"`
}

func handler(ctx context.Context, event MyEvent) (Result, error) {
	client := gotwtr.New(os.Getenv("BEARER"))
	lastUpdatedAt, err := time.Parse(DATETIME_FORMAT, event.UpdatedAt)
	if err != nil {
		log.Panicln(err)
	}

	getListTweets := func(paginationToken string) (*gotwtr.ListTweetsResponse, error) {
		var paginationToken_ string
		if len(paginationToken) != 0 {
			paginationToken_ = paginationToken
		}
		return client.LookUpListTweets(context.Background(), event.TwitterList, &gotwtr.ListTweetsOption{
			PaginationToken: paginationToken_,
			MaxResults:      event.MaxResults,
			Expansions:      []gotwtr.Expansion{gotwtr.ExpansionAuthorID},
			TweetFields:     []gotwtr.TweetField{gotwtr.TweetFieldCreatedAt},
		})
	}

	getResultItems := func(listTweetsResponse *gotwtr.ListTweetsResponse) ([]ResultItem, bool) {
		authorMap := map[string]*gotwtr.User{}
		for _, u := range listTweetsResponse.Includes.Users {
			authorMap[u.ID] = u
		}

		resultItems := []ResultItem{}
		for _, t := range listTweetsResponse.Tweets {
			author := authorMap[t.AuthorID]
			userName := author.UserName
			url := fmt.Sprintf(TWEET_URL, userName, t.ID)
			// 2022-05-04T03:41:21.000Z
			createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
			if err != nil {
				log.Panicln(err)
			}

			if !createdAt.After(lastUpdatedAt) {
				return resultItems, true
			}

			createdAtFormatted := createdAt.Format(DATETIME_FORMAT)
			resultItems = append(resultItems, ResultItem{
				ID:        t.ID,
				UserID:    t.AuthorID,
				UserName:  author.UserName,
				CreatedAt: createdAtFormatted,
				URL:       url,
			})
			log.Printf("%s %s", createdAtFormatted, url)
		}
		return resultItems, false
	}

	listTweetsResponse, err := getListTweets("")
	if err != nil {
		log.Panicln(err)
	}

	resultItems, isLast := getResultItems(listTweetsResponse)

	if !isLast {
		for {
			log.Println("Next Page")
			time.Sleep(1 * time.Second)
			listTweetsResponse, err = getListTweets(listTweetsResponse.Meta.NextToken)
			if err != nil {
				log.Panicln(err)
			}
			nextResultItems, isLast := getResultItems(listTweetsResponse)
			resultItems = append(resultItems, nextResultItems...)
			if isLast {
				break
			}
		}
	}

	tweetCount := len(resultItems)
	if tweetCount == 0 {
		return Result{
			MessageCount: tweetCount,
			UpdatedAt:    event.UpdatedAt,
			Data:         []ResultItem{},
		}, nil
	}

	for i, j := 0, tweetCount-1; i < j; i, j = i+1, j-1 {
		resultItems[i], resultItems[j] = resultItems[j], resultItems[i]
	}

	updatedAt := resultItems[tweetCount-1].CreatedAt
	log.Printf("%s %s %d", event.UpdatedAt, updatedAt, tweetCount)
	return Result{
		MessageCount: tweetCount,
		UpdatedAt:    updatedAt,
		Data:         resultItems,
	}, nil
}

func main() {
	lambda.Start(handler)
}
