package s3control

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketPolicyCreate,
		Read:   resourceBucketPolicyRead,
		Update: resourceBucketPolicyUpdate,
		Delete: resourceBucketPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	}
}

func resourceBucketPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	bucket := d.Get("bucket").(string)

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", d.Get("policy").(string), err)
	}

	input := &s3control.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	}

	_, err = conn.PutBucketPolicy(input)

	if err != nil {
		return fmt.Errorf("error creating S3 Control Bucket Policy (%s): %w", bucket, err)
	}

	d.SetId(bucket)

	return resourceBucketPolicyRead(d, meta)
}

func resourceBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	if parsedArn.AccountID == "" {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): unknown format", d.Id())
	}

	input := &s3control.GetBucketPolicyInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	output, err := conn.GetBucketPolicy(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchBucketPolicy") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error reading S3 Control Bucket Policy (%s): empty response", d.Id())
	}

	d.Set("bucket", d.Id())

	if output.Policy != nil {
		policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), aws.StringValue(output.Policy))

		if err != nil {
			return err
		}

		d.Set("policy", policyToSet)
	} else {
		d.Set("policy", "")
	}

	return nil
}

func resourceBucketPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", d.Get("policy").(string), err)
	}

	input := &s3control.PutBucketPolicyInput{
		Bucket: aws.String(d.Id()),
		Policy: aws.String(policy),
	}

	_, err = conn.PutBucketPolicy(input)

	if err != nil {
		return fmt.Errorf("error updating S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	return resourceBucketPolicyRead(d, meta)
}

func resourceBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	input := &s3control.DeleteBucketPolicyInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	_, err = conn.DeleteBucketPolicy(input)

	if tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchBucketPolicy") {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	return nil
}
