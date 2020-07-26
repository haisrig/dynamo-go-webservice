package db

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var service *dynamodb.DynamoDB

type Employee struct {
	Empid    string
	Name     string
	Age      int
	Location string
}

type EmpData struct {
	Title string
	Data  []Employee
}

func InsertEmployee(empData io.Reader) (string, error) {
	decoder := json.NewDecoder(empData)
	emp := Employee{}
	decoder.Decode(&emp)
	av, err := dynamodbattribute.MarshalMap(emp)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("employees"),
	}
	_, err = getService().PutItem(input)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func ListEmployees() (EmpData, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("employees"),
	}
	result, err := getService().Scan(input)
	if err != nil {
		return EmpData{}, err
	}
	return convertResultToEmployees(result), nil
}

func convertResultToEmployees(result *dynamodb.ScanOutput) EmpData {
	var empList []Employee
	for _, item := range result.Items {
		empList = append(empList, jsonToEmployee(item))
	}
	return EmpData{"Employee Data", empList}
}

func jsonToEmployee(item map[string]*dynamodb.AttributeValue) Employee {
	emp := Employee{}
	err := dynamodbattribute.UnmarshalMap(item, &emp)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return emp
}

func getService() *dynamodb.DynamoDB {
	if service == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		service = dynamodb.New(sess)
	}
	return service
}
