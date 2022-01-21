package services

import (
	"errors"
	"frugal-hero/services/S3Service"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
		return &S3Service.Service{AwsService: s3.New(getAWSSession())}, nil
	}
	return nil, errors.New("the method requested does not exist")
}
