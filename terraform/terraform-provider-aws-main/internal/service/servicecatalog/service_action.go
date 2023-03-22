package servicecatalog

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceServiceAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceActionCreate,
		Read:   resourceServiceActionRead,
		Update: resourceServiceActionUpdate,
		Delete: resourceServiceActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ServiceActionReadyTimeout),
			Read:   schema.DefaultTimeout(ServiceActionReadTimeout),
			Update: schema.DefaultTimeout(ServiceActionUpdateTimeout),
			Delete: schema.DefaultTimeout(ServiceActionDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"definition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assume_role": { // ServiceActionDefinitionKeyAssumeRole
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": { // ServiceActionDefinitionKeyName
							Type:     schema.TypeString,
							Required: true,
						},
						"parameters": { // ServiceActionDefinitionKeyParameters
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: suppressEquivalentJSONEmptyNilDiffs,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      servicecatalog.ServiceActionDefinitionTypeSsmAutomation,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(servicecatalog.ServiceActionDefinitionType_Values(), false),
						},
						"version": { // ServiceActionDefinitionKeyVersion
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServiceActionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.CreateServiceActionInput{
		IdempotencyToken: aws.String(resource.UniqueId()),
		Name:             aws.String(d.Get("name").(string)),
		Definition:       expandServiceActionDefinition(d.Get("definition").([]interface{})[0].(map[string]interface{})),
		DefinitionType:   aws.String(d.Get("definition.0.type").(string)),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	var output *servicecatalog.CreateServiceActionOutput
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		output, err = conn.CreateServiceAction(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateServiceAction(input)
	}

	if err != nil {
		return fmt.Errorf("error creating Service Catalog Service Action: %w", err)
	}

	if output == nil || output.ServiceActionDetail == nil || output.ServiceActionDetail.ServiceActionSummary == nil {
		return fmt.Errorf("error creating Service Catalog Service Action: empty response")
	}

	d.SetId(aws.StringValue(output.ServiceActionDetail.ServiceActionSummary.Id))

	return resourceServiceActionRead(d, meta)
}

func resourceServiceActionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	output, err := WaitServiceActionReady(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Service Catalog Service Action (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Service Catalog Service Action (%s): %w", d.Id(), err)
	}

	if output == nil || output.ServiceActionSummary == nil {
		return fmt.Errorf("error getting Service Catalog Service Action (%s): empty response", d.Id())
	}

	sas := output.ServiceActionSummary

	d.Set("description", sas.Description)
	d.Set("name", sas.Name)

	if output.Definition != nil {
		d.Set("definition", []interface{}{flattenServiceActionDefinition(output.Definition, aws.StringValue(sas.DefinitionType))})
	} else {
		d.Set("definition", nil)
	}

	return nil
}

func resourceServiceActionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.UpdateServiceActionInput{
		Id: aws.String(d.Id()),
	}

	if d.HasChange("accept_language") {
		input.AcceptLanguage = aws.String(d.Get("accept_language").(string))
	}

	if d.HasChange("definition") {
		input.Definition = expandServiceActionDefinition(d.Get("definition").([]interface{})[0].(map[string]interface{}))
	}

	if d.HasChange("description") {
		input.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChange("name") {
		input.Name = aws.String(d.Get("name").(string))
	}

	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := conn.UpdateServiceAction(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.UpdateServiceAction(input)
	}

	if err != nil {
		return fmt.Errorf("error updating Service Catalog Service Action (%s): %w", d.Id(), err)
	}

	return resourceServiceActionRead(d, meta)
}

func resourceServiceActionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.DeleteServiceActionInput{
		Id: aws.String(d.Id()),
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteServiceAction(input)

		if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceInUseException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteServiceAction(input)
	}

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		log.Printf("[INFO] Attempted to delete Service Action (%s) but does not exist", d.Id())
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Service Catalog Service Action (%s): %w", d.Id(), err)
	}

	if err := WaitServiceActionDeleted(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for Service Catalog Service Action (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}

func expandServiceActionDefinition(tfMap map[string]interface{}) map[string]*string {
	if tfMap == nil {
		return nil
	}

	apiObject := make(map[string]*string)

	if v, ok := tfMap["assume_role"].(string); ok && v != "" {
		apiObject[servicecatalog.ServiceActionDefinitionKeyAssumeRole] = aws.String(v)
	}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject[servicecatalog.ServiceActionDefinitionKeyName] = aws.String(v)
	}

	if v, ok := tfMap["parameters"].(string); ok && v != "" {
		apiObject[servicecatalog.ServiceActionDefinitionKeyParameters] = aws.String(v)
	}

	if v, ok := tfMap["version"].(string); ok && v != "" {
		apiObject[servicecatalog.ServiceActionDefinitionKeyVersion] = aws.String(v)
	}

	return apiObject
}

func flattenServiceActionDefinition(apiObject map[string]*string, definitionType string) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v, ok := apiObject[servicecatalog.ServiceActionDefinitionKeyAssumeRole]; ok && v != nil {
		tfMap["assume_role"] = aws.StringValue(v)
	}

	if v, ok := apiObject[servicecatalog.ServiceActionDefinitionKeyName]; ok && v != nil {
		tfMap["name"] = aws.StringValue(v)
	}

	if v, ok := apiObject[servicecatalog.ServiceActionDefinitionKeyParameters]; ok && v != nil {
		tfMap["parameters"] = aws.StringValue(v)
	}

	if v, ok := apiObject[servicecatalog.ServiceActionDefinitionKeyVersion]; ok && v != nil {
		tfMap["version"] = aws.StringValue(v)
	}

	if definitionType != "" {
		tfMap["type"] = definitionType
	}

	return tfMap
}
