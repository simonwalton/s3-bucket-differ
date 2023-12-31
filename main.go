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
	flag.IntVar(&maxKeys, "max-keys", 50, "The maximum number of `keys per page` to retrieve at once.")
}

func main() {
	flag.Parse()

	if len(bucketNameA) == 0 {
		log.Fatal("Please supply the name of bucket A")
	}

	if len(bucketNameB) == 0 {
		log.Fatal("Please supply the name of bucket B")
	}

	// load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	buckets := S3BucketPair{
		buckets: [](*S3Bucket){
			NewBucket(client, bucketNameA, maxKeys),
			NewBucket(client, bucketNameB, maxKeys),
		},
		currIdx: 0,
	}

	doOutput := func(r ComparisonResult) {
		printSummary(r.buckets, r.itemMap)

		if r.firstVisit {
			appendDetail(r.currObject.key, r.itemMap)
		}
	}

	comparer := new(BucketComparer)
	comparer.buckets = &buckets
	correlator := new(BucketCrossCorrelator)
	correlator.buckets = &buckets

	comparer.Compare(correlator, doOutput)
}
