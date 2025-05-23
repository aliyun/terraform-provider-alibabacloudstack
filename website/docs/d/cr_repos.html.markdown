---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_repos"
sidebar_current: "docs-alibabacloudstack-datasource-cr-repos"
description: |-
  Provides a list of Container Registry repositories.
---

# alibabacloudstack_cr_repos

This data source provides a list Container Registry repositories on Alibabacloudstack Cloud.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_cr_repos" "my_repos" {
  name_regex  = "my-repos"
  output_file = "my-repo-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_repos.my_repos.repos}"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) Name of container registry namespace where the repositories are located in.
* `name_regex` - (Optional) A regex string to filter results by repository name.
* `enable_details` - (Optional) Boolean, false by default, only repository attributes are exported. Set to true if domain list and tags belong to this repository are needed. See `tags` in attributes.

* `ids` - (Optional) A list of matched Container Registry Repositories. Its element is set to `names`. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Repositories. Its element is set to `names`.
* `names` - A list of repository names.
* `repos` - A list of matched Container Registry Namespaces. Each element contains the following attributes:
  * `namespace` - Name of container registry namespace where repo is located.
  * `name` - Name of container registry namespace.
  * `summary` - The repository general information.
  * `repo_type` - `PUBLIC` or `PRIVATE`, repository's visibility.
  * `domain_list` - The repository domain list.
    * `public` - Domain of public endpoint.
    * `internal` - Domain of internal endpoint, only in some regions.
    * `vpc` - Domain of vpc endpoint.
  * `tags` - A list of image tags belong to this repository. Each contains several attributes, see `Block Tag`.

  * `summary` - The repository general information. 
  * `repo_type` - `PUBLIC` or `PRIVATE`, repository's visibility. 
  * `domain_list` - The repository domain list. 

### Block Tag

* `tag` - Tag of this image.
* `image_id` - Id of this image.
* `digest` - Digest of this image.
* `status` - Status of this image.
* `image_size` - Status of this image, in bytes.
* `image_update` - Last update time of this image, unix time in nanoseconds.
* `image_create` - Create time of this image, unix time in nanoseconds.