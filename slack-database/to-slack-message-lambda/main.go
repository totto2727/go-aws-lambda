package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

func handler(ctx context.Context, event events.SNSEvent) error {
	api := slack.New(os.Getenv("SLACK_API"))

	for _, m := range event.Records {
		attribute := m.SNS.MessageAttributes["channel"]
		channelAttribute := attribute.(map[string]interface{})
		channelAttributeValue := channelAttribute["Value"].(string)
		channelID := path.Base(channelAttributeValue)
		createdAtAttribute := m.SNS.MessageAttributes["created_at"]
		createdAtAttributeMap := createdAtAttribute.(map[string]interface{})
		createdAtAttributeValue := createdAtAttributeMap["Value"].(string)
		log.Println("channel: " + channelID)

		_, _, err := api.PostMessage(
			channelID,
			slack.MsgOptionText(createdAtAttributeValue+"\n"+m.SNS.Message, false),
			slack.MsgOptionAsUser(true),
		)
		if err != nil {
			log.Panicln(err)
		}

		log.Println(channelID + ": " + m.SNS.Message)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
