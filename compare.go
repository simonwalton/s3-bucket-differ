package main

func compare(bucketA *S3Bucket, bucketB *S3Bucket) {
	itemsRemain := true

	for itemsRemain {
		itemsA := bucketA.NextObjects()
		itemsB := bucketB.NextObjects()

		for _, obj := range itemsA {
			println(*obj.Key)
		}

		for _, obj := range itemsB {
			println(*obj.Key)
		}

		itemsRemain = bucketA.HasMoreItems() || bucketB.HasMoreItems()
	}
}
