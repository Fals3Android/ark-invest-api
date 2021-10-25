# Build command for lambda function
GOOS=linux GOARCH=amd64 go build -o main main.go logger.go get_csv.go

#Local Development
sam local invoke "ArkInvestFunction" -e ./events/events.json
You can modify the events to get different functionality