package outposts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/outposts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceOutpost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutpostRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutpostRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).OutpostsConn

	input := &outposts.ListOutpostsInput{}

	var results []*outposts.Outpost

	err := conn.ListOutpostsPages(input, func(page *outposts.ListOutpostsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, outpost := range page.Outposts {
			if outpost == nil {
				continue
			}

			if v, ok := d.GetOk("id"); ok && v.(string) != aws.StringValue(outpost.OutpostId) {
				continue
			}

			if v, ok := d.GetOk("name"); ok && v.(string) != aws.StringValue(outpost.Name) {
				continue
			}

			if v, ok := d.GetOk("arn"); ok && v.(string) != aws.StringValue(outpost.OutpostArn) {
				continue
			}

			if v, ok := d.GetOk("owner_id"); ok && v.(string) != aws.StringValue(outpost.OwnerId) {
				continue
			}

			results = append(results, outpost)
		}

		return !lastPage
	})

	if err != nil {
		return fmt.Errorf("error listing Outposts Outposts: %w", err)
	}

	if len(results) == 0 {
		return fmt.Errorf("no Outposts Outpost found matching criteria; try different search")
	}

	if len(results) > 1 {
		return fmt.Errorf("multiple Outposts Outpost found matching criteria; try different search")
	}

	outpost := results[0]

	d.SetId(aws.StringValue(outpost.OutpostId))
	d.Set("arn", outpost.OutpostArn)
	d.Set("availability_zone", outpost.AvailabilityZone)
	d.Set("availability_zone_id", outpost.AvailabilityZoneId)
	d.Set("description", outpost.Description)
	d.Set("name", outpost.Name)
	d.Set("owner_id", outpost.OwnerId)
	d.Set("site_id", outpost.SiteId)

	return nil
}
