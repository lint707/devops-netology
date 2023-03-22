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

func ResourcePrincipalPortfolioAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrincipalPortfolioAssociationCreate,
		Read:   resourcePrincipalPortfolioAssociationRead,
		Delete: resourcePrincipalPortfolioAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(PrincipalPortfolioAssociationReadyTimeout),
			Read:   schema.DefaultTimeout(PrincipalPortfolioAssociationReadTimeout),
			Delete: schema.DefaultTimeout(PrincipalPortfolioAssociationDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"portfolio_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      servicecatalog.PrincipalTypeIam,
				ValidateFunc: validation.StringInSlice(servicecatalog.PrincipalType_Values(), false),
			},
		},
	}
}

func resourcePrincipalPortfolioAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	input := &servicecatalog.AssociatePrincipalWithPortfolioInput{
		PortfolioId:  aws.String(d.Get("portfolio_id").(string)),
		PrincipalARN: aws.String(d.Get("principal_arn").(string)),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("principal_type"); ok {
		input.PrincipalType = aws.String(v.(string))
	}

	var output *servicecatalog.AssociatePrincipalWithPortfolioOutput
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		output, err = conn.AssociatePrincipalWithPortfolio(input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.AssociatePrincipalWithPortfolio(input)
	}

	if err != nil {
		return fmt.Errorf("error associating Service Catalog Principal with Portfolio: %w", err)
	}

	if output == nil {
		return fmt.Errorf("error creating Service Catalog Principal Portfolio Association: empty response")
	}

	d.SetId(PrincipalPortfolioAssociationID(d.Get("accept_language").(string), d.Get("principal_arn").(string), d.Get("portfolio_id").(string)))

	return resourcePrincipalPortfolioAssociationRead(d, meta)
}

func resourcePrincipalPortfolioAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	acceptLanguage, principalARN, portfolioID, err := PrincipalPortfolioAssociationParseID(d.Id())

	if err != nil {
		return fmt.Errorf("could not parse ID (%s): %w", d.Id(), err)
	}

	if acceptLanguage == "" {
		acceptLanguage = AcceptLanguageEnglish
	}

	output, err := WaitPrincipalPortfolioAssociationReady(conn, acceptLanguage, principalARN, portfolioID, d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && (tfresource.NotFound(err) || tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException)) {
		log.Printf("[WARN] Service Catalog Principal Portfolio Association (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Service Catalog Principal Portfolio Association (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting Service Catalog Principal Portfolio Association (%s): empty response", d.Id())
	}

	d.Set("accept_language", acceptLanguage)
	d.Set("portfolio_id", portfolioID)
	d.Set("principal_arn", output.PrincipalARN)
	d.Set("principal_type", output.PrincipalType)

	return nil
}

func resourcePrincipalPortfolioAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	acceptLanguage, principalARN, portfolioID, err := PrincipalPortfolioAssociationParseID(d.Id())

	if err != nil {
		return fmt.Errorf("could not parse ID (%s): %w", d.Id(), err)
	}

	if acceptLanguage == "" {
		acceptLanguage = AcceptLanguageEnglish
	}

	input := &servicecatalog.DisassociatePrincipalFromPortfolioInput{
		PortfolioId:    aws.String(portfolioID),
		PrincipalARN:   aws.String(principalARN),
		AcceptLanguage: aws.String(acceptLanguage),
	}

	_, err = conn.DisassociatePrincipalFromPortfolio(input)

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error disassociating Service Catalog Principal from Portfolio (%s): %w", d.Id(), err)
	}

	err = WaitPrincipalPortfolioAssociationDeleted(conn, acceptLanguage, principalARN, portfolioID, d.Timeout(schema.TimeoutDelete))

	if tfresource.NotFound(err) || tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error waiting for Service Catalog Principal Portfolio Disassociation (%s): %w", d.Id(), err)
	}

	return nil
}
