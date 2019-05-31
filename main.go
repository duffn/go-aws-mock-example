package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/urfave/cli"
)

// CloudWatch is a CloudWatch API struct
type CloudWatch struct {
	Client cloudwatchiface.CloudWatchAPI
}

func main() {
	var bucket string

	app := cli.NewApp()
	app.Name = "s3bucketsizegetter"
	app.Usage = "Gets the size of an S3 bucket"

	app.Commands = []cli.Command{
		{
			Name:  "s3",
			Usage: "Gets the size of an S3 bucket",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "bucket",
					Usage:       "name of the S3 bucket",
					Destination: &bucket,
				},
			},
			Action: func(c *cli.Context) error {
				sess := session.Must(session.NewSessionWithOptions(session.Options{
					SharedConfigState: session.SharedConfigEnable,
				}))

				cw := CloudWatch{
					Client: cloudwatch.New(sess),
				}

				result, err := cw.S3BucketSize(bucket)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Bucket %v is approximately %f GB.\n", bucket, result/1e+9)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// S3BucketSize gets the estimated size of an S3 bucket
func (cw *CloudWatch) S3BucketSize(bucket string) (float64, error) {
	result, err := cw.Client.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String("BucketSizeBytes"),
		Namespace:  aws.String("AWS/S3"),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("BucketName"),
				Value: aws.String(bucket),
			},
			&cloudwatch.Dimension{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			},
		},
		Period:     aws.Int64(86400 * 3),
		StartTime:  aws.Time(time.Now().AddDate(0, 0, -3)),
		EndTime:    aws.Time(time.Now()),
		Statistics: aws.StringSlice([]string{"Maximum"}),
	})

	return *result.Datapoints[0].Maximum, err
}
