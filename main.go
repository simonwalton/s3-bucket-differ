package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketNameA, bucketNameB string
	maxKeys                  int
)

func init() {
	flag.StringVar(&bucketNameA, "bucket-a", "", "The `name` of the first bucket.")
	flag.StringVar(&bucketNameB, "bucket-b", "", "The `name` of the second bucket.")
	flag.IntVar(&maxKeys, "max-keys", 5, "The maximum number of `keys per page` to retrieve at once.")
}

func main() {
	flag.Parse()

	if len(bucketNameA) == 0 {
		log.Fatal("Please supply the name of bucket A")
	}

	if len(bucketNameB) == 0 {
		log.Fatal("Please supply the name of bucket B")
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	buckets := S3BucketPair{
		buckets: [](*S3Bucket){
			NewBucket(client, bucketNameA, maxKeys),
			NewBucket(client, bucketNameB, maxKeys),
		},
		currIdx: 0,
	}

	compare(&buckets)
}
