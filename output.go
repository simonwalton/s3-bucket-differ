package main

import (
	"fmt"
	"os"
)

const MAX_OBJECT_LIST_HEAD = 5

func writeDetailHeader(buckets *S3BucketPair) {
	fmt.Printf("%s\t%s\t%s\n", "key", buckets.buckets[0].name, buckets.buckets[1].name)
}

func appendDetail(key string, itemMap *S3CrossBucketItemMap) {
	if itemMap.IsOnlyOneObjectFound(key) {
		fmt.Printf("%s\t%t\t%t\n", key, itemMap.store[key][0] != nil, itemMap.store[key][1] != nil)
	}
}

func printSummary(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap) {
	common := 0
	onlyA := make([](*S3Object), 0)
	onlyB := make([](*S3Object), 0)

	for k, v := range itemMap.store {
		if itemMap.AreBothFoundObjects(k) {
			common++
		}
		if itemMap.IsOnlyGivenObjectFound(k, 0) {
			onlyA = append(onlyA, v[0])
		}
		if itemMap.IsOnlyGivenObjectFound(k, 1) {
			onlyB = append(onlyB, v[1])
		}
	}

	fmt.Fprintf(os.Stderr, "\x1b[2J\n")
	fmt.Fprintf(os.Stderr, "🪣  Bucket Comparison\n")
	fmt.Fprintf(os.Stderr, "  A: %s\n", buckets.buckets[0].name)
	fmt.Fprintf(os.Stderr, "  B: %s\n", buckets.buckets[1].name)

	fmt.Fprintf(os.Stderr, "\nIn common: %d\n", common)
	fmt.Fprintf(os.Stderr, "Only A: (%d)\n", len(onlyA))

	printObjectList(onlyA[:min(len(onlyA), MAX_OBJECT_LIST_HEAD)])

	fmt.Fprintf(os.Stderr, "Only B: (%d)\n", len(onlyB))
	printObjectList(onlyB[:min(len(onlyB), MAX_OBJECT_LIST_HEAD)])
}

func printObjectList(items [](*S3Object)) {
	for _, item := range items {
		fmt.Fprintf(os.Stderr, " - %s\n", item.key)
	}
}
