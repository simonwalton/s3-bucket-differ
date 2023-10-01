package main

type IBucketComparer interface {
	Compare(*IBucketCrossCorrelator, func(r ComparisonResult))
}

type IBucketCrossCorrelator interface {
	CrossCorrelate(*S3CrossBucketItemMap)
	CrossCorrelateItem(*S3CrossBucketItemMap, string)
}

type BucketComparer struct {
	buckets *S3BucketPair
}

type BucketCrossCorrelator struct {
	buckets *S3BucketPair
}

func (c *BucketCrossCorrelator) CrossCorrelateItem(itemMap *S3CrossBucketItemMap, key string) {
	v := itemMap.store[key]

	if itemMap.IsFoundObject(key, 0) && itemMap.IsUncheckedObject(key, 1) {
		obj := c.buckets.buckets[1].GetObjectMetadata(v[0].key)
		itemMap.SetWithKey(v[0].key, obj, 1)
	}
	if itemMap.IsFoundObject(key, 1) && itemMap.IsUncheckedObject(key, 0) {
		obj := c.buckets.buckets[0].GetObjectMetadata(v[1].key)
		itemMap.SetWithKey(v[1].key, obj, 0)
	}
}

func (c *BucketCrossCorrelator) CrossCorrelate(itemMap *S3CrossBucketItemMap) {
	for k := range itemMap.store {
		c.CrossCorrelateItem(itemMap, k)
	}
}

func (c *BucketComparer) Compare(correlator IBucketCrossCorrelator, printer func(r ComparisonResult)) *S3CrossBucketItemMap {
	itemMap := NewS3CrossBucketItemMap()

	for obj := range c.buckets.GetAllObjectsAlternately() {
		firstVisit := itemMap.ObjectKeyExists(obj.object.key)
		itemMap.SetWithItem(obj.object, obj.idx)

		correlator.CrossCorrelate(itemMap)

		printer(ComparisonResult{itemMap, c.buckets, firstVisit, obj.object})
	}

	return itemMap
}
