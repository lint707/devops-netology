package securityhub

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceProductSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceProductSubscriptionCreate,
		Read:   resourceProductSubscriptionRead,
		Delete: resourceProductSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProductSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn
	productArn := d.Get("product_arn").(string)

	log.Printf("[DEBUG] Enabling Security Hub product subscription for product %s", productArn)

	resp, err := conn.EnableImportFindingsForProduct(&securityhub.EnableImportFindingsForProductInput{
		ProductArn: aws.String(productArn),
	})

	if err != nil {
		return fmt.Errorf("Error enabling Security Hub product subscription for product %s: %s", productArn, err)
	}

	d.SetId(fmt.Sprintf("%s,%s", productArn, *resp.ProductSubscriptionArn))

	return resourceProductSubscriptionRead(d, meta)
}

func resourceProductSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	productArn, productSubscriptionArn, err := ProductSubscriptionParseID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Security Hub product subscriptions to find %s", d.Id())

	exists, err := ProductSubscriptionCheckExists(conn, productSubscriptionArn)

	if err != nil {
		return fmt.Errorf("Error reading Security Hub product subscriptions to find %s: %s", d.Id(), err)
	}

	if !exists {
		log.Printf("[WARN] Security Hub product subscriptions (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("product_arn", productArn)
	d.Set("arn", productSubscriptionArn)

	return nil
}

func ProductSubscriptionCheckExists(conn *securityhub.SecurityHub, productSubscriptionArn string) (bool, error) {
	input := &securityhub.ListEnabledProductsForImportInput{}
	exists := false

	err := conn.ListEnabledProductsForImportPages(input, func(page *securityhub.ListEnabledProductsForImportOutput, lastPage bool) bool {
		for _, readProductSubscriptionArn := range page.ProductSubscriptions {
			if aws.StringValue(readProductSubscriptionArn) == productSubscriptionArn {
				exists = true
				return false
			}
		}
		return !lastPage
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func ProductSubscriptionParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ",", 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Expected Security Hub product subscription ID in format <product_arn>,<arn> - received: %s", id)
	}

	return parts[0], parts[1], nil
}

func resourceProductSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn
	log.Printf("[DEBUG] Disabling Security Hub product subscription %s", d.Id())

	_, productSubscriptionArn, err := ProductSubscriptionParseID(d.Id())

	if err != nil {
		return err
	}

	_, err = conn.DisableImportFindingsForProduct(&securityhub.DisableImportFindingsForProductInput{
		ProductSubscriptionArn: aws.String(productSubscriptionArn),
	})

	if err != nil {
		return fmt.Errorf("Error disabling Security Hub product subscription %s: %s", d.Id(), err)
	}

	return nil
}
