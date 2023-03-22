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
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceConstraint() *schema.Resource {
	return &schema.Resource{
		Create: resourceConstraintCreate,
		Read:   resourceConstraintRead,
		Update: resourceConstraintUpdate,
		Delete: resourceConstraintDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ConstraintReadyTimeout),
			Read:   schema.DefaultTimeout(ConstraintReadTimeout),
			Update: schema.DefaultTimeout(ConstraintUpdateTimeout),
			Delete: schema.DefaultTimeout(ConstraintDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
			},
			"portfolio_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ConstraintType_Values(), false),
			},
		},
	}
}

func resourceConstraintCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.CreateConstraintInput{
		IdempotencyToken: aws.String(resource.UniqueId()),
		Parameters:       aws.String(d.Get("parameters").(string)),
		PortfolioId:      aws.String(d.Get("portfolio_id").(string)),
		ProductId:        aws.String(d.Get("product_id").(string)),
		Type:             aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	var output *servicecatalog.CreateConstraintOutput
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		output, err = conn.CreateConstraint(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateConstraint(input)
	}

	if err != nil {
		return fmt.Errorf("error creating Service Catalog Constraint: %w", err)
	}

	if output == nil || output.ConstraintDetail == nil {
		return fmt.Errorf("error creating Service Catalog Constraint: empty response")
	}

	d.SetId(aws.StringValue(output.ConstraintDetail.ConstraintId))

	return resourceConstraintRead(d, meta)
}

func resourceConstraintRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	output, err := WaitConstraintReady(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Service Catalog Constraint (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Service Catalog Constraint (%s): %w", d.Id(), err)
	}

	if output == nil || output.ConstraintDetail == nil {
		return fmt.Errorf("error getting Service Catalog Constraint (%s): empty response", d.Id())
	}

	acceptLanguage := d.Get("accept_language").(string)

	if acceptLanguage == "" {
		acceptLanguage = AcceptLanguageEnglish
	}

	d.Set("accept_language", acceptLanguage)

	d.Set("parameters", output.ConstraintParameters)
	d.Set("status", output.Status)

	detail := output.ConstraintDetail

	d.Set("description", detail.Description)
	d.Set("owner", detail.Owner)
	d.Set("portfolio_id", detail.PortfolioId)
	d.Set("product_id", detail.ProductId)
	d.Set("type", detail.Type)

	return nil
}

func resourceConstraintUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.UpdateConstraintInput{
		Id: aws.String(d.Id()),
	}

	if d.HasChange("accept_language") {
		input.AcceptLanguage = aws.String(d.Get("accept_language").(string))
	}

	if d.HasChange("description") {
		input.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChange("parameters") {
		input.Parameters = aws.String(d.Get("parameters").(string))
	}

	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := conn.UpdateConstraint(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.UpdateConstraint(input)
	}

	if err != nil {
		return fmt.Errorf("error updating Service Catalog Constraint (%s): %w", d.Id(), err)
	}

	return resourceConstraintRead(d, meta)
}

func resourceConstraintDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.DeleteConstraintInput{
		Id: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	_, err := conn.DeleteConstraint(input)

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Service Catalog Constraint (%s): %w", d.Id(), err)
	}

	err = WaitConstraintDeleted(conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutDelete))

	if err != nil && !tfresource.NotFound(err) {
		return fmt.Errorf("error waiting for Service Catalog Constraint (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}
