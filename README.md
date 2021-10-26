# Build command for lambda function
GOOS=linux GOARCH=amd64 go build -o main main.go logger.go get_csv.go DynamoDBCreateItem.go

#Local Development
sam local invoke "ArkInvestFunction" -e ./events/events.json
You can modify the events to get different functionality

sam build
Build the lambda taht will be run in the container of choice

Install DynamoDB and extract zip to your folder of choice and run the 
following to get a local version of DynamoDB running (run from the folder of choice):

java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb 

windows powershell command
java -D"java.library.path=./DynamoDBLocal_lib" -jar DynamoDBLocal.jar -sharedDb 