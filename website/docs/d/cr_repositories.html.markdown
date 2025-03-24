---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_repos"
sidebar_current: "docs-alibabacloudstack-datasource-cr-ee-repos"
description: |-
  Provides a list of Container Registry Enterprise Edition repositories.
---

# alibabacloudstack_cr_ee_repos
-> **NOTE:** Alias name has: `alibabacloudstack_cr_repositories`

This data source provides a list Container Registry Enterprise Edition repositories on Alibaba Cloud.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_cr_ee_repos" "my_repos" {
  instance_id = "cri-xx"
  name_regex  = "my-repos"
  output_file = "my-repo-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_ee_repos.my_repos.repos}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of Container Registry Enterprise Edition instance.
* `namespace` - (Optional) Name of Container Registry Enterprise Edition namespace where the repositories are located in.
* `ids` - (Optional) A list of ids to filter results by repository id.
* `name_regex` - (Optional) A regex string to filter results by repository name.
* `enable_details` - (Optional) Boolean, false by default, only repository attributes are exported. Set to true if tags belong to this repository are needed. See `tags` in attributes.
* `names` - (Optional) A list of repository names.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Enterprise Edition repositories. Its element is a repository id.
* `names` - A list of repository names.
* `repos` - A list of matched Container Registry Enterprise Edition namespaces. Each element contains the following attributes:
  * `instance_id` - ID of Container Registry Enterprise Edition instance.
  * `namespace` - Name of Container Registry Enterprise Edition namespace where repo is located.
  * `id` - ID of Container Registry Enterprise Edition repository.
  * `name` - Name of Container Registry Enterprise Edition repository.
  * `summary` - The repository general information.
  * `repo_type` - `PUBLIC` or `PRIVATE`, repository's visibility.
  * `tags` - A list of image tags belong to this repository. Each contains several attributes, see `Block Tag`.
    * `tag` - Tag of this image.
    * `image_id` - Id of this image.
    * `digest` - Digest of this image.
    * `status` - Status of this image.
    * `image_size` - Status of this image, in bytes.
    * `image_update` - Last update time of this image, unix time in nanoseconds.
    * `image_create` - Create time of this image, unix time in nanoseconds.
* `names` - A list of repository names.