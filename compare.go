package main

import "fmt"

var (
	store = make(map[string]([]*S3Object))
)

func storeItem(item *S3Object, idx int) {
	if item == nil {
		return
	}

	if _, ok := store[item.key]; !ok {
		store[item.key] = make([]*S3Object, 2)
	}

	store[item.key][idx] = item
}

// TODO: compare should not have to care that buckets are paged, so we need the bucket to simply return a next item and we operate one at a time
func compare(bucketA *S3Bucket, bucketB *S3Bucket) {
	itemsRemain := true

	for itemsRemain {
		itemA := bucketA.NextObject()
		itemB := bucketB.NextObject()

		storeItem(itemA, 0)
		storeItem(itemB, 1)

		if itemA != nil {
			println("A " + itemA.key)
		}
		if itemB != nil {
			println("B " + itemB.key)
		}

		drawSummary()

		itemsRemain = itemA != nil || itemB != nil
	}
}

func printObjectList(items [](*S3Object)) {
	for _, item := range items {
		println("- " + item.key)
	}
}

func drawSummary() {
	common := 0
	onlyA := make([](*S3Object), 0)
	onlyB := make([](*S3Object), 0)
	for _, v := range store {
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

	fmt.Printf("\x1b[2J")
	fmt.Printf("In common %d\n", common)
	fmt.Printf("Only A\n")
	printObjectList(onlyA)
	fmt.Printf("Only B\n")
	printObjectList(onlyB)
}
