---
subcategory: "云服务总线 CSB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_service"
sidebar_current: "docs-Alibabacloudstack-resource-csb-service"
description: |-
  Provides a Alibabacloudstack resource to manage CSB Service.
---

# alibabacloudstack_csb_service

This resource will help you to manage CSB Service.

For information about CSB Service and how to use it, see [Create a Service](https://help.aliyun.com/apsara/enterprise/v_3_18_0_30393230/csb/apsarastack-developer-guide/createservice.html)

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_csb_service" "service" {
  csb_id          = "your-csb-id"
  project_id      = "your-project-id"
  service_name    = "example-service"
  service_version = "1.0.0"
  provide_type    = "RESTful"
  skip_auth       = false
  all_visiable    = true
  route_conf_json = jsonencode({
    importConf = {
      accessEndpointJSON = jsonencode({
        # endpoint configuration
      })
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `csb_id` - (Required, ForceNew, Int) The unique ID of the CSB instance. Modifying this parameter will force the creation of a new resource.
* `project_id` - (Required, String) The id of the project to which the service belongs.
* `service_name` - (Required, ForceNew, String) The name of the service. Modifying this parameter will force the creation of a new resource.
* `service_version` - (Required, ForceNew, String) The version of the service (e.g., '1.0.0'). Modifying this parameter will force the creation of a new resource.
* `route_conf_json` - (Required, String) The route configuration in JSON format. This field contains the import configuration and access endpoint settings.
* `alias` - (Optional, String) The alias name of the service.
* `model_version` - (Optional, String) The model version of the service. Default: `2.0`.
* `skip_auth` - (Optional, Bool) Whether to skip authentication for this service. Default: `false`.
* `all_visiable` - (Optional, Bool) Whether the service is visible to all users. Default: `true`.
* `scope` - (Optional, String) The scope of the service (e.g., '0' for private, '1' for public). Default: `0`.
* `consume_types` - (Optional, List of String) List of consume types. Valid values: `Restful`, `WebService`.
* `provide_type` - (Optional, String) The provide type of the service. Valid values: `RESTful`, `SpringCloud`, `HSF`, `WebService`, `DUBBO`, `JDBC`. Default: `Restful`.
* `cas_serv_targets` - (Optional, List of String) List of CAS service targets.
* `qps` - (Optional, Int) Queries per second (read-only monitoring metric). Default: `0`.
* `interface_name` - (Optional, String) The interface name (usually empty for RESTful services).
* `ip_white_str` - (Optional, String) IP whitelist in string format.
* `ip_black_str` - (Optional, String) IP blacklist in string format.
* `err_def_json` - (Optional, String) Error definition in JSON format.
* `access_params_json` - (Optional, String) Access parameters in JSON format.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the resource, in the format `csb_id:service_id`.
* `service_id` - The internal ID of the CSB service.
* `csb_id` - The unique ID of the CSB instance.
* `project_id` - The id of the project to which the service belongs.
* `service_name` - The name of the service.
* `service_version` - The version of the service.
* `alias` - The alias name of the service.
* `model_version` - The model version of the service.
* `skip_auth` - Whether to skip authentication for this service.
* `all_visiable` - Whether the service is visible to all users.
* `scope` - The scope of the service.
* `provide_type` - The provide type of the service.
* `consume_types` - List of consume types.
* `cas_serv_targets` - List of CAS service targets.
* `qps` - Queries per second.
* `status` - Service status (e.g., 0=inactive, 1=active).
* `modified_time` - Last modified timestamp in milliseconds.
* `principal_name` - The principal name associated with the service.
* `owner_id` - The owner ID of the service.
* `user_id` - The user ID of the service creator.
* `interface_name` - The interface name.
* `ssl` - Whether SSL is enabled for the service.
* `ott_flag` - OTF flag status.
* `policy_handler` - The policy handler (e.g., 'accept').
* `ip_white_str` - IP whitelist string.
* `ip_black_str` - IP blacklist string.
* `err_def_json` - Error definition JSON.
* `access_params_json` - Access parameters JSON.
* `route_conf_json` - Route configuration JSON.

## Import

CSB Service can be imported using the combination of `csb_id` and `service_id`, separated by a colon, e.g.

```
$ terraform import alibabacloudstack_csb_service.example <csb_id>:<service_id>
```
