func buildBasicPrometheusInstance(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = "%s"
  tags         = ["test1", "test2"]
}

data "alibabacloudstack_prometheus_v2_instances" "default" {
  name_regex = alibabacloudstack_prometheus_v2_instance.default.cluster_name
}
`, name)
}

func buildBasicNotifyGroup(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  name        = "%s_notify_group"
  type        = "WEBHOOK"
  description = "%s_notify_group_description"
  webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=56b42bc6e7cad53bab514a583847db73c68fa1804b0e72af7167954b66f7aea8"
  webhook_header_params {
    key   = "Content-Type"
    value = "application/json"
  }
}

data "alibabacloudstack_prometheus_v2_notify_groups" "default" {
  name_regex = alibabacloudstack_prometheus_v2_notify_group.default.name
}
`, name, name)
}

func TestAccAlibabacloudStackPrometheusV2Alert_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_alert.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"name":                 "tfacc_alert",
		"alert_type":           "PROMETHEUS",
		"notify_recovered":     "true",
		"tag_set.#":            "2",
		"tag_set.0":            "aaa",
		"tag_set.1":            "ccc",
		"recover_notification": "��ض���\\${alert_source} \\n�ָ�ʱ�䣺\\${alert_time}",
	})
	serviceFunc := func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tfacc_alert"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		instanceConfig := buildBasicPrometheusInstance(name)
		groupConfig := buildBasicNotifyGroup(name)

		return fmt.Sprintf(`
%s
%s

resource "alibabacloudstack_prometheus_v2_alert" "default" {
  name               = "%s"
  alert_type         = "PROMETHEUS"
  notify_recovered   = true
  is_check_all       = false
  tag_set            = ["aaa", "ccc"]
  recover_notification = "��ض���\\${alert_source} \\n�ָ�ʱ�䣺\\${alert_time}"

  trigger_rule = {
    clusterIds = [data.alibabacloudstack_prometheus_v2_instances.default.instances.0.id]
    promql     = "select testfield from testtable where testfield >= 0"
    period     = "5m"
    severity   = "warning"
    cron       = "0 /5 * * * ?"
    timeType   = 1
  }

  notification = {
    message = "����������$${condition}\\n���м�¼��$${alert_result}"
  }

  alert_notify_params {
    notify_types      = ["EMAIL", "SMS"]
    notify_group_ids  = [data.alibabacloudstack_prometheus_v2_notify_groups.default.groups.0.id]
    notify_interval   = "10m"
  }
}
`, instanceConfig, groupConfig, name)
	})

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             "tfacc_alert_updated",
					"notify_recovered": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             "tfacc_alert_updated",
						"notify_recovered": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
