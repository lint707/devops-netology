package servicequotas_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

// This resource is different than many since quotas are pre-existing
// and the resource is only designed to help with increases.
// In the basic case, we test that the resource can match the existing quota
// without unexpected changes.
func TestAccServiceQuotasServiceQuota_basic(t *testing.T) {
	const dataSourceName = "data.aws_servicequotas_service_quota.test"
	const resourceName = "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheck(t)
			preCheckServiceQuotaSet(setQuotaServiceCode, setQuotaQuotaCode, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, servicequotas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceQuotaConfig_sameValue(setQuotaServiceCode, setQuotaQuotaCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "adjustable", dataSourceName, "adjustable"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "default_value", dataSourceName, "default_value"),
					resource.TestCheckResourceAttrPair(resourceName, "quota_code", dataSourceName, "quota_code"),
					resource.TestCheckResourceAttrPair(resourceName, "quota_name", dataSourceName, "quota_name"),
					resource.TestCheckResourceAttrPair(resourceName, "service_code", dataSourceName, "service_code"),
					resource.TestCheckResourceAttrPair(resourceName, "service_name", dataSourceName, "service_name"),
					resource.TestCheckResourceAttrPair(resourceName, "value", dataSourceName, "value"),
					resource.TestCheckNoResourceAttr(resourceName, "request_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceQuotasServiceQuota_basic_Unset(t *testing.T) {
	const dataSourceName = "data.aws_servicequotas_service_quota.test"
	const resourceName = "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheck(t)
			preCheckServiceQuotaUnset(unsetQuotaServiceCode, unsetQuotaQuotaCode, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, servicequotas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceQuotaConfig_sameValue(unsetQuotaServiceCode, unsetQuotaQuotaCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "adjustable", dataSourceName, "adjustable"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "default_value", dataSourceName, "default_value"),
					resource.TestCheckResourceAttrPair(resourceName, "quota_code", dataSourceName, "quota_code"),
					resource.TestCheckResourceAttrPair(resourceName, "quota_name", dataSourceName, "quota_name"),
					resource.TestCheckResourceAttrPair(resourceName, "service_code", dataSourceName, "service_code"),
					resource.TestCheckResourceAttrPair(resourceName, "service_name", dataSourceName, "service_name"),
					resource.TestCheckResourceAttrPair(resourceName, "value", dataSourceName, "value"),
					resource.TestCheckNoResourceAttr(resourceName, "request_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceQuotasServiceQuota_Value_increaseOnCreate(t *testing.T) {
	quotaCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_QUOTA_CODE")
	if quotaCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_QUOTA_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	serviceCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_SERVICE_CODE")
	if serviceCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_SERVICE_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	value := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_VALUE")
	if value == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_VALUE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	resourceName := "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicequotas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceQuotaConfig_value(serviceCode, quotaCode, value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
				),
			},
		},
	})
}

func TestAccServiceQuotasServiceQuota_Value_increaseOnUpdate(t *testing.T) {
	quotaCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_QUOTA_CODE")
	if quotaCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_QUOTA_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	serviceCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_SERVICE_CODE")
	if serviceCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_SERVICE_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	value := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_VALUE")
	if value == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_VALUE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	dataSourceName := "aws_servicequotas_service_quota.test"
	resourceName := "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicequotas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceQuotaConfig_sameValue(serviceCode, quotaCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttrPair(resourceName, "value", dataSourceName, "value"),
					resource.TestCheckNoResourceAttr(resourceName, "request_id"),
				),
			},
			{
				Config: testAccServiceQuotaConfig_value(serviceCode, quotaCode, value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
				),
			},
		},
	})
}

func TestAccServiceQuotasServiceQuota_permissionError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t); acctest.PreCheckAssumeRoleARN(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicequotas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccServiceQuotaConfig_permissionError("elasticloadbalancing", "L-53DA6B97"),
				ExpectError: regexp.MustCompile(`DEPENDENCY_ACCESS_DENIED_ERROR`),
			},
		},
	})
}

// nosemgrep:ci.servicequotas-in-func-name
func testAccServiceQuotaConfig_sameValue(serviceCode, quotaCode string) string {
	return fmt.Sprintf(`
data "aws_servicequotas_service_quota" "test" {
  quota_code   = %[1]q
  service_code = %[2]q
}

resource "aws_servicequotas_service_quota" "test" {
  quota_code   = data.aws_servicequotas_service_quota.test.quota_code
  service_code = data.aws_servicequotas_service_quota.test.service_code
  value        = data.aws_servicequotas_service_quota.test.value
}
`, quotaCode, serviceCode)
}

func testAccServiceQuotaConfig_value(serviceCode, quotaCode, value string) string {
	return fmt.Sprintf(`
resource "aws_servicequotas_service_quota" "test" {
  quota_code   = %[1]q
  service_code = %[2]q
  value        = %[3]s
}
`, quotaCode, serviceCode, value)
}

func testAccServiceQuotaConfig_permissionError(serviceCode, quotaCode string) string {
	policy := `{
  "Version": "2012-10-17",
  "Statement": [
    {
  	  "Effect": "Allow",
  	  "Action": [
  	    "servicequotas:GetServiceQuota"
  	  ],
  	  "Resource": "*"
    },
    {
  	  "Effect": "Deny",
  	  "Action": [
  	    "elasticloadbalancing:*"
  	  ],
  	  "Resource": "*"
    }
  ]
}`

	return acctest.ConfigCompose(
		acctest.ConfigAssumeRolePolicy(policy),
		fmt.Sprintf(`
resource "aws_servicequotas_service_quota" "test" {
  service_code = %[1]q
  quota_code   = %[2]q
  value        = 1
}
`, serviceCode, quotaCode))
}
