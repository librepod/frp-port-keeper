package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/file"
)

type ProxyRecord struct {
	Port      int
	// IP        string
	CreatedAt time.Time
}

type PortRecord struct {
	Proxy string
	// IP        string
	CreatedAt time.Time
}

var DB gokv.Store

func init() {
	fmt.Println("Initializing store...")
	DB = createStore()
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// Redis Store
// func createStore() gokv.Store {
// 	options := redis.DefaultOptions // Address: "localhost:6379", Password: "", DB: 0
// 	db, err := redis.NewClient(options)
// 	if err != nil {
// 		fmt.Println("Error occured connecting to redis: ", err)
// 		panic(err)
// 	}
// 	return db
// }

// File Store
func createStore() gokv.Store {
	options := file.DefaultOptions
	store, err := file.NewStore(options)
	if err != nil {
		panic(err)
	}
	return store
}

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
