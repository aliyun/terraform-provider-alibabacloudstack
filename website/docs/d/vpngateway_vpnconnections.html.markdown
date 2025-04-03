---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnconnections"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpnconnections"
description: |- 
  Provides a list of vpngateway vpnconnections owned by an alibabacloudstack account.
---

# alibabacloudstack_vpngateway_vpnconnections
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_connections`

This data source provides a list of vpngateway vpnconnections in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpngateway_vpnconnections" "example" {
  vpn_gateway_id      = "vgw-1234567890abcdef"
  customer_gateway_id = "cgw-abcdefgh12345678"
  name_regex         = "example-vpn-connection"
  output_file        = "vpn_connections_output.txt"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of IDs of the IPsec-VPN connections that you want to query. If this parameter is specified, only the specified resources will be returned.
* `vpn_gateway_id` - (Optional, ForceNew) The ID of the VPN gateway associated with the IPsec-VPN connections. This is a required parameter if you want to filter connections by the VPN gateway.
* `customer_gateway_id` - (Optional, ForceNew) The ID of the customer gateway associated with the IPsec-VPN connections. This is a required parameter if you want to filter connections by the customer gateway.
* `name_regex` - (Optional, ForceNew) A regular expression used to filter the results based on the name of the IPsec-VPN connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of IDs of the matching IPsec-VPN connections.
* `names` - A list of names of the matching IPsec-VPN connections.
* `connections` - A list of IPsec-VPN connection details. Each element contains the following attributes:
  * `id` - The unique identifier of the IPsec-VPN connection.
  * `customer_gateway_id` - The ID of the customer gateway associated with the IPsec-VPN connection.
  * `vpn_gateway_id` - The ID of the VPN gateway associated with the IPsec-VPN connection.
  * `name` - The name of the IPsec-VPN connection.
  * `local_subnet` - The CIDR block of the virtual private cloud (VPC) that the IPsec-VPN connection is attached to.
  * `remote_subnet` - The CIDR block of the on-premises data center or another network connected via the IPsec-VPN connection.
  * `create_time` - The time when the IPsec-VPN connection was created.
  * `effect_immediately` - Indicates whether IPsec-VPN negotiations are initiated immediately after the configuration is updated. Valid values: `true` or `false`.
  * `status` - The current status of the IPsec-VPN connection. Possible values include:
    * `ike_sa_not_established` - The IKE Security Association (SA) has not been established.
    * `ike_sa_established` - The IKE SA has been established.
    * `ipsec_sa_not_established` - The IPsec SA has not been established.
    * `ipsec_sa_established` - The IPsec SA has been established.
  * `ike_config` - The configurations of phase-one negotiation. It includes the following attributes:
    * `psk` - The pre-shared key used for authentication between the IPsec-VPN gateway and the customer gateway.
    * `ike_version` - The version of the IKE protocol. For example, `ikev1` or `ikev2`.
    * `ike_mode` - The negotiation mode of IKE phase-one. For example, `main` or `aggressive`.
    * `ike_enc_alg` - The encryption algorithm used during IKE phase-one negotiation. For example, `aes-128-cbc`, `aes-256-cbc`.
    * `ike_auth_alg` - The authentication algorithm used during IKE phase-one negotiation. For example, `sha1`, `sha256`.
    * `ike_pfs` - The Diffie-Hellman key exchange algorithm used by IKE phase-one negotiation. For example, `group2`, `group14`.
    * `ike_lifetime` - The lifetime of the Security Association (SA) as the result of IKE phase-one negotiation. Measured in seconds.
    * `ike_local_id` - The identification of the local (VPN gateway) endpoint.
    * `ike_remote_id` - The identification of the remote (customer gateway) endpoint.
  * `ipsec_config` - The configurations of phase-two negotiation. It includes the following attributes:
    * `ipsec_enc_alg` - The encryption algorithm used during IPsec phase-two negotiation. For example, `aes-128-cbc`, `aes-256-cbc`.
    * `ipsec_auth_alg` - The authentication algorithm used during IPsec phase-two negotiation. For example, `hmac-sha1-96`, `hmac-sha256-128`.
    * `ipsec_pfs` - The Diffie-Hellman key exchange algorithm used by IPsec phase-two negotiation. For example, `group2`, `group14`.
    * `ipsec_lifetime` - The lifetime of the Security Association (SA) as the result of IPsec phase-two negotiation. Measured in seconds.