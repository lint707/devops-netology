package servicecatalog

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourcePortfolio() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePortfolioRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(ConstraintReadTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "en",
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourcePortfolioRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.DescribePortfolioInput{
		Id: aws.String(d.Get("id").(string)),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	output, err := conn.DescribePortfolio(input)

	if err != nil {
		return fmt.Errorf("error getting Service Catalog Portfolio (%s): %w", d.Get("id").(string), err)
	}

	if output == nil || output.PortfolioDetail == nil {
		return fmt.Errorf("error getting Service Catalog Portfolio (%s): empty response", d.Get("id").(string))
	}

	detail := output.PortfolioDetail

	d.SetId(aws.StringValue(detail.Id))

	if err := d.Set("created_time", aws.TimeValue(detail.CreatedTime).Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Error setting created_time: %s", err)
	}

	d.Set("arn", detail.ARN)
	d.Set("description", detail.Description)
	d.Set("name", detail.DisplayName)
	d.Set("provider_name", detail.ProviderName)

	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	tags := KeyValueTags(output.Tags)

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
