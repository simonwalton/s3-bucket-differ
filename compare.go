package main

import (
	"fmt"
	"os"
)

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

func compare(buckets *S3BucketPair) {
	itemMap := NewS3CrossBucketItemMap()
	itemsRemain := true

	for itemsRemain {
		obj, idx := buckets.NextAlternateObject()
		if obj == nil {
			break
		}

		alreadyVisited := itemMap.ObjectKeyExists(obj.key)
		itemMap.SetWithItem(obj, idx)

		crossCorrelate(buckets, itemMap)
		printSummary(buckets, itemMap)

		if !alreadyVisited {
			appendDetail(obj.key, itemMap)
		}

		itemsRemain = obj != nil
	}

	crossCorrelate(buckets, itemMap)
	printSummary(buckets, itemMap)
}

func printObjectList(items [](*S3Object)) {
	for _, item := range items {
		fmt.Fprintf(os.Stderr, " - %s\n", item.key)
	}
}

func appendDetail(key string, itemMap *S3CrossBucketItemMap) {
	fmt.Printf("%s\t%t\t%t\n", key, itemMap.store[key][0] != nil, itemMap.store[key][1] != nil)
}

func printSummary(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap) {
	common := 0
	onlyA := make([](*S3Object), 0)
	onlyB := make([](*S3Object), 0)

	for _, v := range itemMap.store {
		if v[0] != nil && v[1] != nil {
			common++
		}
		if v[0] != nil && v[1] == nil {
			onlyA = append(onlyA, v[0])
		}
		if v[0] == nil && v[1] != nil {
			onlyB = append(onlyB, v[1])
		}
	}

	fmt.Fprintf(os.Stderr, "\x1b[2J\n")
	fmt.Fprintf(os.Stderr, "🪣  Bucket Comparison\n")
	fmt.Fprintf(os.Stderr, "  A: %s\n", buckets.buckets[0].name)
	fmt.Fprintf(os.Stderr, "  B: %s\n", buckets.buckets[1].name)

	fmt.Fprintf(os.Stderr, "\nIn common: %d\n", common)
	fmt.Fprintf(os.Stderr, "Only A: (%d)\n", len(onlyA))
	printObjectList(onlyA)
	fmt.Fprintf(os.Stderr, "Only B: (%d)\n", len(onlyB))
	printObjectList(onlyB)
}
