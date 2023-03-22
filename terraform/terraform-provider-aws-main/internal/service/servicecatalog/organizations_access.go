package servicecatalog

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceOrganizationsAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationsAccessCreate,
		Read:   resourceOrganizationsAccessRead,
		Delete: resourceOrganizationsAccessDelete,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(OrganizationsAccessStableTimeout),
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOrganizationsAccessCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	d.SetId(meta.(*conns.AWSClient).AccountID)

	// During create, if enabled = "true", then Enable Access and vice versa
	// During delete, the opposite

	if _, ok := d.GetOk("enabled"); ok {
		_, err := conn.EnableAWSOrganizationsAccess(&servicecatalog.EnableAWSOrganizationsAccessInput{})

		if err != nil {
			return fmt.Errorf("error enabling Service Catalog AWS Organizations Access: %w", err)
		}

		return resourceOrganizationsAccessRead(d, meta)
	}

	_, err := conn.DisableAWSOrganizationsAccess(&servicecatalog.DisableAWSOrganizationsAccessInput{})

	if err != nil {
		return fmt.Errorf("error disabling Service Catalog AWS Organizations Access: %w", err)
	}

	return resourceOrganizationsAccessRead(d, meta)
}

func resourceOrganizationsAccessRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	output, err := WaitOrganizationsAccessStable(conn, d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		// theoretically this should not be possible
		log.Printf("[WARN] Service Catalog Organizations Access (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Service Catalog AWS Organizations Access (%s): %w", d.Id(), err)
	}

	if output == "" {
		return fmt.Errorf("error getting Service Catalog AWS Organizations Access (%s): empty response", d.Id())
	}

	if output == servicecatalog.AccessStatusEnabled {
		d.Set("enabled", true)
		return nil
	}

	d.Set("enabled", false)
	return nil
}

func resourceOrganizationsAccessDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	// During create, if enabled = "true", then Enable Access and vice versa
	// During delete, the opposite

	if _, ok := d.GetOk("enabled"); !ok {
		_, err := conn.EnableAWSOrganizationsAccess(&servicecatalog.EnableAWSOrganizationsAccessInput{})

		if err != nil {
			return fmt.Errorf("error enabling Service Catalog AWS Organizations Access: %w", err)
		}

		return nil
	}

	_, err := conn.DisableAWSOrganizationsAccess(&servicecatalog.DisableAWSOrganizationsAccessInput{})

	if err != nil {
		return fmt.Errorf("error disabling Service Catalog AWS Organizations Access: %w", err)
	}

	return nil
}
