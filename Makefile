.PHONY: build

build:
	sam build --use-container

deploy:
	sam build --use-container
	sam deploy

deploy-use-cache:
	sam deploy

test-to-slack-message-lambda:
	sam build --use-container ToSlackMessageLambda
	sam local invoke -e to-slack-message-lambda/event.json ToSlackMessageLambda

test-get-twitter-list-tweet-lambda:
	sam build --use-container GetTwitterListTweet
	sam local invoke GetTwitterListTweet -e get-twitter-list-tweet-lambda/event.json

test-launch-get-message-state-machine-lambda:
	sam build --use-container LaunchGetMessageStateMachineLambda
	sam local invoke LaunchGetMessageStateMachineLambda -e launch-get-message-state-machine-lambda/event.json

