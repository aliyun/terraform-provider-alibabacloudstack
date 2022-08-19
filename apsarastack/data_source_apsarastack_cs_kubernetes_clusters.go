package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func dataSourceApsaraStackCSKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackCSKubernetesClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_internet_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"master_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_numbers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_cidr_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_data_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_data_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"master_auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"worker_auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"worker_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connections": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_server_internet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"api_server_intranet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"master_public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackCSKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "Cs"
	request.Version = "2015-12-15"

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	request.ServiceCode = "cs"
	request.ApiName = "DescribeClustersV1"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "Cs", "RegionId": client.RegionId, "Action": "DescribeClustersV1", "Version": cs.CSAPIVersion, "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	Cresponse := ClustersV1{}
	Clusterresponse := ClustersV1{}

	raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
		return csClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cs_kubernetes_clusters", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	resp, _ := raw.(*responses.CommonResponse)
	request.TransToAcsRequest()
	log.Printf("clusterResponse1 %v", resp)
	log.Printf("clusterResponse2 %v", resp.GetHttpContentBytes())
	err = json.Unmarshal(resp.GetHttpContentBytes(), &Clusterresponse)
	if err != nil {
		return WrapError(err)
	}
	//var nullc ClustersV1
	//if Clusterresponsenullc{
	//	return WrapErrorf(err,"Response is nil")
	//}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	var clusterids map[string]string
	if v, ok := d.GetOk("ids"); ok {
		clusterids = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			clusterids[vv.(string)] = vv.(string)
		}
	}

	log.Printf("Entering kubeconfig %v 52", clusterids)
	for _, cresp := range Clusterresponse.Clusters {
		if clusterids != nil && clusterids[cresp.ClusterID] == "" {
			continue
		}
		if r != nil && !r.MatchString(cresp.Name) {
			continue
		}
		Cresponse.Clusters = append(Cresponse.Clusters, cresp)
	}
	Clusterresponse = Cresponse
	log.Printf("Clusterresponse idfiltered %v", Clusterresponse)
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, kc := range Clusterresponse.Clusters {
		if r != nil && !r.MatchString(kc.Name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                kc.ClusterID,
			"name":              kc.Name,
			"vpc_id":            kc.VpcID,
			"security_group_id": kc.SecurityGroupID,
			"availability_zone": kc.ZoneID,
			"state":             kc.State,
			//"master_instance_types":       []string{kc.Parameters.MasterInstanceType},
			//"nat_gateway_id":              kc.Parameters.NatGatewayID,
			"vswitch_ids": []string{kc.VswitchID},
			//"master_disk_category":        kc.MasterSystemDiskCategory,
			"cluster_network_type": kc.NetworkMode,
			"pod_cidr":             kc.SubnetCidr,
			//"worker_data_disk_size":       kc.Parameters.WorkerDataDiskSize,
			//"worker_disk_category":        kc.Parameters.WorkerDataDiskCategory,
			//"worker_instance_types":       []string{kc.Parameters.WorkerInstanceType},
			//"worker_instance_charge_type": kc.Parameters.WorkerInstanceChargeType,
			//"node_cidr_mask":              kc.MetaData.Capabilities.NodeCIDRMask,
		}

		ids = append(ids, kc.ClusterID)
		names = append(names, kc.Name)

		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		log.Printf("Entered kubeconfig")

		for i, k := range clusterids {
			log.Printf("IDS %v", clusterids)
			request.Method = "POST"
			request.Product = "Cs"
			request.Version = "2015-12-15"

			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.ServiceCode = "cs"
			request.ApiName = "DescribeClusterUserKubeconfig"
			request.Headers = map[string]string{"RegionId": client.RegionId, "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			request.RegionId = client.RegionId
			request.QueryParams = map[string]string{
				"AccessKeyId":      client.AccessKey,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Cs",
				"RegionId":         client.RegionId,
				"Action":           "DescribeClustersV1",
				"Version":          cs.CSAPIVersion,
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"PrivateIpAddress": "false",
				"ClusterId":        k,
			}
			log.Printf("request body %v", request)
			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cs_kubernetes_clusters", request.GetActionName(), ApsaraStackSdkGoERROR)
			}
			resp, _ := raw.(*responses.CommonResponse)
			//request.TransToAcsRequest()
			var conf KubeConfig
			var kubeconf Config
			err = json.Unmarshal(resp.GetHttpContentBytes(), &conf)
			if err != nil {
				return WrapError(err)
			}
			err = yaml.Unmarshal([]byte(conf.Config), &kubeconf)
			if err != nil {
				log.Fatalf("cannot unmarshal data: %v", err)
			}
			yamls, err := yaml.Marshal(kubeconf)
			filename := fmt.Sprint(file, i, ".yaml")
			err = ioutil.WriteFile(filename, yamls, 0777)
			//// handle this error
			if err != nil {
				// print it out
				fmt.Println(err)
			}

			log.Printf("kubeconfig check %v ", conf.Config)
		}

	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

type KubeConfig struct {
	ServerRole      string `json:"serverRole"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	AsapiRequestID  string `json:"asapiRequestId"`
	Domain          string `json:"domain"`
	API             string `json:"api"`
	Config          string `json:"config"`
}

type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}
