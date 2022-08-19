package apsarastack

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

type Tag struct {
	Key   string
	Value string
}

type JsonTag struct {
	TagKey   string
	TagValue string
}

type AddTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     common.Region
	Tag          []Tag
}

type RemoveTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     common.Region
	Tag          []Tag
}
