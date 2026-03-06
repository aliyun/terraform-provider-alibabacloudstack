package alibabacloudstack

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DataworksService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DataworksService) DescribeDataWorksFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId": parts[0],
		"FolderId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFolder", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *DataworksService) GetFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":  parts[0],
		"ProjectId": parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFolder", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DataworksService) DescribeDataWorksConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId": parts[1],
		"Name":      parts[2],
	}
	response, err = s.client.DoTeaRequest("GET", "dataworks-public", "2020-05-18", "ListConnections", "", nil, nil, request)
	addDebug("ListConnections", response, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data.Connections", response)
	i := v.([]interface{})
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data.Connections", response)
	}
	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DataworksService) DescribeDataWorksUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId": parts[0],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "ListProjectMembers", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data.ProjectMemberList", response)
	i := v.([]interface{})
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data.ProjectMemberList", response)
	}

	for _, o := range i {
		object := o.(map[string]interface{})
		if object["ProjectMemberId"].(string) == convertAscmUid2MemberUid(parts[1]) {
			return object, nil
		}
	}
	return object, errmsgs.GetNotFoundErrorFromString("User " + id + " not existed.")

}

func (s *DataworksService) DescribeDataWorksUserRoleBinding(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	user, err := s.DescribeDataWorksUser(parts[1] + ":" + parts[2])
	if err != nil {
		return object, err
	}
	for _, i := range user["ProjectRoleList"].([]interface{}) {
		item := i.(map[string]interface{})
		if item["ProjectRoleCode"].(string) == parts[0] {
			return item, nil
		}
	}
	
	return object, errmsgs.GetNotFoundErrorFromString("no binding found")
}

func (s *DataworksService) DescribeDataWorksRemind(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"RemindId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetRemind", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DataworksService) DescribeDataWorksProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ProjectId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetProjectDetail", "", nil, nil, request)
	if err != nil {
		if e, ok := err.(*errmsgs.ComplexError); ok {
			err = e.Cause
		}
		if e, ok := err.(*tea.SDKError); ok {
			if strings.Contains(*e.Message, "does not exist.") {
				return object, errmsgs.GetNotFoundErrorFromString(*e.Message)
			}
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *DataworksService) ProjectStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksProject(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		var status string
		if v, err := toInt(object["Status"]); err == nil {
			status = strconv.Itoa(v)
		}
		if slices.Contains(failStates, status) {
			return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
		}

		return object, status, nil
	}
}

func (s *DataworksService) DescribeDataWorksBusiness(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	var response map[string]interface{}
	request := map[string]interface{}{
		"ProjectId":  parts[0],
		"BusinessId": parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetBusiness", "", nil, nil, request)
	if err != nil {
		if e, ok := err.(*errmsgs.ComplexError); ok {
			err = e.Cause
		}
		if e, ok := err.(*tea.SDKError); ok {
			if strings.Contains(*e.Message, "does not exist.") {
				return object, errmsgs.GetNotFoundErrorFromString(*e.Message)
			}
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})
	if _, existed := object["BusinessId"]; !existed {
		return object, errmsgs.GetNotFoundErrorFromString("Business " + id + " not existed.")
	}
	return object, nil
}

func (s *DataworksService) DescribeDataWorksFile(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId": parts[0],
		"FileId":    parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFile", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *DataworksService) GetFolderPath(id string) (path string, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":  parts[1],
		"ProjectId": parts[0],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFolder", "", nil, nil, request)
	if err != nil {
		return
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return
	}
	object := v.(map[string]interface{})
	parts = strings.SplitN(object["FolderPath"].(string), "/", 3)

	if len(parts) != 3 {
		return "", errmsgs.WrapError(fmt.Errorf("Invalid Folder Path %s.", object["FolderPath"].(string)))
	}
	parts[0] = "业务流程"
	parts[2] = strings.TrimPrefix(parts[2], "folder")
	return strings.Join(parts, "/"), nil
}


func (s *DataworksService) DescribeDataWorksBaseline(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, err
	}

	request := map[string]interface{}{
		"BaselineId": parts[1],
		"ProjectId":  parts[0],
	}


	response, err := s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetBaseline", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "The baseline does not exist.") {
			return object, errmsgs.GetNotFoundErrorFromString("Not Found Baseline "+ id)
		}
		return object, errmsgs.WrapError(err)
	}

	data, exists := response["Data"]
	if !exists {
		return object, errmsgs.GetNotFoundErrorFromString("Not Found Baseline "+ id)
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return object, errmsgs.GetNotFoundErrorFromString("Not Found Baseline "+ id)
	}
	
	return dataMap, nil
}