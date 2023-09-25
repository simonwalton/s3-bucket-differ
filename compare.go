package main

func crossCorrelate(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap) {
	for k, v := range itemMap.store {
		if itemMap.IsFoundObject(k, 0) && itemMap.IsUncheckedObject(k, 1) {
			obj := buckets.buckets[1].GetObjectMetadata(v[0].key)
			itemMap.SetWithKey(v[0].key, obj, 1)
		}
		if itemMap.IsFoundObject(k, 1) && itemMap.IsUncheckedObject(k, 0) {
			obj := buckets.buckets[0].GetObjectMetadata(v[1].key)
			itemMap.SetWithKey(v[1].key, obj, 0)
		}
	}
}

func compare(buckets *S3BucketPair) chan ComparisonResult {
	ch := make(chan ComparisonResult)
	itemMap := NewS3CrossBucketItemMap()
	obj, idx := buckets.NextAlternateObject()

	go func(ch chan ComparisonResult) {
		for obj != nil {
			firstVisit := itemMap.ObjectKeyExists(obj.key)
			itemMap.SetWithItem(obj, idx)

			crossCorrelate(buckets, itemMap)

			ch <- ComparisonResult{itemMap, buckets, firstVisit, obj}

			obj, idx = buckets.NextAlternateObject()
		}

		close(ch)
	}(ch)

	return ch
}
