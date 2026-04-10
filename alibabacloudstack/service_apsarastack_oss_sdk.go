package alibabacloudstack

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

// OssSdkService *connectivity.AlibabacloudStackClient
type OssSdkService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s OssSdkService) GetOssClient() (*oss.Client, error) {
	endpoint := s.client.Config.Endpoints[connectivity.OSSCode]
	schma := strings.ToLower(s.client.Config.Protocol)
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
	}

	// OSS SDK v2 uses a different configuration approach
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.client.Config.AccessKey,
			s.client.Config.SecretKey,
			s.client.Config.SecurityToken,
		)).
		WithEndpoint(endpoint).
		WithRegion(s.client.RegionId)

	client := oss.NewClient(cfg)
	return client, nil
}

func (s OssSdkService) GetBucketClient(bucketName string) (*oss.Client, error) {
	bucketInfo, err := s.DescribeOssBucket(bucketName)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if bucketInfo.Name == nil || *bucketInfo.Name == "" {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketName + " Not Found")
	}

	bucketEndpoint := bucketInfo.ExtranetEndpoint
	schma := strings.ToLower(s.client.Config.Protocol)
	if bucketEndpoint == nil || !strings.HasPrefix(*bucketEndpoint, "http") {
		endpoint := fmt.Sprintf("%s://%s.%s", schma, *bucketInfo.Name, *bucketEndpoint)
		bucketEndpoint = &endpoint
	}

	// OSS SDK v2 uses a different configuration approach
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.client.Config.AccessKey,
			s.client.Config.SecretKey,
			s.client.Config.SecurityToken,
		)).
		WithEndpoint(*bucketEndpoint).
		WithRegion(s.client.RegionId)

	client := oss.NewClient(cfg)

	return client, nil
}

func (s OssSdkService) DescribeOssBucket(bucketName string) (*oss.BucketInfo, error) {
	// Describe Oss BucketInfo for bucketName
	client, err := s.GetOssClient()
	if err != nil {
		return nil, err
	}
	var request *oss.GetBucketInfoRequest
	request.Bucket = &bucketName
	bucketResult, err := client.GetBucketInfo(context.Background(), request)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if bucketResult.StatusCode == 404 || &bucketResult.BucketInfo == nil {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketName + " Not Found")
	}
	return &bucketResult.BucketInfo, nil
}
