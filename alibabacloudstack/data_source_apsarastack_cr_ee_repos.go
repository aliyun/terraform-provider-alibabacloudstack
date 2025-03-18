package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCrEeRepos() *schema.Resource {
	return &schema.Resource{
		Read:	dataSourceAlibabacloudStackCrEeReposRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:		schema.TypeString,
				ForceNew:	true,
				Required:	true,
			},
			"namespace": {
				Type:		schema.TypeString,
				Optional:	true,
				ValidateFunc:	validation.StringIsValidRegExp,
			},
			"name_regex": {
				Type:		schema.TypeString,
				Optional:	true,
				ValidateFunc:	validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:		schema.TypeString,
				Optional:	true,
			},
			"enable_details": {
				Type:		schema.TypeBool,
				Optional:	true,
				Default:	false,
			},

			// Computed values
			"ids": {
				Type:		schema.TypeList,
				Optional:	true,
				Computed:	true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"repos": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"namespace": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"id": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"name": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"summary": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"repo_type": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"tags": {
							Type:		schema.TypeList,
							Computed:	true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag": {
										Type:		schema.TypeString,
										Computed:	true,
									},
									"image_id": {
										Type:		schema.TypeString,
										Computed:	true,
									},
									"digest": {
										Type:		schema.TypeString,
										Computed:	true,
									},
									"status": {
										Type:		schema.TypeString,
										Computed:	true,
									},
									"image_size": {
										Type:		schema.TypeInt,
										Computed:	true,
									},
									"image_update": {
										Type:		schema.TypeString,
										Computed:	true,
									},
									"image_create": {
										Type:		schema.TypeString,
										Computed:	true,
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

func dataSourceAlibabacloudStackCrEeReposRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	pageNo := 1
	pageSize := 100
	instanceId := d.Get("instance_id").(string)

	var namespaces []string
	if namespace, ok := d.GetOk("namespace"); ok {
		namespaces = append(namespaces, namespace.(string))
	} else {
		for {
			resp, err := crService.ListCrEeNamespaces(instanceId, pageNo, pageSize)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			for _, n := range resp.Namespaces {
				namespaces = append(namespaces, n.NamespaceName)
			}
			if len(resp.Namespaces) < pageSize {
				break
			}
			pageNo++
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var idsMap map[string]string
	if v, ok := d.GetOk("ids"); ok {
		idsMap = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var enableDetails bool
	if v, ok := d.GetOk("enable_details"); ok {
		enableDetails = v.(bool)
	}

	var repos []cr_ee.RepositoriesItem
	for _, namespace := range namespaces {
		pageNo = 1
		for {
			resp, err := crService.ListCrEeRepos(instanceId, namespace, pageNo, pageSize)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			for _, r := range resp.Repositories {
				if nameRegex != nil && !nameRegex.MatchString(r.RepoName) {
					continue
				}
				if idsMap != nil && idsMap[r.RepoId] == "" {
					continue
				}

				repos = append(repos, r)
			}

			if len(resp.Repositories) < pageSize {
				break
			}
			pageNo++
		}
	}

	tags := make([][]map[string]interface{}, len(repos))
	if enableDetails {
		for i, repo := range repos {
			pageNo = 1
			var images []cr_ee.ImagesItem
			for {
				resp, err := crService.ListCrEeRepoTags(instanceId, repo.RepoId, pageNo, pageSize)
				if err != nil {
					return errmsgs.WrapError(err)
				}
				images = append(images, resp.Images...)

				if len(resp.Images) < pageSize {
					break
				}
				pageNo++
			}

			repoTags := make([]map[string]interface{}, len(images))
			for j, image := range images {
				m := make(map[string]interface{})
				m["tag"] = image.Tag
				m["image_id"] = image.ImageId
				m["digest"] = image.Digest
				m["status"] = image.Status
				m["image_size"] = image.ImageSize
				m["image_update"] = image.ImageUpdate
				m["image_create"] = image.ImageCreate
				repoTags[j] = m
			}
			tags[i] = repoTags
		}
	}

	ids := make([]string, len(repos))
	names := make([]string, len(repos))
	reposMaps := make([]map[string]interface{}, len(repos))
	for i, r := range repos {
		ids[i] = r.RepoId
		names[i] = r.RepoName
		m := make(map[string]interface{})
		m["instance_id"] = r.InstanceId
		m["namespace"] = r.RepoNamespaceName
		m["id"] = r.RepoId
		m["name"] = r.RepoName
		m["summary"] = r.Summary
		m["repo_type"] = r.RepoType
		if enableDetails {
			m["tags"] = tags[i]
		}
		reposMaps[i] = m
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("repos", reposMaps); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			if err := writeToFile(output.(string), reposMaps); err != nil {
			return err
		}
	}

	return nil
}
