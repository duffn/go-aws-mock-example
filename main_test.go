package main

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type mockedCloudWatchMetrics struct {
	cloudwatchiface.CloudWatchAPI
	Resp cloudwatch.GetMetricStatisticsOutput
}

func (m mockedCloudWatchMetrics) GetMetricStatistics(in *cloudwatch.GetMetricStatisticsInput) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return &m.Resp, nil
}

func TestS3BucketSize(t *testing.T) {
	expected := 5.1428e+10
	mockedResp := cloudwatch.GetMetricStatisticsOutput{
		Datapoints: []*cloudwatch.Datapoint{
			{
				Maximum:   aws.Float64(5.1428e+10),
				Timestamp: aws.Time(time.Now()),
				Unit:      aws.String("Bytes"),
			},
		},
		Label: aws.String("BucketSizeBytes"),
	}

	cw := CloudWatch{
		Client: mockedCloudWatchMetrics{Resp: mockedResp},
	}

	resp, err := cw.S3BucketSize("my-bucket")

	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if expected != resp {
		t.Errorf("expected %v, go %v", expected, resp)
	}
}
