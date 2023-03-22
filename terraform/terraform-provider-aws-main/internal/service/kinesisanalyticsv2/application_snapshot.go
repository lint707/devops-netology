package kinesisanalyticsv2

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceApplicationSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSnapshotCreate,
		Read:   resourceApplicationSnapshotRead,
		Delete: resourceApplicationSnapshotDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`), "must only include alphanumeric, underscore, period, or hyphen characters"),
				),
			},

			"application_version_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"snapshot_creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 256),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`), "must only include alphanumeric, underscore, period, or hyphen characters"),
				),
			},
		},
	}
}

func resourceApplicationSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KinesisAnalyticsV2Conn
	applicationName := d.Get("application_name").(string)
	snapshotName := d.Get("snapshot_name").(string)

	input := &kinesisanalyticsv2.CreateApplicationSnapshotInput{
		ApplicationName: aws.String(applicationName),
		SnapshotName:    aws.String(snapshotName),
	}

	log.Printf("[DEBUG] Creating Kinesis Analytics v2 Application Snapshot: %s", input)

	_, err := conn.CreateApplicationSnapshot(input)

	if err != nil {
		return fmt.Errorf("error creating Kinesis Analytics v2 Application Snapshot (%s/%s): %w", applicationName, snapshotName, err)
	}

	d.SetId(applicationSnapshotCreateID(applicationName, snapshotName))

	_, err = waitSnapshotCreated(conn, applicationName, snapshotName, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("error waiting for Kinesis Analytics v2 Application Snapshot (%s) creation: %w", d.Id(), err)
	}

	return resourceApplicationSnapshotRead(d, meta)
}

func resourceApplicationSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KinesisAnalyticsV2Conn

	applicationName, snapshotName, err := applicationSnapshotParseID(d.Id())

	if err != nil {
		return err
	}

	snapshot, err := FindSnapshotDetailsByApplicationAndSnapshotNames(conn, applicationName, snapshotName)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Kinesis Analytics v2 Application Snapshot (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Kinesis Analytics v2 Application Snapshot (%s): %w", d.Id(), err)
	}

	d.Set("application_name", applicationName)
	d.Set("application_version_id", snapshot.ApplicationVersionId)
	d.Set("snapshot_creation_timestamp", aws.TimeValue(snapshot.SnapshotCreationTimestamp).Format(time.RFC3339))
	d.Set("snapshot_name", snapshotName)

	return nil
}

func resourceApplicationSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KinesisAnalyticsV2Conn

	applicationName, snapshotName, err := applicationSnapshotParseID(d.Id())

	if err != nil {
		return err
	}

	snapshotCreationTimestamp, err := time.Parse(time.RFC3339, d.Get("snapshot_creation_timestamp").(string))
	if err != nil {
		return fmt.Errorf("error parsing snapshot_creation_timestamp: %w", err)
	}

	log.Printf("[DEBUG] Deleting Kinesis Analytics v2 Application Snapshot (%s)", d.Id())
	_, err = conn.DeleteApplicationSnapshot(&kinesisanalyticsv2.DeleteApplicationSnapshotInput{
		ApplicationName:           aws.String(applicationName),
		SnapshotCreationTimestamp: aws.Time(snapshotCreationTimestamp),
		SnapshotName:              aws.String(snapshotName),
	})

	if tfawserr.ErrCodeEquals(err, kinesisanalyticsv2.ErrCodeResourceNotFoundException) {
		return nil
	}

	if tfawserr.ErrMessageContains(err, kinesisanalyticsv2.ErrCodeInvalidArgumentException, "does not exist") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Kinesis Analytics v2 Application Snapshot (%s): %w", d.Id(), err)
	}

	_, err = waitSnapshotDeleted(conn, applicationName, snapshotName, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("error waiting for Kinesis Analytics v2 Application Snapshot (%s) deletion: %w", d.Id(), err)
	}

	return nil
}
