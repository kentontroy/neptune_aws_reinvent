package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	defaultLocalSource  = "../../../data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
	defaultAWSConfig    = "317913635185_cldr_poweruser"
	defaultAWSRegion    = "us-east-2"
	defaultAWSBucket    = "kdavis-bucket"
	defaultAWSBucketKey = "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
)

func main() {
	source := flag.String("source", defaultLocalSource, "The full path to the source file")
	awsConfig := flag.String("aws_config", defaultAWSConfig, "The section in the AWS Config file")
	awsRegion := flag.String("aws_region", defaultAWSRegion, "The AWS Region")
	awsBucket := flag.String("aws_bucket", defaultAWSBucket, "The AWS Bucket")
	awsBucketKey := flag.String("aws_bucket_key", defaultAWSBucketKey, "The AWS Bucket Key")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(*awsRegion),
		config.WithSharedConfigProfile(*awsConfig),
	)
	if err != nil {
		log.Fatalf("Unable to load the AWS SDK config, %v", err)
	}
	log.Println("Successful connection to AWS")

	s3Client := s3.NewFromConfig(cfg)

	file, err := os.Open(*source)
	if err != nil {
		log.Fatalf("Failed to open file %q, %v", source, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	input := &s3.PutObjectInput{
		Bucket:        aws.String(*awsBucket),
		Key:           aws.String(*awsBucketKey),
		Body:          file,
		ContentLength: &fileSize,
		ContentType:   aws.String("text/csv"),
		ACL:           types.ObjectCannedACLPublicRead,
	}

	_, err = s3Client.PutObject(context.TODO(), input)
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", "data.csv", awsBucket, err)
	}

	fmt.Printf("Successfully uploaded %s to %s/%s\n", *source, *awsBucket, *awsBucketKey)
}
