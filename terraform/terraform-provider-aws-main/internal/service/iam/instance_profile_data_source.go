package iam

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstanceProfileRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInstanceProfileRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	name := d.Get("name").(string)

	req := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	log.Printf("[DEBUG] Reading IAM Instance Profile: %s", req)
	resp, err := conn.GetInstanceProfile(req)
	if err != nil {
		return fmt.Errorf("Error getting instance profiles: %w", err)
	}
	if resp == nil {
		return fmt.Errorf("no IAM instance profile found")
	}

	instanceProfile := resp.InstanceProfile

	d.SetId(aws.StringValue(instanceProfile.InstanceProfileId))
	d.Set("arn", instanceProfile.Arn)
	d.Set("create_date", fmt.Sprintf("%v", instanceProfile.CreateDate))
	d.Set("path", instanceProfile.Path)

	if len(instanceProfile.Roles) > 0 {
		role := instanceProfile.Roles[0]
		d.Set("role_arn", role.Arn)
		d.Set("role_id", role.RoleId)
		d.Set("role_name", role.RoleName)
	}

	return nil
}
