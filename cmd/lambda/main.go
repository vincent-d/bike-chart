package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vincent-d/bike-count/pkg/charts"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	name := request.QueryStringParameters["name"]
	start := request.QueryStringParameters["start"]
	end := request.QueryStringParameters["end"]
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)
	var err error
	if start != "" {
		startDate, err = time.Parse("02-01-2006", start)
		if err != nil {
			log.Print("Invalid start date")
			startDate = endDate.AddDate(0, -6, 0)
		}
	}
	if end != "" {
		endDate, err = time.Parse("02-01-2006", end)
		if err != nil {
			log.Print("Invalid end date")
			endDate = time.Now()
		}
	}

	page, err := charts.GetChartPage(name, startDate, endDate)
	if err != nil {
		log.Print("Error when rendering chart")
		return events.APIGatewayProxyResponse{}, errors.New("error when rendering chart")
	}
	var body string
	if page != nil {
		buf := new(bytes.Buffer)
		page.Render(buf)
		body = buf.String()
	} else {
		body = fmt.Sprintf("<h2>No totem found for name %s</h2>", name)
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            body,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
