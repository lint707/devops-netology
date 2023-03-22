package waf

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceGeoMatchSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeoMatchSetCreate,
		Read:   resourceGeoMatchSetRead,
		Update: resourceGeoMatchSetUpdate,
		Delete: resourceGeoMatchSetDelete,
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
			"geo_match_constraint": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceGeoMatchSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	log.Printf("[INFO] Creating GeoMatchSet: %s", d.Get("name").(string))

	wr := NewRetryer(conn)
	out, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		params := &waf.CreateGeoMatchSetInput{
			ChangeToken: token,
			Name:        aws.String(d.Get("name").(string)),
		}

		return conn.CreateGeoMatchSet(params)
	})
	if err != nil {
		return fmt.Errorf("Error creating GeoMatchSet: %s", err)
	}
	resp := out.(*waf.CreateGeoMatchSetOutput)

	d.SetId(aws.StringValue(resp.GeoMatchSet.GeoMatchSetId))

	return resourceGeoMatchSetUpdate(d, meta)
}

func resourceGeoMatchSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn
	log.Printf("[INFO] Reading GeoMatchSet: %s", d.Get("name").(string))
	params := &waf.GetGeoMatchSetInput{
		GeoMatchSetId: aws.String(d.Id()),
	}

	resp, err := conn.GetGeoMatchSet(params)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
			log.Printf("[WARN] WAF GeoMatchSet (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", resp.GeoMatchSet.Name)
	d.Set("geo_match_constraint", FlattenGeoMatchConstraint(resp.GeoMatchSet.GeoMatchConstraints))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "waf",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("geomatchset/%s", d.Id()),
	}
	d.Set("arn", arn.String())

	return nil
}

func resourceGeoMatchSetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	if d.HasChange("geo_match_constraint") {
		o, n := d.GetChange("geo_match_constraint")
		oldT, newT := o.(*schema.Set).List(), n.(*schema.Set).List()

		err := updateGeoMatchSetResource(d.Id(), oldT, newT, conn)
		if err != nil {
			return fmt.Errorf("Error updating GeoMatchSet: %s", err)
		}
	}

	return resourceGeoMatchSetRead(d, meta)
}

func resourceGeoMatchSetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	oldConstraints := d.Get("geo_match_constraint").(*schema.Set).List()
	if len(oldConstraints) > 0 {
		noConstraints := []interface{}{}
		err := updateGeoMatchSetResource(d.Id(), oldConstraints, noConstraints, conn)
		if err != nil {
			return fmt.Errorf("Error updating GeoMatchConstraint: %s", err)
		}
	}

	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.DeleteGeoMatchSetInput{
			ChangeToken:   token,
			GeoMatchSetId: aws.String(d.Id()),
		}

		return conn.DeleteGeoMatchSet(req)
	})
	if err != nil {
		return fmt.Errorf("Error deleting GeoMatchSet: %s", err)
	}

	return nil
}

func updateGeoMatchSetResource(id string, oldT, newT []interface{}, conn *waf.WAF) error {
	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.UpdateGeoMatchSetInput{
			ChangeToken:   token,
			GeoMatchSetId: aws.String(id),
			Updates:       DiffGeoMatchSetConstraints(oldT, newT),
		}

		log.Printf("[INFO] Updating GeoMatchSet constraints: %s", req)
		return conn.UpdateGeoMatchSet(req)
	})
	if err != nil {
		return fmt.Errorf("Error updating GeoMatchSet: %s", err)
	}

	return nil
}
