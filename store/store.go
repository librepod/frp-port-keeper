package store

import (
	"encoding/json"
	"fmt"
	"time"

	// "os"

	"github.com/philippgille/gokv"

	// "github.com/philippgille/gokv/file"
	// "github.com/philippgille/gokv/encoding"
	// "github.com/philippgille/gokv/s3"
	"github.com/philippgille/gokv/redis"
)

type MachineRecord struct {
	Port      int
	IP        string
	CreatedAt time.Time
}

type PortRecord struct {
	MachineID string
	IP        string
	CreatedAt time.Time
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func CreateStore() gokv.Store {

	options := redis.DefaultOptions // Address: "localhost:6379", Password: "", DB: 0

	// Create client
	store, err := redis.NewClient(options)
	if err != nil {
		fmt.Println("Error occured connecting to redis: ", err)
		panic(err)
	}
	return store
}

// func CreateStore() gokv.Store {
//   if !checkConnection() {
// 		panic("No connection to S3 could be established. Probably not running in a proper test environment.")
// 	}
// 	options := s3.DefaultOptions
// 	store, err := s3.NewStore(options)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return store
// }

// func CreateStore() gokv.Store {
// 	options := file.DefaultOptions
// 	store, err := file.NewStore(options)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return store
// }

// checkConnection returns true if a connection could be made, false otherwise.
// func checkConnection() bool {
// 	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
// 	os.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
// 	// No S3ForcePathStyle required because we're not accessing a specific bucket here.
// 	sess, err := session.NewSession(aws.NewConfig().WithRegion(endpoints.EuCentral1RegionID).WithEndpoint(customEndpoint))
// 	if err != nil {
// 		log.Printf("An error occurred during testing the connection to the server: %v\n", err)
// 		return false
// 	}
// 	svc := awss3.New(sess)
//
// 	_, err = svc.ListBuckets(&awss3.ListBucketsInput{})
// 	if err != nil {
// 		log.Printf("An error occurred during testing the connection to the server: %v\n", err)
// 		return false
// 	}
//
// 	return true
// }

// func CreateClient() s3.Client {
// 	// TODO
// 	minioClient, err := s3.New(s3.Config{
// 		Endpoint: os.Getenv("MINIO_ENDPOINT"),
// 		AccessKeyID: os.Getenv("MINIO_ACCESS_KEY"),
// 		SecretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
// 		UseSSL: false,
//     Codec: encoding.JSON,
// 	})
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	bucketName := "gokv"
// 	client, err := s3.NewClient(minioClient, bucketName, encoding.JSON)
//
// 	return client
// }

// func createClient() s3.Client {
// 	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
// 	os.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
// 	options := s3.Options{
// 		BucketName:             "gokv",
// 		Region:                 endpoints.EuCentral1RegionID,
// 		CustomEndpoint:         customEndpoint,
// 		UsePathStyleAddressing: true,
// 		Codec:                  encoding.JSON,
// 	}
// 	client, err := s3.NewClient(options)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return client
// }
