package main

import "testing"

func Test_compare(t *testing.T) {
	type args struct {
		bucketA *S3Bucket
		bucketB *S3Bucket
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compare(tt.args.bucketA, tt.args.bucketB)
		})
	}
}
