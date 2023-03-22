package glue

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourcePolicyPut(glue.ExistConditionNotExist),
		Read:   resourceResourcePolicyRead,
		Update: resourceResourcePolicyPut(glue.ExistConditionMustExist),
		Delete: resourceResourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
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
			"enable_hybrid": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(glue.EnableHybridValues_Values(), false),
			},
		},
	}
}

func resourceResourcePolicyPut(condition string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		conn := meta.(*conns.AWSClient).GlueConn

		policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

		if err != nil {
			return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
		}

		input := &glue.PutResourcePolicyInput{
			PolicyInJson:          aws.String(policy),
			PolicyExistsCondition: aws.String(condition),
		}

		if v, ok := d.GetOk("enable_hybrid"); ok {
			input.EnableHybrid = aws.String(v.(string))
		}

		_, err = conn.PutResourcePolicy(input)
		if err != nil {
			return fmt.Errorf("error putting policy request: %s", err)
		}
		d.SetId(meta.(*conns.AWSClient).Region)
		return resourceResourcePolicyRead(d, meta)
	}
}

func resourceResourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GlueConn

	resourcePolicy, err := conn.GetResourcePolicy(&glue.GetResourcePolicyInput{})
	if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
		log.Printf("[WARN] Glue Resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading policy request: %w", err)
	}

	if aws.StringValue(resourcePolicy.PolicyInJson) == "" {
		//Since the glue resource policy is global we expect it to be deleted when the policy is empty
		d.SetId("")
	} else {
		policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), aws.StringValue(resourcePolicy.PolicyInJson))

		if err != nil {
			return err
		}

		d.Set("policy", policyToSet)
	}
	return nil
}

func resourceResourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GlueConn

	_, err := conn.DeleteResourcePolicy(&glue.DeleteResourcePolicyInput{})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			return nil
		}
		return fmt.Errorf("error deleting policy request: %w", err)
	}

	return nil
}
