package kms

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func DataSourcePublicKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePublicKeyRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_master_key_spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_algorithms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"grant_tokens": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validKey,
			},
			"key_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_algorithms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePublicKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KMSConn
	keyId := d.Get("key_id").(string)

	input := &kms.GetPublicKeyInput{
		KeyId: aws.String(keyId),
	}

	if v, ok := d.GetOk("grant_tokens"); ok {
		input.GrantTokens = aws.StringSlice(v.([]string))
	}

	output, err := conn.GetPublicKey(input)

	if err != nil {
		return fmt.Errorf("error while describing KMS public key (%s): %w", keyId, err)
	}

	d.SetId(aws.StringValue(output.KeyId))

	d.Set("arn", output.KeyId)
	d.Set("customer_master_key_spec", output.CustomerMasterKeySpec)
	d.Set("key_usage", output.KeyUsage)
	d.Set("public_key", base64.StdEncoding.EncodeToString(output.PublicKey))
	d.Set("public_key_pem", string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: output.PublicKey,
	})))

	if err := d.Set("encryption_algorithms", flex.FlattenStringList(output.EncryptionAlgorithms)); err != nil {
		return fmt.Errorf("error setting encryption_algorithms: %w", err)
	}

	if err := d.Set("signing_algorithms", flex.FlattenStringList(output.SigningAlgorithms)); err != nil {
		return fmt.Errorf("error setting signing_algorithms: %w", err)
	}

	return nil
}
