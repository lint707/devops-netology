package imagebuilder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceComponent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComponentRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"change_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_os_versions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tftags.TagsSchemaComputed(),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComponentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &imagebuilder.GetComponentInput{}

	if v, ok := d.GetOk("arn"); ok {
		input.ComponentBuildVersionArn = aws.String(v.(string))
	}

	output, err := conn.GetComponent(input)

	if err != nil {
		return fmt.Errorf("error getting Image Builder Component: %w", err)
	}

	if output == nil || output.Component == nil {
		return fmt.Errorf("error getting Image Builder Component: empty result")
	}

	component := output.Component

	d.SetId(aws.StringValue(component.Arn))

	d.Set("arn", component.Arn)
	d.Set("change_description", component.ChangeDescription)
	d.Set("data", component.Data)
	d.Set("date_created", component.DateCreated)
	d.Set("description", component.Description)
	d.Set("encrypted", component.Encrypted)
	d.Set("kms_key_id", component.KmsKeyId)
	d.Set("name", component.Name)
	d.Set("owner", component.Owner)
	d.Set("platform", component.Platform)
	d.Set("supported_os_versions", aws.StringValueSlice(component.SupportedOsVersions))

	if err := d.Set("tags", KeyValueTags(component.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	d.Set("type", component.Type)
	d.Set("version", component.Version)

	return nil
}
