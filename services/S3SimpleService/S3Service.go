package S3SimpleService

import (
	"fmt"
	"frugal-hero/outputs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BucketStatus struct {
	bucketName string
	isEmpty    bool
	err        awserr.Error
}

type Service struct {
	Session session.Session
}

func (s *Service) getAllBuckets(s3Service *s3.S3) (*s3.ListBucketsOutput, error) {
	result, err := s3Service.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) isBucketEmpty(s3Service *s3.S3, bucket string) BucketStatus {

	params := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1),
	}

	obj, objErr := s3Service.ListObjects(params)

	if objErr != nil {
		return BucketStatus{bucketName: bucket, isEmpty: false, err: objErr.(awserr.Error)}
	}

	if len(obj.Contents) == 0 {
		fmt.Printf("Bucket: `%v` \t Status: empty\n", bucket)
		return BucketStatus{bucketName: bucket, isEmpty: true}
	}
	return BucketStatus{bucketName: bucket, isEmpty: false}
}

func (s *Service) Inspect(output outputs.OutputInterface) {
	defer output.Write()
	s3Service := s3.New(&s.Session)
	result, err := s.getAllBuckets(s3Service)

	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return
	}

	fmt.Printf("Total number of buckets: %v\n", len(result.Buckets))
	fmt.Println("Fetching all the empty buckets...")

	for _, bucket := range result.Buckets {
		s.isBucketEmpty(s3Service, *bucket.Name)
	}
}
