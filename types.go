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
	buckets []*S3Bucket
	currIdx int
}

func (b *S3BucketPair) NextAlternateObject() (*S3Object, int) {
	idx := b.currIdx
	obj := b.buckets[idx].NextObject()

	if obj != nil {
		b.currIdx = (b.currIdx + 1) % 2
	}

	return obj, idx
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
	obj := b.GetInPageCache(key)
	if obj != nil {
		return obj
	}

	resp, err := b.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &b.name,
		Key:    &key,
	})

	if err != nil {
		return nil
	}

	return NewS3ObjectFromHeadObject(key, resp)
}

func (b *S3Bucket) GetInPageCache(key string) *S3Object {
	for _, e := range b.pageCache {
		if *e.Key == key {
			return NewS3Object(&e)
		}
	}

	return nil
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
		page, err := bucket.paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v", err)
		}

		for _, obj := range page.Contents {
			items = append(items, obj)
		}
	}

	return items
}

var (
	BUCKET_UNCHECKED = new(S3Object)
)

type S3CrossBucketItemMap struct {
	store map[string]([]*S3Object)
}

func NewS3CrossBucketItemMap() *S3CrossBucketItemMap {
	itemMap := new(S3CrossBucketItemMap)
	itemMap.store = make(map[string]([]*S3Object))

	return itemMap
}

func (m *S3CrossBucketItemMap) SetWithItem(item *S3Object, idx int) {
	if item != nil {
		m.SetWithKey(item.key, item, idx)
	}
}

func (m *S3CrossBucketItemMap) SetWithItems(itemA *S3Object, itemB *S3Object) {
	if itemA != nil {
		m.SetWithKey(itemA.key, itemA, 0)
	}
	if itemB != nil {
		m.SetWithKey(itemB.key, itemB, 1)
	}
}

func (m *S3CrossBucketItemMap) SetWithKey(key string, item *S3Object, idx int) {
	if _, ok := m.store[key]; !ok {
		m.store[key] = make([]*S3Object, 2)
		m.store[key][0] = BUCKET_UNCHECKED
		m.store[key][1] = BUCKET_UNCHECKED
	}

	m.store[key][idx] = item
}

func (m *S3CrossBucketItemMap) ObjectKeyExists(key string) bool {
	_, ok := m.store[key]
	return ok
}

func (m *S3CrossBucketItemMap) IsFoundObject(key string, idx int) bool {
	return m.store[key] != nil && m.store[key][idx] != nil && m.store[key][idx] != BUCKET_UNCHECKED
}

func (m *S3CrossBucketItemMap) IsUncheckedObject(key string, idx int) bool {
	return m.store[key] != nil && m.store[key][idx] == BUCKET_UNCHECKED
}

func (m *S3CrossBucketItemMap) IsAbsentObject(key string, idx int) bool {
	return m.store[key] != nil && m.store[key][idx] == nil
}
