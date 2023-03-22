package cloudcontrol

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[A-Za-z0-9]{2,64}::[A-Za-z0-9]{2,64}::[A-Za-z0-9]{2,64}`), "must be three alphanumeric sections separated by double colons (::)"),
			},
			"type_version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CloudControlConn

	identifier := d.Get("identifier").(string)
	resourceDescription, err := FindResourceByID(ctx, conn,
		identifier,
		d.Get("type_name").(string),
		d.Get("type_version_id").(string),
		d.Get("role_arn").(string),
	)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Cloud Control API Resource (%s): %w", identifier, err))
	}

	d.SetId(aws.StringValue(resourceDescription.Identifier))

	d.Set("properties", resourceDescription.Properties)

	return nil
}
