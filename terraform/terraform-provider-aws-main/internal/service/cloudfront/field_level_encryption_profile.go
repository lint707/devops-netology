package cloudfront

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceFieldLevelEncryptionProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceFieldLevelEncryptionProfileCreate,
		Read:   resourceFieldLevelEncryptionProfileRead,
		Update: resourceFieldLevelEncryptionProfileUpdate,
		Delete: resourceFieldLevelEncryptionProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"caller_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_entities": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_patterns": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"items": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"provider_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"public_key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFieldLevelEncryptionProfileCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	apiObject := &cloudfront.FieldLevelEncryptionProfileConfig{
		CallerReference: aws.String(resource.UniqueId()),
		Name:            aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("comment"); ok {
		apiObject.Comment = aws.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_entities"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		apiObject.EncryptionEntities = expandEncryptionEntities(v.([]interface{})[0].(map[string]interface{}))
	}

	input := &cloudfront.CreateFieldLevelEncryptionProfileInput{
		FieldLevelEncryptionProfileConfig: apiObject,
	}

	log.Printf("[DEBUG] Creating CloudFront Field-level Encryption Profile: (%s)", input)
	output, err := conn.CreateFieldLevelEncryptionProfile(input)

	if err != nil {
		return fmt.Errorf("error creating CloudFront Field-level Encryption Profile (%s): %w", d.Id(), err)
	}

	d.SetId(aws.StringValue(output.FieldLevelEncryptionProfile.Id))

	return resourceFieldLevelEncryptionProfileRead(d, meta)
}

func resourceFieldLevelEncryptionProfileRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	output, err := FindFieldLevelEncryptionProfileByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudFront Field-level Encryption Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading CloudFront Field-level Encryption Profile (%s): %w", d.Id(), err)
	}

	apiObject := output.FieldLevelEncryptionProfile.FieldLevelEncryptionProfileConfig
	d.Set("caller_reference", apiObject.CallerReference)
	d.Set("comment", apiObject.Comment)
	if apiObject.EncryptionEntities != nil {
		if err := d.Set("encryption_entities", []interface{}{flattenEncryptionEntities(apiObject.EncryptionEntities)}); err != nil {
			return fmt.Errorf("error setting encryption_entities: %w", err)
		}
	} else {
		d.Set("encryption_entities", nil)
	}
	d.Set("etag", output.ETag)
	d.Set("name", apiObject.Name)

	return nil
}

func resourceFieldLevelEncryptionProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	apiObject := &cloudfront.FieldLevelEncryptionProfileConfig{
		CallerReference: aws.String(d.Get("caller_reference").(string)),
		Name:            aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("comment"); ok {
		apiObject.Comment = aws.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_entities"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		apiObject.EncryptionEntities = expandEncryptionEntities(v.([]interface{})[0].(map[string]interface{}))
	}

	input := &cloudfront.UpdateFieldLevelEncryptionProfileInput{
		FieldLevelEncryptionProfileConfig: apiObject,
		Id:                                aws.String(d.Id()),
		IfMatch:                           aws.String(d.Get("etag").(string)),
	}

	log.Printf("[DEBUG] Updating CloudFront Field-level Encryption Profile: (%s)", input)
	_, err := conn.UpdateFieldLevelEncryptionProfile(input)

	if err != nil {
		return fmt.Errorf("error updating CloudFront Field-level Encryption Profile (%s): %w", d.Id(), err)
	}

	return resourceFieldLevelEncryptionProfileRead(d, meta)
}

func resourceFieldLevelEncryptionProfileDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudFrontConn

	log.Printf("[DEBUG] Deleting CloudFront Field-level Encryption Profile: (%s)", d.Id())
	_, err := conn.DeleteFieldLevelEncryptionProfile(&cloudfront.DeleteFieldLevelEncryptionProfileInput{
		Id:      aws.String(d.Id()),
		IfMatch: aws.String(d.Get("etag").(string)),
	})

	if tfawserr.ErrCodeEquals(err, cloudfront.ErrCodeNoSuchFieldLevelEncryptionProfile) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting CloudFront Field-level Encryption Profile (%s): %w", d.Id(), err)
	}

	return nil
}

func expandEncryptionEntities(tfMap map[string]interface{}) *cloudfront.EncryptionEntities {
	if tfMap == nil {
		return nil
	}

	apiObject := &cloudfront.EncryptionEntities{}

	if v, ok := tfMap["items"].(*schema.Set); ok && v.Len() > 0 {
		items := expandEncryptionEntityItems(v.List())
		apiObject.Items = items
		apiObject.Quantity = aws.Int64(int64(len(items)))
	}

	return apiObject
}

func expandEncryptionEntity(tfMap map[string]interface{}) *cloudfront.EncryptionEntity {
	if tfMap == nil {
		return nil
	}

	apiObject := &cloudfront.EncryptionEntity{}

	if v, ok := tfMap["field_patterns"].([]interface{}); ok && len(v) > 0 {
		apiObject.FieldPatterns = expandFieldPatterns(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["provider_id"].(string); ok && v != "" {
		apiObject.ProviderId = aws.String(v)
	}

	if v, ok := tfMap["public_key_id"].(string); ok && v != "" {
		apiObject.PublicKeyId = aws.String(v)
	}

	return apiObject
}

func expandEncryptionEntityItems(tfList []interface{}) []*cloudfront.EncryptionEntity {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*cloudfront.EncryptionEntity

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandEncryptionEntity(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func expandFieldPatterns(tfMap map[string]interface{}) *cloudfront.FieldPatterns {
	if tfMap == nil {
		return nil
	}

	apiObject := &cloudfront.FieldPatterns{}

	if v, ok := tfMap["items"].(*schema.Set); ok && v.Len() > 0 {
		items := flex.ExpandStringSet(v)
		apiObject.Items = items
		apiObject.Quantity = aws.Int64(int64(len(items)))
	}

	return apiObject
}

func flattenEncryptionEntities(apiObject *cloudfront.EncryptionEntities) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Items; len(v) > 0 {
		tfMap["items"] = flattenEncryptionEntityItems(v)
	}

	return tfMap
}

func flattenEncryptionEntity(apiObject *cloudfront.EncryptionEntity) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := flattenFieldPatterns(apiObject.FieldPatterns); len(v) > 0 {
		tfMap["field_patterns"] = []interface{}{v}
	}

	if v := apiObject.ProviderId; v != nil {
		tfMap["provider_id"] = aws.StringValue(v)
	}

	if v := apiObject.PublicKeyId; v != nil {
		tfMap["public_key_id"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenEncryptionEntityItems(apiObjects []*cloudfront.EncryptionEntity) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		if v := flattenEncryptionEntity(apiObject); len(v) > 0 {
			tfList = append(tfList, v)
		}
	}

	return tfList
}

func flattenFieldPatterns(apiObject *cloudfront.FieldPatterns) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Items; len(v) > 0 {
		tfMap["items"] = aws.StringValueSlice(v)
	}

	return tfMap
}
