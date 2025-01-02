package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	yaml "gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCSKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCSKubernetesClustersRead,

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
						"pod_cidr": {
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

func dataSourceAlibabacloudStackCSKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClustersV1", "/api/v1/clusters")
	request.QueryParams["SignatureVersion"] = "1.0"
	request.QueryParams["ProductName"] = "CS"

	Cresponse := ClustersV1{}
	Clusterresponse := ClustersV1{}

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_kubernetes_clusters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	request.TransToAcsRequest()
	log.Printf("clusterResponse1 %v", response)
	log.Printf("clusterResponse2 %v", response.GetHttpContentBytes())
	err = json.Unmarshal(response.GetHttpContentBytes(), &Clusterresponse)
	if err != nil {
		return errmsgs.WrapError(err)
	}

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
			"id":                   kc.ClusterID,
			"name":                 kc.Name,
			"vpc_id":               kc.VpcID,
			"security_group_id":    kc.SecurityGroupID,
			"availability_zone":    kc.ZoneID,
			"state":                kc.State,
			"cluster_network_type": kc.NetworkMode,
			"pod_cidr":             kc.SubnetCidr,
		}

		ids = append(ids, kc.ClusterID)
		names = append(names, kc.Name)

		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		log.Printf("Entered kubeconfig")

		for i, k := range clusterids {
			log.Printf("IDS %v", clusterids)
			request := client.NewCommonRequest("POST", "Cs", "2015-12-15", "DescribeClusterUserKubeconfig", "")
			request.QueryParams["ClusterId"] = k
			request.QueryParams["PrivateIpAddress"] = "false"

			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			response, ok := raw.(*responses.CommonResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_kubernetes_clusters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			var conf KubeConfig
			var kubeconf Config
			err = json.Unmarshal(response.GetHttpContentBytes(), &conf)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			err = yaml.Unmarshal([]byte(conf.Config), &kubeconf)
			if err != nil {
				log.Fatalf("cannot unmarshal data: %v", err)
			}
			yamls, err := yaml.Marshal(kubeconf)
			filename := fmt.Sprint(file, i, ".yaml")
			err = ioutil.WriteFile(filename, yamls, 0777)
			if err != nil {
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
