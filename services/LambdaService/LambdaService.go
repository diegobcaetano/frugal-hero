package LambdaService

import (
	"fmt"
	"frugal-hero/outputs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/lambda"
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
	metricsWaitGroup   sync.WaitGroup
	functionsWaitGroup sync.WaitGroup
	Session            session.Session
}

func (s *Service) getFunctions(
	lambdaService *lambda.Lambda,
	functionChannel chan *lambda.ListFunctionsOutput,
	marker *string) {
	defer s.functionsWaitGroup.Done()
	params := &lambda.ListFunctionsInput{
		MaxItems: aws.Int64(500),
		Marker:   marker,
	}

	lambdaList, lambdaErr := lambdaService.ListFunctions(params)

	if lambdaErr != nil {
		fmt.Print("Unable to fetch ", lambdaErr)
		os.Exit(1)
	}

	functionChannel <- lambdaList

	if lambdaList.NextMarker != nil {
		s.functionsWaitGroup.Add(1)
		go s.getFunctions(lambdaService, functionChannel, lambdaList.NextMarker)
	}
}

func (s *Service) getFunctionMetrics(
	cw *cloudwatch.CloudWatch,
	f *lambda.FunctionConfiguration,
	functionMetricsChannel chan FunctionStatus) {

	defer s.metricsWaitGroup.Done()

	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		Statistics: []*string{aws.String("Sum")},
		StartTime:  aws.Time(time.Now().UTC().Add(time.Second * -3600 * 24 * 7)),
		EndTime:    aws.Time(time.Now().UTC()),
		Period:     aws.Int64(3600 * 24 * 7),
		Dimensions: []*cloudwatch.Dimension{{
			Name:  aws.String("FunctionName"),
			Value: f.FunctionName}},
	}
	result, err := cw.GetMetricStatistics(params)
	if err != nil {
		functionMetricsChannel <- FunctionStatus{
			functionName:  aws.StringValue(f.FunctionName),
			hasInvocation: false,
			err:           err.(awserr.Error)}
		return
	}

	if len(result.Datapoints) == 0 {
		functionMetricsChannel <- FunctionStatus{functionName: aws.StringValue(f.FunctionName), hasInvocation: false}
		return
	}
	functionMetricsChannel <- FunctionStatus{functionName: aws.StringValue(f.FunctionName), hasInvocation: true}
}

func (s *Service) Inspect(output outputs.OutputInterface) {
	cwService := cloudwatch.New(&s.Session)
	lambdaService := lambda.New(&s.Session /* ,&aws.Config{Region: aws.String("us-west-2")}*/)

	var countNonInvokedFunctions int16
	defer func() {
		fmt.Printf("Total of non-invoked functions found: %v \n", countNonInvokedFunctions)
		output.Write()
	}()

	functionsChannel := make(chan *lambda.ListFunctionsOutput)
	functionMetricsChannel := make(chan FunctionStatus)
	s.functionsWaitGroup.Add(1)
	go s.getFunctions(lambdaService, functionsChannel, nil)

	fmt.Println("Fetching all the functions that are not invoked in the last 7 days...")

	go func() {
		s.functionsWaitGroup.Wait()
		close(functionsChannel)
	}()

	for functionList := range functionsChannel {
		for _, f := range functionList.Functions {
			s.metricsWaitGroup.Add(1)
			go s.getFunctionMetrics(cwService, f, functionMetricsChannel)
		}
	}

	go func() {
		s.metricsWaitGroup.Wait()
		close(functionMetricsChannel)
	}()
	for status := range functionMetricsChannel {
		if !status.hasInvocation {
			countNonInvokedFunctions++
			output.Read([]byte(fmt.Sprintf("Function Name: \t %v\n", status.functionName)))
		}
	}
}
