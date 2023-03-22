package waf

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceSizeConstraintSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSizeConstraintSetCreate,
		Read:   resourceSizeConstraintSetRead,
		Update: resourceSizeConstraintSetUpdate,
		Delete: resourceSizeConstraintSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: SizeConstraintSetSchema(),
	}
}

func resourceSizeConstraintSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	log.Printf("[INFO] Creating SizeConstraintSet: %s", d.Get("name").(string))

	wr := NewRetryer(conn)
	out, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		params := &waf.CreateSizeConstraintSetInput{
			ChangeToken: token,
			Name:        aws.String(d.Get("name").(string)),
		}

		return conn.CreateSizeConstraintSet(params)
	})
	if err != nil {
		return fmt.Errorf("Error creating SizeConstraintSet: %s", err)
	}
	resp := out.(*waf.CreateSizeConstraintSetOutput)

	d.SetId(aws.StringValue(resp.SizeConstraintSet.SizeConstraintSetId))

	return resourceSizeConstraintSetUpdate(d, meta)
}

func resourceSizeConstraintSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn
	log.Printf("[INFO] Reading SizeConstraintSet: %s", d.Get("name").(string))
	params := &waf.GetSizeConstraintSetInput{
		SizeConstraintSetId: aws.String(d.Id()),
	}

	resp, err := conn.GetSizeConstraintSet(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == waf.ErrCodeNonexistentItemException {
			log.Printf("[WARN] WAF SizeConstraintSet (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", resp.SizeConstraintSet.Name)
	d.Set("size_constraints", FlattenSizeConstraints(resp.SizeConstraintSet.SizeConstraints))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "waf",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("sizeconstraintset/%s", d.Id()),
	}
	d.Set("arn", arn.String())

	return nil
}

func resourceSizeConstraintSetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	if d.HasChange("size_constraints") {
		o, n := d.GetChange("size_constraints")
		oldConstraints, newConstraints := o.(*schema.Set).List(), n.(*schema.Set).List()

		err := updateSizeConstraintSetResource(d.Id(), oldConstraints, newConstraints, conn)
		if err != nil {
			return fmt.Errorf("Error updating SizeConstraintSet: %s", err)
		}
	}

	return resourceSizeConstraintSetRead(d, meta)
}

func resourceSizeConstraintSetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFConn

	oldConstraints := d.Get("size_constraints").(*schema.Set).List()

	if len(oldConstraints) > 0 {
		noConstraints := []interface{}{}
		err := updateSizeConstraintSetResource(d.Id(), oldConstraints, noConstraints, conn)
		if err != nil {
			return fmt.Errorf("Error deleting SizeConstraintSet: %s", err)
		}
	}

	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.DeleteSizeConstraintSetInput{
			ChangeToken:         token,
			SizeConstraintSetId: aws.String(d.Id()),
		}
		return conn.DeleteSizeConstraintSet(req)
	})
	if err != nil {
		return fmt.Errorf("Error deleting SizeConstraintSet: %s", err)
	}

	return nil
}

func updateSizeConstraintSetResource(id string, oldS, newS []interface{}, conn *waf.WAF) error {
	wr := NewRetryer(conn)
	_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
		req := &waf.UpdateSizeConstraintSetInput{
			ChangeToken:         token,
			SizeConstraintSetId: aws.String(id),
			Updates:             DiffSizeConstraints(oldS, newS),
		}

		log.Printf("[INFO] Updating WAF Size Constraint constraints: %s", req)
		return conn.UpdateSizeConstraintSet(req)
	})
	if err != nil {
		return fmt.Errorf("Error updating SizeConstraintSet: %s", err)
	}

	return nil
}
