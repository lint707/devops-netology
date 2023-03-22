package glue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceDataCatalogEncryptionSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataCatalogEncryptionSettingsRead,
		Schema: map[string]*schema.Schema{
			"catalog_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_catalog_encryption_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_password_encryption": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aws_kms_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"return_connection_password_encrypted": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"encryption_at_rest": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"catalog_encryption_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sse_aws_kms_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDataCatalogEncryptionSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GlueConn

	catalogID := d.Get("catalog_id").(string)
	output, err := conn.GetDataCatalogEncryptionSettings(&glue.GetDataCatalogEncryptionSettingsInput{
		CatalogId: aws.String(catalogID),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Glue Data Catalog Encryption Settings (%s): %w", catalogID, err))
	}

	d.SetId(catalogID)
	d.Set("catalog_id", d.Id())
	if output.DataCatalogEncryptionSettings != nil {
		if err := d.Set("data_catalog_encryption_settings", []interface{}{flattenDataCatalogEncryptionSettings(output.DataCatalogEncryptionSettings)}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data_catalog_encryption_settings: %w", err))
		}
	} else {
		d.Set("data_catalog_encryption_settings", nil)
	}

	return nil
}
