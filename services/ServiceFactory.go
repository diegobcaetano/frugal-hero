package services

import (
	"errors"
	"frugal-hero/services/LambdaService"
	"frugal-hero/services/S3Service"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getAWSSession() *session.Session {
	return session.Must(
		session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
}

func GetService(name string) (IService, error) {
	switch name {
	case "s3":
		return &S3Service.Service{Session: *getAWSSession()}, nil
	case "lambda":
		return &LambdaService.Service{Session: *getAWSSession()}, nil
	}
	return nil, errors.New("the method requested does not exist")
}
