package waf

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceRegexMatchSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegexMatchSetCreate,
		Read:   resourceRegexMatchSetRead,
		Update: resourceRegexMatchSetUpdate,
		Delete: resourceRegexMatchSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regex_match_tuple": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      RegexMatchSetTupleHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_to_match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Type:     schema.TypeString,
										Optional: true,
										StateFunc: func(v interface{}) string {
											return strings.ToLower(v.(string))
										},
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"regex_pattern_set_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"text_transformation": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRegexMatchSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	log.Printf("[INFO] Creating WAF Regex Match Set: %s", d.Get("name").(string))

	wr := NewRetryer(conn)
	out, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		params := &waf.CreateRegexMatchSetInput{
			ChangeToken: token,
			Name:        aws.String(d.Get("name").(string)),
		}
		return conn.CreateRegexMatchSet(params)
	})
	if err != nil {
		return fmt.Errorf("Failed creating WAF Regex Match Set: %s", err)
	}
	resp := out.(*waf.CreateRegexMatchSetOutput)

	d.SetId(aws.StringValue(resp.RegexMatchSet.RegexMatchSetId))

	return resourceRegexMatchSetUpdate(d, meta)
}

func resourceRegexMatchSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn
	log.Printf("[INFO] Reading WAF Regex Match Set: %s", d.Get("name").(string))
	params := &waf.GetRegexMatchSetInput{
		RegexMatchSetId: aws.String(d.Id()),
	}

	resp, err := conn.GetRegexMatchSet(params)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
			log.Printf("[WARN] WAF Regex Match Set (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", resp.RegexMatchSet.Name)
	d.Set("regex_match_tuple", FlattenRegexMatchTuples(resp.RegexMatchSet.RegexMatchTuples))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "waf",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("regexmatchset/%s", d.Id()),
	}
	d.Set("arn", arn.String())

	return nil
}

func resourceRegexMatchSetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	log.Printf("[INFO] Updating WAF Regex Match Set: %s", d.Get("name").(string))

	if d.HasChange("regex_match_tuple") {
		o, n := d.GetChange("regex_match_tuple")
		oldT, newT := o.(*schema.Set).List(), n.(*schema.Set).List()
		err := updateRegexMatchSetResource(d.Id(), oldT, newT, conn)
		if err != nil {
			return fmt.Errorf("Failed updating WAF Regex Match Set: %s", err)
		}
	}

	return resourceRegexMatchSetRead(d, meta)
}

func resourceRegexMatchSetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	oldTuples := d.Get("regex_match_tuple").(*schema.Set).List()
	if len(oldTuples) > 0 {
		noTuples := []interface{}{}
		err := updateRegexMatchSetResource(d.Id(), oldTuples, noTuples, conn)
		if err != nil {
			return fmt.Errorf("Error updating WAF Regex Match Set: %s", err)
		}
	}

	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.DeleteRegexMatchSetInput{
			ChangeToken:     token,
			RegexMatchSetId: aws.String(d.Id()),
		}
		log.Printf("[INFO] Deleting WAF Regex Match Set: %s", req)
		return conn.DeleteRegexMatchSet(req)
	})
	if err != nil {
		return fmt.Errorf("Failed deleting WAF Regex Match Set: %s", err)
	}

	return nil
}

func updateRegexMatchSetResource(id string, oldT, newT []interface{}, conn *waf.WAF) error {
	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.UpdateRegexMatchSetInput{
			ChangeToken:     token,
			RegexMatchSetId: aws.String(id),
			Updates:         DiffRegexMatchSetTuples(oldT, newT),
		}

		return conn.UpdateRegexMatchSet(req)
	})
	if err != nil {
		return fmt.Errorf("Failed updating WAF Regex Match Set: %s", err)
	}

	return nil
}
