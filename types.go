package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Object struct {
	key  string
	size int64
}

func NewS3Object(obj *types.Object) *S3Object {
	return &S3Object{
		key:  *obj.Key,
		size: obj.Size,
	}
}

func NewS3ObjectFromHeadObject(key string, obj *s3.HeadObjectOutput) *S3Object {
	return &S3Object{
		key:  key,
		size: obj.ContentLength,
	}
}

type S3BucketPair struct {
	a *S3Bucket
	b *S3Bucket
}

type S3ObjectPair struct {
	pair [2]*S3Object
}

type S3Bucket struct {
	name      string
	paginator *s3.ListObjectsV2Paginator
	pageCache []types.Object
	pageIdx   int
	client    *s3.Client
}

func NewBucket(client *s3.Client, bucketName string, maxKeys int) *S3Bucket {
	b := new(S3Bucket)
	b.name = bucketName
	b.pageIdx = 0
	b.client = client

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

func (b *S3Bucket) GetObjectMetadata(key string) *S3Object {
	resp, err := b.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &b.name,
		Key:    &key,
	})

	if err != nil {
		return nil
	}

	return NewS3ObjectFromHeadObject(key, resp)
}

func (b *S3Bucket) NextObject() *S3Object {
	if b.pageIdx == len(b.pageCache) {
		InsertNextPageIntoCache(b)
		b.pageIdx = 0
	}

	if len(b.pageCache) == 0 {
		return nil
	}

	obj := &b.pageCache[b.pageIdx]
	b.pageIdx++

	return NewS3Object(obj)
}

func InsertNextPageIntoCache(bucket *S3Bucket) {
	bucket.pageCache = NextObjects(bucket)
}

func NextObjects(bucket *S3Bucket) []types.Object {
	items := make([]types.Object, 0)

	if bucket.paginator.HasMorePages() {
		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := bucket.paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v", err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			items = append(items, obj)
		}
	}

	return items
}

type S3CrossBucketItemMap struct {
	store map[string]([]*S3Object)
}

func NewS3CrossBucketItemMap() *S3CrossBucketItemMap {
	itemMap := new(S3CrossBucketItemMap)
	itemMap.store = make(map[string]([]*S3Object))
	return itemMap
}

func (m *S3CrossBucketItemMap) Set(item *S3Object, idx int) {
	if item == nil {
		return
	}

	if _, ok := m.store[item.key]; !ok {
		m.store[item.key] = make([]*S3Object, 2)
	}

	m.store[item.key][idx] = item
}
