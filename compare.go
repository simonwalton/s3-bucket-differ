package main

func crossCorrelateItem(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap, key string) {
	v := itemMap.store[key]

	if itemMap.IsFoundObject(key, 0) && itemMap.IsUncheckedObject(key, 1) {
		obj := buckets.buckets[1].GetObjectMetadata(v[0].key)
		itemMap.SetWithKey(v[0].key, obj, 1)
	}
	if itemMap.IsFoundObject(key, 1) && itemMap.IsUncheckedObject(key, 0) {
		obj := buckets.buckets[0].GetObjectMetadata(v[1].key)
		itemMap.SetWithKey(v[1].key, obj, 0)
	}
}

func crossCorrelate(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap) {
	for k := range itemMap.store {
		crossCorrelateItem(buckets, itemMap, k)
	}
}

func compare(buckets *S3BucketPair, printer func(r ComparisonResult)) {
	itemMap := NewS3CrossBucketItemMap()

	for obj := range buckets.GetAllObjectsAlternately() {
		firstVisit := itemMap.ObjectKeyExists(obj.object.key)
		itemMap.SetWithItem(obj.object, obj.idx)

		crossCorrelate(buckets, itemMap)

		printer(ComparisonResult{itemMap, buckets, firstVisit, obj.object})
	}
}
