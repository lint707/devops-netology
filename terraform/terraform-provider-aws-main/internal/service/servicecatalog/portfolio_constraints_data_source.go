package servicecatalog

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourcePortfolioConstraints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePortfolioConstraintsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(PortfolioConstraintsReadyTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"constraint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portfolio_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"portfolio_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePortfolioConstraintsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ServiceCatalogConn

	output, err := WaitPortfolioConstraintsReady(conn, d.Get("accept_language").(string), d.Get("portfolio_id").(string), d.Get("product_id").(string), d.Timeout(schema.TimeoutRead))

	if err != nil {
		return fmt.Errorf("error describing Service Catalog Portfolio Constraints: %w", err)
	}

	if len(output) == 0 {
		return fmt.Errorf("error getting Service Catalog Portfolio Constraints: no results, change your input")
	}

	acceptLanguage := d.Get("accept_language").(string)

	if acceptLanguage == "" {
		acceptLanguage = AcceptLanguageEnglish
	}

	d.Set("accept_language", acceptLanguage)
	d.Set("portfolio_id", d.Get("portfolio_id").(string))
	d.Set("product_id", d.Get("product_id").(string))

	if err := d.Set("details", flattenConstraintDetails(output)); err != nil {
		return fmt.Errorf("error setting details: %w", err)
	}

	d.SetId(PortfolioConstraintsID(d.Get("accept_language").(string), d.Get("portfolio_id").(string), d.Get("product_id").(string)))

	return nil
}

func flattenConstraintDetail(apiObject *servicecatalog.ConstraintDetail) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.ConstraintId; v != nil {
		tfMap["constraint_id"] = aws.StringValue(v)
	}

	if v := apiObject.Description; v != nil {
		tfMap["description"] = aws.StringValue(v)
	}

	if v := apiObject.Owner; v != nil {
		tfMap["owner"] = aws.StringValue(v)
	}

	if v := apiObject.PortfolioId; v != nil {
		tfMap["portfolio_id"] = aws.StringValue(v)
	}

	if v := apiObject.ProductId; v != nil {
		tfMap["product_id"] = aws.StringValue(v)
	}

	if v := apiObject.Type; v != nil {
		tfMap["type"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenConstraintDetails(apiObjects []*servicecatalog.ConstraintDetail) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenConstraintDetail(apiObject))
	}

	return tfList
}
