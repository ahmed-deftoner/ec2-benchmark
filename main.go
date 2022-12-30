package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Todo struct {
	Id          string `json:"id" dynamodbav:"id"`
	Name        string `json:"name" dynamodbav:"name"`
	Description string `json:"description" dynamodbav:"description"`
	Status      bool   `json:"status" dynamodbav:"status"`
}

const TableName = "Todos"

var db dynamodb.Client
var ctx context.Context

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.TODO()
	db = *dynamodb.NewFromConfig(sdkConfig)
}

func listItems() ([]Todo, error) {
	todos := make([]Todo, 0)
	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedTodos []Todo
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedTodos)
		if err != nil {
			return nil, err
		}

		todos = append(todos, fetchedTodos...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}
	}

	return todos, nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Print("Received GET todos request on 8000")

	todos, err := listItems()
	if err != nil {
		log.Fatal("error while getting todos")
	}

	json.NewEncoder(w).Encode(todos)
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/todo", homePage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
