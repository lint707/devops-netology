package apigateway

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceCreate,
		Read:   resourceResourceRead,
		Update: resourceResourceUpdate,
		Delete: resourceResourceDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected REST-API-ID/RESOURCE-ID", d.Id())
				}
				restApiID := idParts[0]
				resourceID := idParts[1]
				d.Set("rest_api_id", restApiID)
				d.SetId(resourceID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"path_part": {
				Type:     schema.TypeString,
				Required: true,
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn
	log.Printf("[DEBUG] Creating API Gateway Resource for API %s", d.Get("rest_api_id").(string))

	var err error
	resource, err := conn.CreateResource(&apigateway.CreateResourceInput{
		ParentId:  aws.String(d.Get("parent_id").(string)),
		PathPart:  aws.String(d.Get("path_part").(string)),
		RestApiId: aws.String(d.Get("rest_api_id").(string)),
	})

	if err != nil {
		return fmt.Errorf("Error creating API Gateway Resource: %s", err)
	}

	d.SetId(aws.StringValue(resource.Id))

	return resourceResourceRead(d, meta)
}

func resourceResourceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	log.Printf("[DEBUG] Reading API Gateway Resource %s", d.Id())
	resource, err := conn.GetResource(&apigateway.GetResourceInput{
		ResourceId: aws.String(d.Id()),
		RestApiId:  aws.String(d.Get("rest_api_id").(string)),
	})

	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
			log.Printf("[WARN] API Gateway Resource (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading API Gateway Resource (%s): %w", d.Id(), err)
	}

	d.Set("parent_id", resource.ParentId)
	d.Set("path_part", resource.PathPart)
	d.Set("path", resource.Path)

	return nil
}

func resourceResourceUpdateOperations(d *schema.ResourceData) []*apigateway.PatchOperation {
	operations := make([]*apigateway.PatchOperation, 0)
	if d.HasChange("path_part") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/pathPart"),
			Value: aws.String(d.Get("path_part").(string)),
		})
	}

	if d.HasChange("parent_id") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/parentId"),
			Value: aws.String(d.Get("parent_id").(string)),
		})
	}
	return operations
}

func resourceResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	log.Printf("[DEBUG] Updating API Gateway Resource %s", d.Id())
	_, err := conn.UpdateResource(&apigateway.UpdateResourceInput{
		ResourceId:      aws.String(d.Id()),
		RestApiId:       aws.String(d.Get("rest_api_id").(string)),
		PatchOperations: resourceResourceUpdateOperations(d),
	})

	if err != nil {
		return err
	}

	return resourceResourceRead(d, meta)
}

func resourceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn
	log.Printf("[DEBUG] Deleting API Gateway Resource: %s", d.Id())

	log.Printf("[DEBUG] schema is %#v", d)
	_, err := conn.DeleteResource(&apigateway.DeleteResourceInput{
		ResourceId: aws.String(d.Id()),
		RestApiId:  aws.String(d.Get("rest_api_id").(string)),
	})

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting API Gateway Resource: %s", err)
	}
	return nil
}
