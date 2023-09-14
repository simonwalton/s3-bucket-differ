package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Object struct {
	key string
	sha string
}

type S3Bucket struct {
	name      string
	paginator *s3.ListObjectsV2Paginator
}

func (b *S3Bucket) HasMoreItems() bool {
	return b.paginator.HasMorePages()
}

func NewBucket(client *s3.Client, bucketName string, maxKeys int) *S3Bucket {
	b := new(S3Bucket)
	b.name = bucketName

	params := &s3.ListObjectsV2Input{
		Bucket: &b.name,
	}
	b.paginator = s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})
	return b
}

func (b *S3Bucket) NextObjects() []types.Object {
	var i int
	items := make([]types.Object, 0)

	if b.paginator.HasMorePages() {
		i++

		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := b.paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			items = append(items, obj)
		}
	}

	return items
}
