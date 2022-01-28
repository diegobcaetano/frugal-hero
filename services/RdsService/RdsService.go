package RdsService

import (
	"fmt"
	"frugal-hero/outputs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
	"os"
	"sync"
	"time"
)

type FunctionStatus struct {
	functionName  string
	hasInvocation bool
	err           awserr.Error
}

type Service struct {
	metricsWaitGroup sync.WaitGroup
	dbWaitGroup      sync.WaitGroup
	Session          session.Session
}

func (s *Service) getDatabases(
	rdsService *rds.RDS,
	dbChannel chan *rds.DescribeDBInstancesOutput,
	marker *string) {
	defer s.dbWaitGroup.Done()
	params := &rds.DescribeDBInstancesInput{
		MaxRecords: aws.Int64(100),
		Marker:     marker,
	}

	dbList, dbErr := rdsService.DescribeDBInstances(params)

	if dbErr != nil {
		fmt.Print("Unable to fetch ", dbErr)
		os.Exit(1)
	}

	for i, db := range dbList.DBInstances {
		fmt.Printf("%v - %v \n", i, *db.DBInstanceIdentifier)
	}

	//functionChannel <- dbList

	if dbList.Marker != nil {
		s.dbWaitGroup.Add(1)
		go s.getDatabases(rdsService, dbChannel, dbList.Marker)
	}
}

func (s *Service) getFunctionMetrics(
	cw *cloudwatch.CloudWatch,
	rds *rds.DBInstance,
	functionMetricsChannel chan FunctionStatus) {
	defer s.metricsWaitGroup.Done()
	fmt.Print("oasdweqokdfeoiw")
	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("{AWS/RDS,Per-Database Metrics}"),
		MetricName: aws.String("DatabaseConnections"),
		Statistics: []*string{aws.String("Sum")},
		StartTime:  aws.Time(time.Now().UTC().Add(time.Second * -3600 * 24 * 7)),
		EndTime:    aws.Time(time.Now().UTC()),
		Period:     aws.Int64(3600 * 24 * 7),
		Dimensions: []*cloudwatch.Dimension{{
			Name:  aws.String("FunctionName"),
			Value: rds.DBInstanceIdentifier}},
	}
	result, err := cw.GetMetricStatistics(params)
	if err != nil {
		//functionMetricsChannel <- FunctionStatus{
		//	functionName:  aws.StringValue(rds.FunctionName),
		//	hasInvocation: false,
		//	err:           err.(awserr.Error)}
		//return
	}

	fmt.Println(result)

	//if len(result.Datapoints) == 0 {
	//	functionMetricsChannel <- FunctionStatus{functionName: aws.StringValue(rds.FunctionName), hasInvocation: false}
	//	return
	//}
	//functionMetricsChannel <- FunctionStatus{functionName: aws.StringValue(rds.FunctionName), hasInvocation: true}
}

func (s *Service) Inspect(output outputs.OutputInterface) {
	cwService := cloudwatch.New(&s.Session)
	rdsService := rds.New(&s.Session /* ,&aws.Config{Region: aws.String("us-west-2")}*/)

	//var countNonInvokedFunctions int16
	defer func() {
		//fmt.Printf("Total of non-invoked functions found: %v \n", countNonInvokedFunctions)
		//output.Write()
	}()

	dbChannel := make(chan *rds.DescribeDBInstancesOutput)
	dbMetricsChannel := make(chan FunctionStatus)
	s.dbWaitGroup.Add(1)
	go s.getDatabases(rdsService, dbChannel, nil)

	fmt.Println("Fetching all the functions that are not invoked in the last 7 days...")

	go func() {
		s.dbWaitGroup.Wait()
		close(dbChannel)
	}()

	for dbList := range dbChannel {
		print("wefweweff")
		for _, db := range dbList.DBInstances {
			s.metricsWaitGroup.Add(1)
			s.getFunctionMetrics(cwService, db, dbMetricsChannel)
		}
	}

	//go func() {
	//	s.metricsWaitGroup.Wait()
	//	close(dbMetricsChannel)
	//}()
	//for status := range dbMetricsChannel {
	//	if !status.hasInvocation {
	//		countNonInvokedFunctions++
	//		output.Read([]byte(fmt.Sprintf("Function Name: \t %v\n", status.functionName)))
	//	}
	//}
}
