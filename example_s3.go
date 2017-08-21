package main

import (
    "flag"

    "github.com/docker/goamz/aws"
    "github.com/docker/goamz/s3"
)

var (
    access_key string
    secret_key string
    bucketName string
    fileName   string
)

func init() {
    flag.StringVar(&access_key, "a", "", "AWS_ACCESS_KEY_ID")
    flag.StringVar(&secret_key, "s", "", "AWS_SECRET_ACCESS_KEY")
    flag.StringVar(&bucketName, "b", "", "Bucket Name")
}

func main() {

    flag.Parse()

    // The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
    auth, err := aws.EnvAuth()
    if err != nil {
        panic(err.Error())
    }

    // Open Bucket
    // s := s3.New(auth, aws.APNortheast)
    s := s3.New(auth, aws.USWest2)
    bucket := s.Bucket(bucketName)

    data := []byte("Hello, Goamz!!")
    err = bucket.Put("/test-upload/sample.txt", data, "text/plain", s3.BucketOwnerFull, s3.Options{})
    if err != nil {
        panic(err.Error())
    }
}
