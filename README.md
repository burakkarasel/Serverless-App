# Serverless App

---

### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [Author Info](#author-info)

---

## Description

Aim of this project is creating a DynamoDB table, Lambda Function, API Gateway and serving HTTP requests for creating, updating, getting and deleting user.

## Technologies

### Main Technologies

- [Go](https://go.dev/)
- [AWS API Gateway](https://aws.amazon.com/api-gateway/)
- [AWS Lambda](https://aws.amazon.com/lambda/)
- [AWS DynamoDB](https://aws.amazon.com/dynamodb/)

### Libraries

- [aws/aws-lambda-go](https://github.com/aws/aws-lambda-go)
- [aws/aws-sdk-go](https://github.com/aws/aws-sdk-go)
- [asaskevich/govalidator](https://github.com/asaskevich/govalidator)

[Back To The Top](#Serverless-App)

---

## How To Use

### Tools

- [Go](https://go.dev/dl/)
- [AWS](https://aws.amazon.com/)

### Build and Zip the App

- Build

```
go build cmd/web/*.go
```

- Zip

```
zip -jrm build/main.zip build/main
```

## Deploy The App

- First Create an AWS account if you don't have one from [here](https://signin.aws.amazon.com/)

- Then Create yourself a new Lambda function with the zipped app

- Then Create a new table in DynamoDB

- Then Create a new API Gateway for the app and deploy it!

## Author Info

- Twitter - [@dev_bck](https://twitter.com/dev_bck)

[Back To The Top](#Serverless-App)
