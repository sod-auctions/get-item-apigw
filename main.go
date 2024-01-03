package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sod-auctions/auctions-db"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

var database *auctions_db.Database

func init() {
	log.SetFlags(0)
	var err error
	database, err = auctions_db.NewDatabase(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
}

type Item struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	MediaURL string `json:"mediaUrl"`
	Quality  string `json:"quality"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, _ := strconv.Atoi(event.QueryStringParameters["id"])

	item, err := database.GetItem(int32(id))
	if err != nil {
		log.Printf("An error occurred: %v\n", err)

		errorMessage := ErrorMessage{Error: "An internal error occurred"}
		body, _ := json.Marshal(errorMessage)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "http://localhost:3000",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
			},
			Body: string(body),
		}, nil
	}

	mItem := &Item{
		Id:       item.Id,
		Name:     item.Name,
		MediaURL: item.MediaURL,
		Quality:  item.Rarity,
	}

	body, _ := json.Marshal(mItem)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "http://localhost:3000",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
