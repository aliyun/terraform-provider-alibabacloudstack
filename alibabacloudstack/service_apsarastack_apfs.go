package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type EfsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EfsService) DescribeFileSystem(id string) (response map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["FileSystemId"] = id
	response, err = s.client.DoTeaRequest("POST", "EFS", "2017-06-26", "DescribeFileSystems", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	fileSystems, err := jsonpath.Get("$.FileSystems.FileSystem", response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	for _, item := range fileSystems.([]interface{}) {
		fileSystem := item.(map[string]interface{})
		if fileSystem["FileSystemId"].(string) == id {
			return fileSystem, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Apfs File System " + id)
}

func (s *EfsService) ApfsFileSystemStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeFileSystem(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		for _, failState := range failStates {
			status := object["Status"].(string)
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}
		return object, object["Status"].(string), nil
	}
}
