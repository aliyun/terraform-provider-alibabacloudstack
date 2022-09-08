---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_share_permission"
sidebar_current: "docs-alibabacloudstack-resource-image-share-permission"
description: |-
  Provides an ECS image share permission resource.
---

# alibabacloudstack\_image\_share\_permission

Manage image sharing permissions. You can share your custom image to other AlibabacloudStack users. The user can use the shared custom image to create ECS instances or replace the system disk of the instance.

-> **NOTE:** You can only share your own custom images to other AlibabacloudStack users.

-> **NOTE:** Each custom image can be shared with up to 50 AlibabacloudStack accounts. You can submit a ticket to share with more users.

-> **NOTE:** After creating an ECS instance using a shared image, once the custom image owner releases the image sharing relationship or deletes the custom image, the instance cannot initialize the system disk.

## Example Usage

```
resource "alibabacloudstack_image_share_permission" "default" {
  image_id           = "m-bp1gxyh***"
  account_id         = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, ForceNew) The source image ID.
* `account_id` - (Required, ForceNew) Alibabacloudstack Account ID. It is used to share images.
   
   

### Attributes Reference0
 
 The following attributes are exported:
 
* `id` - ID of the image. It formats as `<image_id>:<account_id>`
