package wafv2

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleGroupCreate,
		Read:   resourceRuleGroupRead,
		Update: resourceRuleGroupUpdate,
		Delete: resourceRuleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected ID/NAME/SCOPE", d.Id())
				}
				id := idParts[0]
				name := idParts[1]
				scope := idParts[2]
				d.SetId(id)
				d.Set("name", name)
				d.Set("scope", scope)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"custom_response_body": customResponseBodySchema(),
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"lock_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-_]+$`), "must contain only alphanumeric hyphen and underscore characters"),
				),
			},
			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					wafv2.ScopeCloudfront,
					wafv2.ScopeRegional,
				}, false),
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow": allowConfigSchema(),
									"block": blockConfigSchema(),
									"count": countConfigSchema(),
								},
							},
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 128),
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_label":        ruleLabelsSchema(),
						"statement":         rootStatementSchema(rootStatementSchemaLevel),
						"visibility_config": visibilityConfigSchema(),
					},
				},
			},
			"tags":              tftags.TagsSchema(),
			"tags_all":          tftags.TagsSchemaComputed(),
			"visibility_config": visibilityConfigSchema(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceRuleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))
	var resp *wafv2.CreateRuleGroupOutput

	params := &wafv2.CreateRuleGroupInput{
		Name:             aws.String(d.Get("name").(string)),
		Scope:            aws.String(d.Get("scope").(string)),
		Capacity:         aws.Int64(int64(d.Get("capacity").(int))),
		Rules:            expandRules(d.Get("rule").(*schema.Set).List()),
		VisibilityConfig: expandVisibilityConfig(d.Get("visibility_config").([]interface{})),
	}

	if v, ok := d.GetOk("custom_response_body"); ok && v.(*schema.Set).Len() > 0 {
		params.CustomResponseBodies = expandCustomResponseBodies(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		params.Tags = Tags(tags.IgnoreAWS())
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.CreateRuleGroup(params)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, wafv2.ErrCodeWAFUnavailableEntityException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		resp, err = conn.CreateRuleGroup(params)
	}

	if err != nil {
		return fmt.Errorf("Error creating WAFv2 RuleGroup: %s", err)
	}

	if resp == nil || resp.Summary == nil {
		return fmt.Errorf("Error creating WAFv2 RuleGroup")
	}

	d.SetId(aws.StringValue(resp.Summary.Id))

	return resourceRuleGroupRead(d, meta)
}

func resourceRuleGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	params := &wafv2.GetRuleGroupInput{
		Id:    aws.String(d.Id()),
		Name:  aws.String(d.Get("name").(string)),
		Scope: aws.String(d.Get("scope").(string)),
	}

	resp, err := conn.GetRuleGroup(params)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, wafv2.ErrCodeWAFNonexistentItemException) {
			log.Printf("[WARN] WAFv2 RuleGroup (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if resp == nil || resp.RuleGroup == nil {
		return fmt.Errorf("Error getting WAFv2 RuleGroup")
	}

	d.Set("name", resp.RuleGroup.Name)
	d.Set("capacity", resp.RuleGroup.Capacity)
	d.Set("description", resp.RuleGroup.Description)
	d.Set("arn", resp.RuleGroup.ARN)
	d.Set("lock_token", resp.LockToken)

	if err := d.Set("custom_response_body", flattenCustomResponseBodies(resp.RuleGroup.CustomResponseBodies)); err != nil {
		return fmt.Errorf("Error setting custom_response_body: %w", err)
	}

	if err := d.Set("rule", flattenRules(resp.RuleGroup.Rules)); err != nil {
		return fmt.Errorf("Error setting rule: %s", err)
	}

	if err := d.Set("visibility_config", flattenVisibilityConfig(resp.RuleGroup.VisibilityConfig)); err != nil {
		return fmt.Errorf("Error setting visibility_config: %s", err)
	}

	arn := aws.StringValue(resp.RuleGroup.ARN)
	tags, err := ListTags(conn, arn)
	if err != nil {
		return fmt.Errorf("Error listing tags for WAFv2 RuleGroup (%s): %s", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceRuleGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn

	log.Printf("[INFO] Updating WAFv2 RuleGroup %s", d.Id())

	u := &wafv2.UpdateRuleGroupInput{
		Id:               aws.String(d.Id()),
		Name:             aws.String(d.Get("name").(string)),
		Scope:            aws.String(d.Get("scope").(string)),
		LockToken:        aws.String(d.Get("lock_token").(string)),
		Rules:            expandRules(d.Get("rule").(*schema.Set).List()),
		VisibilityConfig: expandVisibilityConfig(d.Get("visibility_config").([]interface{})),
	}

	if v, ok := d.GetOk("custom_response_body"); ok && v.(*schema.Set).Len() > 0 {
		u.CustomResponseBodies = expandCustomResponseBodies(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("description"); ok {
		u.Description = aws.String(v.(string))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.UpdateRuleGroup(u)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, wafv2.ErrCodeWAFUnavailableEntityException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.UpdateRuleGroup(u)
	}

	if err != nil {
		return fmt.Errorf("Error updating WAFv2 RuleGroup: %s", err)
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceRuleGroupRead(d, meta)
}

func resourceRuleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn

	log.Printf("[INFO] Deleting WAFv2 RuleGroup %s", d.Id())

	r := &wafv2.DeleteRuleGroupInput{
		Id:        aws.String(d.Id()),
		Name:      aws.String(d.Get("name").(string)),
		Scope:     aws.String(d.Get("scope").(string)),
		LockToken: aws.String(d.Get("lock_token").(string)),
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteRuleGroup(r)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, wafv2.ErrCodeWAFAssociatedItemException) {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrCodeEquals(err, wafv2.ErrCodeWAFUnavailableEntityException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteRuleGroup(r)
	}

	if err != nil {
		return fmt.Errorf("Error deleting WAFv2 RuleGroup: %s", err)
	}

	return nil
}
