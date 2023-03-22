package apigateway

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceRestAPIPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceRestAPIPolicyPut,
		Read:   resourceRestAPIPolicyRead,
		Update: resourceRestAPIPolicyPut,
		Delete: resourceRestAPIPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceRestAPIPolicyPut(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	restApiId := d.Get("rest_api_id").(string)
	log.Printf("[DEBUG] Setting API Gateway REST API Policy: %s", restApiId)

	operations := make([]*apigateway.PatchOperation, 0)

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
	}

	operations = append(operations, &apigateway.PatchOperation{
		Op:    aws.String(apigateway.OpReplace),
		Path:  aws.String("/policy"),
		Value: aws.String(policy),
	})

	res, err := conn.UpdateRestApi(&apigateway.UpdateRestApiInput{
		RestApiId:       aws.String(restApiId),
		PatchOperations: operations,
	})

	if err != nil {
		return fmt.Errorf("error setting API Gateway REST API Policy %w", err)
	}

	log.Printf("[DEBUG] API Gateway REST API Policy Set: %s", restApiId)

	d.SetId(aws.StringValue(res.Id))

	return resourceRestAPIPolicyRead(d, meta)
}

func resourceRestAPIPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	log.Printf("[DEBUG] Reading API Gateway REST API Policy %s", d.Id())

	api, err := conn.GetRestApi(&apigateway.GetRestApiInput{
		RestApiId: aws.String(d.Id()),
	})
	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		log.Printf("[WARN] API Gateway REST API Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading API Gateway REST API Policy (%s): %w", d.Id(), err)
	}

	normalizedPolicy, err := structure.NormalizeJsonString(`"` + aws.StringValue(api.Policy) + `"`)
	if err != nil {
		return fmt.Errorf("error normalizing API Gateway REST API policy JSON: %w", err)
	}

	policy, err := strconv.Unquote(normalizedPolicy)
	if err != nil {
		return fmt.Errorf("error unescaping API Gateway REST API policy: %w", err)
	}

	policyToSet, err := verify.SecondJSONUnlessEquivalent(d.Get("policy").(string), policy)

	if err != nil {
		return fmt.Errorf("while setting policy (%s), encountered: %w", policyToSet, err)
	}

	d.Set("policy", policyToSet)

	d.Set("rest_api_id", api.Id)

	return nil
}

func resourceRestAPIPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	restApiId := d.Get("rest_api_id").(string)
	log.Printf("[DEBUG] Deleting API Gateway REST API Policy: %s", restApiId)

	operations := make([]*apigateway.PatchOperation, 0)

	operations = append(operations, &apigateway.PatchOperation{
		Op:    aws.String(apigateway.OpReplace),
		Path:  aws.String("/policy"),
		Value: aws.String(""),
	})

	_, err := conn.UpdateRestApi(&apigateway.UpdateRestApiInput{
		RestApiId:       aws.String(restApiId),
		PatchOperations: operations,
	})

	if err != nil {
		return fmt.Errorf("error deleting API Gateway REST API policy: %w", err)
	}

	log.Printf("[DEBUG] API Gateway REST API Policy Deleted: %s", restApiId)

	return nil
}
