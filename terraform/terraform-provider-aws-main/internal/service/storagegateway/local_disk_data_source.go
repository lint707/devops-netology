package storagegateway

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceLocalDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalDiskRead,

		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_node": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disk_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gateway_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func dataSourceLocalDiskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).StorageGatewayConn

	input := &storagegateway.ListLocalDisksInput{
		GatewayARN: aws.String(d.Get("gateway_arn").(string)),
	}

	log.Printf("[DEBUG] Reading Storage Gateway Local Disk: %s", input)
	output, err := conn.ListLocalDisks(input)
	if err != nil {
		return fmt.Errorf("error reading Storage Gateway Local Disk: %w", err)
	}

	if output == nil || len(output.Disks) == 0 {
		return errors.New("no results found for query, try adjusting your search criteria")
	}

	var matchingDisks []*storagegateway.Disk

	for _, disk := range output.Disks {
		if v, ok := d.GetOk("disk_node"); ok && v.(string) == aws.StringValue(disk.DiskNode) {
			matchingDisks = append(matchingDisks, disk)
			continue
		}
		if v, ok := d.GetOk("disk_path"); ok && v.(string) == aws.StringValue(disk.DiskPath) {
			matchingDisks = append(matchingDisks, disk)
			continue
		}
	}

	if len(matchingDisks) == 0 {
		return errors.New("no results found for query, try adjusting your search criteria")
	}

	if len(matchingDisks) > 1 {
		return errors.New("multiple results found for query, try adjusting your search criteria")
	}

	disk := matchingDisks[0]

	d.SetId(aws.StringValue(disk.DiskId))
	d.Set("disk_id", disk.DiskId)
	d.Set("disk_node", disk.DiskNode)
	d.Set("disk_path", disk.DiskPath)

	return nil
}
