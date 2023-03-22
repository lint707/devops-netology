package guardduty

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceDetector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDetectorRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finding_publishing_frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDetectorRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GuardDutyConn

	detectorId := d.Get("id").(string)

	if detectorId == "" {
		input := &guardduty.ListDetectorsInput{}

		resp, err := conn.ListDetectors(input)
		if err != nil {
			return fmt.Errorf("error listing GuardDuty Detectors: %w", err)
		}

		if resp == nil || len(resp.DetectorIds) == 0 {
			return fmt.Errorf("no GuardDuty Detectors found")
		}
		if len(resp.DetectorIds) > 1 {
			return fmt.Errorf("multiple GuardDuty Detectors found; please use the `id` argument to look up a single detector")
		}

		detectorId = aws.StringValue(resp.DetectorIds[0])
	}

	getInput := &guardduty.GetDetectorInput{
		DetectorId: aws.String(detectorId),
	}

	getResp, err := conn.GetDetector(getInput)
	if err != nil {
		return err
	}

	if getResp == nil {
		return fmt.Errorf("cannot receive GuardDuty Detector details")
	}

	d.SetId(detectorId)
	d.Set("status", getResp.Status)
	d.Set("service_role_arn", getResp.ServiceRole)
	d.Set("finding_publishing_frequency", getResp.FindingPublishingFrequency)

	return nil
}
