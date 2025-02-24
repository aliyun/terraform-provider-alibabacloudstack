package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KmsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *KmsService) DescribeKmsKey(id string) (object kms.KeyMetadata, err error) {
	request := kms.CreateDescribeKeyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.KeyId = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	bresponse, ok := raw.(*kms.DescribeKeyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.AliasNotFound", "Forbidden.KeyNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KmsKey", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.KeyMetadata.KeyState == "PendingDeletion" {
		log.Printf("[WARN] Removing KmsKey  %s because it's already gone", id)
		return bresponse.KeyMetadata, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KmsKey", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse.KeyMetadata, nil
}

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (*kms.DecryptResponse, error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	request := kms.CreateDecryptRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.CiphertextBlob = ciphertextBlob
	request.EncryptionContext = string(context[:])

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Decrypt(request)
	})
	bresponse, ok := raw.(*kms.DecryptResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, context, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return bresponse, err
}

func (s *KmsService) Decrypt2(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, errmsgs.WrapError(err)
	}

	request := map[string]interface{}{
		"CiphertextBlob":   ciphertextBlob,
		"EncryptionContext": string(context[:]),
		"Product":          "Kms",
		"OrganizationId":   s.client.Department,
	}

	var response map[string]interface{}
	response, err = s.client.DoTeaRequest("POST", "Kms", "2016-01-20", "Decrypt", "", nil, nil, request)
	if err != nil {
		return plaintext, err
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

func (s *KmsService) DescribeKmsSecret(id string) (object kms.DescribeSecretResponse, err error) {
	request := kms.CreateDescribeSecretRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SecretName = id
	request.FetchTags = "true"

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeSecret(request)
	})
	bresponse, ok := raw.(*kms.DescribeSecretResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.errmsgs.ResourceNotfound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KmsSecret", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *bresponse, nil
}

func (s *KmsService) GetSecretValue(id string) (object kms.GetSecretValueResponse, err error) {
	request := kms.CreateGetSecretValueRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.SecretName = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.GetSecretValue(request)
	})
	bresponse, ok := raw.(*kms.GetSecretValueResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.errmsgs.ResourceNotfound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("kmssecret", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *bresponse, nil
}

func (s *KmsService) setResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]JsonTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, JsonTag{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := kms.CreateUntagResourceRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		if resourceType == "key" {
			request.KeyId = d.Id()
		}
		if resourceType == "secret" {
			request.SecretName = d.Id()
		}
		remove, err := json.Marshal(removed)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.TagKeys = string(remove)
		raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UntagResource(request)
		})
		addDebug(request.GetActionName(), raw)
		bresponse, ok := raw.(*kms.UntagResourceResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	if len(added) > 0 {
		request := kms.CreateTagResourceRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		if resourceType == "key" {
			request.KeyId = d.Id()
		}
		if resourceType == "secret" {
			request.SecretName = d.Id()
		}
		add, err := json.Marshal(added)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.Tags = string(add)
		raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.TagResource(request)
		})
		addDebug(request.GetActionName(), raw)
		bresponse, ok := raw.(*kms.TagResourceResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func (s *KmsService) DescribeKmsAlias(id string) (object kms.KeyMetadata, err error) {
	request := kms.CreateDescribeKeyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.KeyId = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	bresponse, ok := raw.(*kms.DescribeKeyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.AliasNotFound", "Forbidden.KeyNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KmsAlias", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse.KeyMetadata, nil
}

func (s *KmsService) DescribeKmsKeyVersion(id string) (object kms.DescribeKeyVersionResponse, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := kms.CreateDescribeKeyVersionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.KeyId = parts[0]
	request.KeyVersionId = parts[1]

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKeyVersion(request)
	})
	bresponse, ok := raw.(*kms.DescribeKeyVersionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *bresponse, nil
}
