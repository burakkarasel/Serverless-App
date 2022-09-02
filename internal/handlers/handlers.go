package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/burakkarasel/Serverless-App/internal/db"
)

var ErrorMethodNotAllowed = "method not allowed"

// ErrorBody holds the error body
type ErrorBody struct {
	ErrorMessage *string `json"error,omitempty"`
}

// GetUser gets the user from database and returns response
func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := db.FetchUser(email, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMessage: aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := db.FetchUsers(tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMessage: aws.String(err.Error())})
	}

	return apiResponse(http.StatusOK, result)
}

// NewUser creates a new user in DB and returns it
func NewUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := db.CreateUser(req, tableName, dynaClient)

	if err != nil {
		apiResponse(http.StatusBadRequest, ErrorBody{ErrorMessage: aws.String(err.Error())})
	}

	return apiResponse(http.StatusOK, result)
}

// UpdateUser updates a specific user and returns it
func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := db.UpdateUser(req, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMessage: aws.String(err.Error())})
	}

	return apiResponse(http.StatusOK, result)
}

// DeleteUser deletes a specific user in DB
func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := db.DeleteUser(req, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMessage: aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)
}

// UnhandledMethod returns error message for unsupported request types
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusBadRequest, ErrorMethodNotAllowed)
}
