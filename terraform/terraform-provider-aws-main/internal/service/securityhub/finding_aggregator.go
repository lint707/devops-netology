package securityhub

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

const (
	allRegions                = "ALL_REGIONS"
	allRegionsExceptSpecified = "ALL_REGIONS_EXCEPT_SPECIFIED"
	specifiedRegions          = "SPECIFIED_REGIONS"
)

func ResourceFindingAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceFindingAggregatorCreate,
		Read:   resourceFindingAggregatorRead,
		Update: resourceFindingAggregatorUpdate,
		Delete: resourceFindingAggregatorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"linking_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					allRegions,
					allRegionsExceptSpecified,
					specifiedRegions,
				}, false),
			},
			"specified_regions": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceFindingAggregatorCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	linkingMode := d.Get("linking_mode").(string)

	req := &securityhub.CreateFindingAggregatorInput{
		RegionLinkingMode: &linkingMode,
	}

	if v, ok := d.GetOk("specified_regions"); ok && (linkingMode == allRegionsExceptSpecified || linkingMode == specifiedRegions) {
		req.Regions = flex.ExpandStringSet(v.(*schema.Set))
	}

	log.Printf("[DEBUG] Creating Security Hub finding aggregator")

	resp, err := conn.CreateFindingAggregator(req)

	if err != nil {
		return fmt.Errorf("Error creating finding aggregator for Security Hub: %s", err)
	}

	d.SetId(aws.StringValue(resp.FindingAggregatorArn))

	return resourceFindingAggregatorRead(d, meta)
}

func resourceFindingAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	aggregatorArn := d.Id()

	log.Printf("[DEBUG] Reading Security Hub finding aggregator to find %s", aggregatorArn)

	aggregator, err := FindingAggregatorCheckExists(conn, aggregatorArn)

	if err != nil {
		return fmt.Errorf("Error reading Security Hub finding aggregator to find %s: %s", aggregatorArn, err)
	}

	if aggregator == nil {
		log.Printf("[WARN] Security Hub finding aggregator (%s) not found, removing from state", aggregatorArn)
		d.SetId("")
		return nil
	}

	d.Set("linking_mode", aggregator.RegionLinkingMode)

	if len(aggregator.Regions) > 0 {
		d.Set("specified_regions", flex.FlattenStringList(aggregator.Regions))
	}

	return nil
}

func FindingAggregatorCheckExists(conn *securityhub.SecurityHub, findingAggregatorArn string) (*securityhub.GetFindingAggregatorOutput, error) {
	input := &securityhub.ListFindingAggregatorsInput{}

	var found *securityhub.GetFindingAggregatorOutput
	var err error = nil

	err = conn.ListFindingAggregatorsPages(input, func(page *securityhub.ListFindingAggregatorsOutput, lastPage bool) bool {
		for _, aggregator := range page.FindingAggregators {
			if aws.StringValue(aggregator.FindingAggregatorArn) == findingAggregatorArn {
				getInput := &securityhub.GetFindingAggregatorInput{
					FindingAggregatorArn: &findingAggregatorArn,
				}
				found, err = conn.GetFindingAggregator(getInput)
				return false
			}
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return found, nil
}

func resourceFindingAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	aggregatorArn := d.Id()

	linkingMode := d.Get("linking_mode").(string)

	req := &securityhub.UpdateFindingAggregatorInput{
		FindingAggregatorArn: &aggregatorArn,
		RegionLinkingMode:    &linkingMode,
	}

	if v, ok := d.GetOk("specified_regions"); ok && (linkingMode == allRegionsExceptSpecified || linkingMode == specifiedRegions) {
		req.Regions = flex.ExpandStringSet(v.(*schema.Set))
	}

	resp, err := conn.UpdateFindingAggregator(req)

	if err != nil {
		return fmt.Errorf("Error updating Security Hub finding aggregator (%s): %w", aggregatorArn, err)
	}

	d.SetId(aws.StringValue(resp.FindingAggregatorArn))

	return resourceFindingAggregatorRead(d, meta)
}

func resourceFindingAggregatorDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	aggregatorArn := d.Id()

	log.Printf("[DEBUG] Disabling Security Hub finding aggregator %s", aggregatorArn)

	_, err := conn.DeleteFindingAggregator(&securityhub.DeleteFindingAggregatorInput{
		FindingAggregatorArn: &aggregatorArn,
	})

	if err != nil {
		return fmt.Errorf("Error disabling Security Hub finding aggregator %s: %s", aggregatorArn, err)
	}

	return nil
}
