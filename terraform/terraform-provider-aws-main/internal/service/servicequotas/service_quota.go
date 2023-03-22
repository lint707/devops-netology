package servicequotas

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceServiceQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceQuotaCreate,
		Read:   resourceServiceQuotaRead,
		Update: resourceServiceQuotaUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"quota_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with alphabetic character"),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "must contain only alphanumeric and hyphen characters"),
				),
			},
			"quota_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with alphabetic character"),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "must contain only alphanumeric and hyphen characters"),
				),
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Required: true,
			},
		},
	}
}

func resourceServiceQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceQuotasConn

	quotaCode := d.Get("quota_code").(string)
	serviceCode := d.Get("service_code").(string)
	value := d.Get("value").(float64)

	d.SetId(fmt.Sprintf("%s/%s", serviceCode, quotaCode))

	// A Service Quota will always have a default value, but will only have a current value if it has been set.
	// If it is not set, `GetServiceQuota` will return "NoSuchResourceException"
	defaultQuota, err := findServiceQuotaDefaultByID(conn, serviceCode, quotaCode)
	if err != nil {
		return fmt.Errorf("error getting Default Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
	}
	quotaValue := aws.Float64Value(defaultQuota.Value)

	serviceQuota, err := findServiceQuotaByID(conn, serviceCode, quotaCode)
	if err != nil && !tfresource.NotFound(err) {
		return fmt.Errorf("error getting Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
	}
	if serviceQuota != nil {
		quotaValue = aws.Float64Value(serviceQuota.Value)
	}

	if value > quotaValue {
		input := &servicequotas.RequestServiceQuotaIncreaseInput{
			DesiredValue: aws.Float64(value),
			QuotaCode:    aws.String(quotaCode),
			ServiceCode:  aws.String(serviceCode),
		}

		output, err := conn.RequestServiceQuotaIncrease(input)

		if err != nil {
			return fmt.Errorf("error requesting Service Quota (%s) increase: %w", d.Id(), err)
		}

		if output == nil || output.RequestedQuota == nil {
			return fmt.Errorf("error requesting Service Quota (%s) increase: empty result", d.Id())
		}

		d.Set("request_id", output.RequestedQuota.Id)
	}

	return resourceServiceQuotaRead(d, meta)
}

func resourceServiceQuotaRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceQuotasConn

	serviceCode, quotaCode, err := resourceServiceQuotaParseID(d.Id())

	if err != nil {
		return err
	}

	// A Service Quota will always have a default value, but will only have a current value if it has been set.
	// If it is not set, `GetServiceQuota` will return "NoSuchResourceException"
	defaultQuota, err := findServiceQuotaDefaultByID(conn, serviceCode, quotaCode)
	if err != nil {
		return fmt.Errorf("error getting Default Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
	}

	d.Set("adjustable", defaultQuota.Adjustable)
	d.Set("arn", defaultQuota.QuotaArn)
	d.Set("default_value", defaultQuota.Value)
	d.Set("quota_code", defaultQuota.QuotaCode)
	d.Set("quota_name", defaultQuota.QuotaName)
	d.Set("service_code", defaultQuota.ServiceCode)
	d.Set("service_name", defaultQuota.ServiceName)
	d.Set("value", defaultQuota.Value)

	serviceQuota, err := findServiceQuotaByID(conn, serviceCode, quotaCode)
	if err != nil && !tfresource.NotFound(err) {
		return fmt.Errorf("error getting Service Quota for (%s/%s): %w", serviceCode, quotaCode, err)
	}

	if err == nil {
		d.Set("arn", serviceQuota.QuotaArn)
		d.Set("value", serviceQuota.Value)
	}

	requestID := d.Get("request_id").(string)

	if requestID != "" {
		input := &servicequotas.GetRequestedServiceQuotaChangeInput{
			RequestId: aws.String(requestID),
		}

		output, err := conn.GetRequestedServiceQuotaChange(input)

		if tfawserr.ErrCodeEquals(err, servicequotas.ErrCodeNoSuchResourceException) {
			d.Set("request_id", "")
			d.Set("request_status", "")
			return nil
		}

		if err != nil {
			return fmt.Errorf("error getting Service Quotas Requested Service Quota Change (%s): %w", requestID, err)
		}

		if output == nil || output.RequestedQuota == nil {
			return fmt.Errorf("error getting Service Quotas Requested Service Quota Change (%s): empty result", requestID)
		}

		requestStatus := aws.StringValue(output.RequestedQuota.Status)
		d.Set("request_status", requestStatus)

		switch requestStatus {
		case servicequotas.RequestStatusApproved, servicequotas.RequestStatusCaseClosed, servicequotas.RequestStatusDenied:
			d.Set("request_id", "")
		case servicequotas.RequestStatusCaseOpened, servicequotas.RequestStatusPending:
			d.Set("value", output.RequestedQuota.DesiredValue)
		}
	}

	return nil
}

func resourceServiceQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceQuotasConn

	value := d.Get("value").(float64)
	serviceCode, quotaCode, err := resourceServiceQuotaParseID(d.Id())

	if err != nil {
		return err
	}

	input := &servicequotas.RequestServiceQuotaIncreaseInput{
		DesiredValue: aws.Float64(value),
		QuotaCode:    aws.String(quotaCode),
		ServiceCode:  aws.String(serviceCode),
	}

	output, err := conn.RequestServiceQuotaIncrease(input)

	if err != nil {
		return fmt.Errorf("error requesting Service Quota (%s) increase: %w", d.Id(), err)
	}

	if output == nil || output.RequestedQuota == nil {
		return fmt.Errorf("error requesting Service Quota (%s) increase: empty result", d.Id())
	}

	d.Set("request_id", output.RequestedQuota.Id)

	return resourceServiceQuotaRead(d, meta)
}

func resourceServiceQuotaParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected SERVICE-CODE/QUOTA-CODE", id)
	}

	return parts[0], parts[1], nil
}
