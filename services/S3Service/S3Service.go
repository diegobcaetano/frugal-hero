package S3Service

import (
	"fmt"
	"frugal-hero/outputs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

type BucketStatus struct {
	bucketName string
	isEmpty    bool
	err awserr.Error
}

type Service struct {
	waitGroup  sync.WaitGroup
	AwsService *s3.S3
}

func (s *Service) getAllBuckets() (*s3.ListBucketsOutput, error) {
	result, err := s.AwsService.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) isBucketEmpty(bucket string, c chan BucketStatus) {
	defer s.waitGroup.Done()
	params := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1),
	}

	obj, objErr := s.AwsService.ListObjects(params)

	if objErr != nil {
		c <- BucketStatus{bucketName: bucket, isEmpty: false, err: objErr.(awserr.Error)}
		return
	}

	if len(obj.Contents) == 0 {
		c <- BucketStatus{bucketName: bucket, isEmpty: true}
	}
	c <- BucketStatus{bucketName: bucket, isEmpty: false}
}

func (s *Service) Inspect(output outputs.OutputInterface) {
	result, err := s.getAllBuckets()
	defer output.Write()

	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return
	}

	fmt.Printf("Total number of buckets: %v\n", len(result.Buckets))
	c := make(chan BucketStatus)

	for _, bucket := range result.Buckets {
		s.waitGroup.Add(1)
		go s.isBucketEmpty(*bucket.Name, c)
	}

	go func() {
		s.waitGroup.Wait()
		close(c)
	}()

	for status := range c {
		if status.isEmpty {
			output.Read([]byte(fmt.Sprintf("Bucket: `%v` \t Status: empty\n", status.bucketName)))
		}
	}
}