package alibabacloudstack

import (
	"github.com/denverdino/aliyungo/ecs"
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JsonTag struct {
	TagKey   string
	TagValue string
}

type AddTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     string
	Tag          []Tag
}

type RemoveTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     string
	Tag          []Tag
}
