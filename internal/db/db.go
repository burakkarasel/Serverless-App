package db

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/burakkarasel/Serverless-App/internal/models"
	"github.com/burakkarasel/Serverless-App/internal/validators"
)

var (
	ErrorFailedToFetchRecord     = errors.New("failed to fetch record")
	ErrorFailedToUnmarshalRecord = errors.New("failed to unmarshal record")
	ErrorFailedToMarshalRecord   = errors.New("failed to marshal record")
	ErrorFailedToPutRecord       = errors.New("failed to put record")
	ErrorInvalidUserData         = errors.New("invalid user data")
	ErrorInvalidEmailData        = errors.New("invalid email")
	ErrorFailedToDeleteRecord    = errors.New("failed to delete record")
	ErrorUserAlreadyExists       = errors.New("user already exists")
	ErrorUserDoesntExist         = errors.New("user doesnt exists")
)

// FetchUser gets a single user from DB
func FetchUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	// here we create our get function for dynamodb
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	// then we get the record from db
	output, err := dynaClient.GetItem(input)

	// any error occurs we return it
	if err != nil {
		return nil, ErrorFailedToFetchRecord
	}

	// then we unmarshal the result into a user variable
	item := new(models.User)
	err = dynamodbattribute.UnmarshalMap(output.Item, item)

	// if any error occurs we return it
	if err != nil {
		return nil, ErrorFailedToUnmarshalRecord
	}

	// finally we return the record and nil
	return item, nil
}

// FetchUsers gets all the user from the DB
func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]models.User, error) {
	// here we create our get function for dynamodb
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	// then we scan with the input
	output, err := dynaClient.Scan(input)

	if err != nil {
		return nil, ErrorFailedToFetchRecord
	}

	// then we unmarshal list of maps
	items := new([]models.User)
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, items)

	if err != nil {
		return nil, ErrorFailedToUnmarshalRecord
	}

	// finally we return unmarshaled values and nil
	return items, nil
}

// CreateUser creates a new user in DB and returns it
func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	// we unmarshal user data from request
	var u models.User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, ErrorInvalidUserData
	}

	// then we check if email is valid
	if !validators.IsEmailValid(u.Email) {
		return nil, ErrorInvalidEmailData
	}

	// check user if exists
	check, _ := FetchUser(u.Email, tableName, dynaClient)

	if check != nil && len(check.Email) != 0 {
		return nil, ErrorUserAlreadyExists
	}

	// i marshal my user input
	marshaled, err := dynamodbattribute.MarshalMap(u)

	// check for error
	if err != nil {
		return nil, ErrorFailedToMarshalRecord
	}

	// then i create a put item input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      marshaled,
	}

	// then i check for error and output
	_, err = dynaClient.PutItem(input)

	// any error occurs i return it
	if err != nil {
		return nil, ErrorFailedToPutRecord
	}

	// finally i return the created record and nil
	return &u, nil
}

// UpdateUser updates a user in DB for given email
func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*models.User, error) {
	// first we unmarshal user from request
	var u models.User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, ErrorFailedToUnmarshalRecord
	}

	// then we check for given email if user exists
	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)

	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, ErrorUserDoesntExist
	}

	// if its exists we marshal the new user data
	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, ErrorFailedToMarshalRecord
	}

	// then we create input
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	// then we update the user
	_, err = dynaClient.PutItem(input)

	if err != nil {
		return nil, ErrorFailedToPutRecord
	}

	// if no error occurs we return the new user and nil
	return &u, nil
}

// DeleteUser deletes a user in DB for given email
func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	// first we get the email from query
	email := req.QueryStringParameters["email"]

	// then we create the input query
	input := dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	// then we run the query
	_, err := dynaClient.DeleteItem(&input)

	if err != nil {
		return ErrorFailedToDeleteRecord
	}

	// finally if no error occurs we return nil
	return nil
}
