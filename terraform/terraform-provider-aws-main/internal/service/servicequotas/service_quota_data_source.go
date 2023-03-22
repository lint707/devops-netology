package servicequotas

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceServiceQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceQuotaRead,

		Schema: map[string]*schema.Schema{
			"adjustable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_value": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"global_quota": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quota_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"quota_code", "quota_name"},
			},
			"quota_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"quota_code", "quota_name"},
			},
			"service_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceQuotaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceQuotasConn

	quotaCode := d.Get("quota_code").(string)
	quotaName := d.Get("quota_name").(string)
	serviceCode := d.Get("service_code").(string)

	var err error
	var defaultQuota *servicequotas.ServiceQuota

	// A Service Quota will always have a default value, but will only have a current value if it has been set.
	// If it is not set, `GetServiceQuota` will return "NoSuchResourceException"
	if quotaName != "" {
		defaultQuota, err = findServiceQuotaDefaultByName(conn, serviceCode, quotaName)
		if err != nil {
			return fmt.Errorf("error getting Default Service Quota for (%s/%s): %w", serviceCode, quotaName, err)
		}

		quotaCode = aws.StringValue(defaultQuota.QuotaCode)
	} else {
		defaultQuota, err = findServiceQuotaDefaultByID(conn, serviceCode, quotaCode)
		if err != nil {
			return fmt.Errorf("error getting Default Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
		}
	}

	d.SetId(aws.StringValue(defaultQuota.QuotaArn))
	d.Set("adjustable", defaultQuota.Adjustable)
	d.Set("arn", defaultQuota.QuotaArn)
	d.Set("default_value", defaultQuota.Value)
	d.Set("global_quota", defaultQuota.GlobalQuota)
	d.Set("quota_code", defaultQuota.QuotaCode)
	d.Set("quota_name", defaultQuota.QuotaName)
	d.Set("service_code", defaultQuota.ServiceCode)
	d.Set("service_name", defaultQuota.ServiceName)
	d.Set("value", defaultQuota.Value)

	serviceQuota, err := findServiceQuotaByID(conn, serviceCode, quotaCode)
	if tfresource.NotFound(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error getting Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
	}

	d.Set("arn", serviceQuota.QuotaArn)
	d.Set("value", serviceQuota.Value)

	return nil
}
