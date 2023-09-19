package main

import (
	"fmt"
	"os"
)

func compare(buckets *S3BucketPair) {
	itemMap := NewS3CrossBucketItemMap()
	itemsRemain := true

	for itemsRemain {
		itemA := buckets.a.NextObject()
		itemB := buckets.b.NextObject()

		itemMap.Set(itemA, 0)
		itemMap.Set(itemB, 1)

		drawSummary(buckets, itemMap)

		itemsRemain = itemA != nil || itemB != nil
	}
}

func printObjectList(items [](*S3Object)) {
	for _, item := range items {
		fmt.Fprintf(os.Stderr, " - %s\n", item.key)
	}
}

func drawSummary(buckets *S3BucketPair, itemMap *S3CrossBucketItemMap) {
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
	fmt.Fprintf(os.Stderr, "  A: %s\n", buckets.a.name)
	fmt.Fprintf(os.Stderr, "  B: %s\n", buckets.b.name)

	fmt.Fprintf(os.Stderr, "\nIn common: %d\n", common)
	fmt.Fprintf(os.Stderr, "Only A: (%d)\n", len(onlyA))
	printObjectList(onlyA)
	fmt.Fprintf(os.Stderr, "Only B: (%d)\n", len(onlyB))
	printObjectList(onlyB)
}
