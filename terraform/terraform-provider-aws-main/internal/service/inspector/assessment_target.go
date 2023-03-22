package inspector

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceAssessmentTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAssessmentTargetCreate,
		Read:   resourceAssessmentTargetRead,
		Update: resourceAssessmentTargetUpdate,
		Delete: resourceAssessmentTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAssessmentTargetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	input := &inspector.CreateAssessmentTargetInput{
		AssessmentTargetName: aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("resource_group_arn"); ok {
		input.ResourceGroupArn = aws.String(v.(string))
	}

	resp, err := conn.CreateAssessmentTarget(input)
	if err != nil {
		return fmt.Errorf("error creating Inspector Assessment Target: %s", err)
	}

	d.SetId(aws.StringValue(resp.AssessmentTargetArn))

	return resourceAssessmentTargetRead(d, meta)
}

func resourceAssessmentTargetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	assessmentTarget, err := DescribeAssessmentTarget(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error describing Inspector Assessment Target (%s): %s", d.Id(), err)
	}

	if assessmentTarget == nil {
		log.Printf("[WARN] Inspector Assessment Target (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("arn", assessmentTarget.Arn)
	d.Set("name", assessmentTarget.Name)
	d.Set("resource_group_arn", assessmentTarget.ResourceGroupArn)

	return nil
}

func resourceAssessmentTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	input := inspector.UpdateAssessmentTargetInput{
		AssessmentTargetArn:  aws.String(d.Id()),
		AssessmentTargetName: aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("resource_group_arn"); ok {
		input.ResourceGroupArn = aws.String(v.(string))
	}

	_, err := conn.UpdateAssessmentTarget(&input)
	if err != nil {
		return fmt.Errorf("error updating Inspector Assessment Target (%s): %s", d.Id(), err)
	}

	return resourceAssessmentTargetRead(d, meta)
}

func resourceAssessmentTargetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn
	input := &inspector.DeleteAssessmentTargetInput{
		AssessmentTargetArn: aws.String(d.Id()),
	}
	err := resource.Retry(60*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteAssessmentTarget(input)

		if tfawserr.ErrCodeEquals(err, inspector.ErrCodeAssessmentRunInProgressException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteAssessmentTarget(input)
	}
	if err != nil {
		return fmt.Errorf("Error deleting Inspector Assessment Target: %s", err)
	}
	return nil
}

func DescribeAssessmentTarget(conn *inspector.Inspector, arn string) (*inspector.AssessmentTarget, error) {
	input := &inspector.DescribeAssessmentTargetsInput{
		AssessmentTargetArns: []*string{aws.String(arn)},
	}

	output, err := conn.DescribeAssessmentTargets(input)

	if tfawserr.ErrCodeEquals(err, inspector.ErrCodeInvalidInputException) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var assessmentTarget *inspector.AssessmentTarget
	for _, target := range output.AssessmentTargets {
		if aws.StringValue(target.Arn) == arn {
			assessmentTarget = target
			break
		}
	}

	return assessmentTarget, nil
}
