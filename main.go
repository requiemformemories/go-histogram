package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
	ErrBadRequest      = errors.New("bad request")
)

type RequestBody struct {
	Values []int `json:"values"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	//JSON(string) -> OBJECT
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       ErrBadRequest.Error(),
			StatusCode: 400,
		}, nil
	}

	// Todo: parse data -> svg
	width := len(body.Values) * 5
	height := max(body.Values...)
	writer := &strings.Builder{}
	canvas := svg.New(writer)
	canvas.Start(width, height)
	for i, v := range body.Values {
		canvas.Rect(i*5, height-v, 5, v, "fill=\"#CBD5E0\"", "onmouseover=\"evt.target.setAttribute('fill', '#A9AED9');\"", "onmouseout=\"evt.target.setAttribute('fill', '#CBD5E0');\"")
	}
	canvas.End()

	return events.APIGatewayProxyResponse{
		Body:       writer.String(),
		Headers:    map[string]string{"Content-Type": "image/svg+xml"},
		StatusCode: 200,
	}, nil

}

func max(nums ...int) int {
	var max int
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

func main() {
	lambda.Start(Handler)
}
