package alibabacloudstack

import (

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type DataworksPublicService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DataworksPublicService) DescribeDataWorksFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":   parts[0],
		"ProjectId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFolder", "", nil, request)
	if err != nil {
		return object, err	}
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

func (s *DataworksPublicService) GetFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":   parts[0],
		"ProjectId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetFolder", "", nil, request)
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

func (s *DataworksPublicService) DescribeDataWorksConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId":  parts[1],
		"Name":       parts[2],
	}
	response, err = s.client.DoTeaRequest("GET", "dataworks-public", "2020-05-18", "ListConnections", "", nil, request)
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

func (s *DataworksPublicService) DescribeDataWorksUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "ListProjectMembers", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Data.ProjectMemberList", response)
	i := v.([]interface{})
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data.ProjectMemberList", response)
	}

	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DataworksPublicService) DescribeDataWorksUserRoleBinding(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "ListProjectRoles", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.ProjectRoleList", response)
	if v == nil || err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.ProjectRoleList", response)
	}
	i := v.([]interface{})

	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dataworks", id)), errmsgs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DataworksPublicService) DescribeDataWorksRemind(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"RemindId":   id,
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetRemind", "", nil, request)
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

func (s *DataworksPublicService) DescribeDataWorksProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ProjectId":  parts[1],
	}
	response, err = s.client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "GetProjectDetail", "", nil, request)
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
