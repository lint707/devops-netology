package kms

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func ResourceCiphertext() *schema.Resource {

	return &schema.Resource{
		Create: resourceCiphertextCreate,
		Read:   schema.Noop,
		Delete: schema.Noop,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCiphertextCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KMSConn

	//lintignore:R017 // Allow legacy unstable ID usage in managed resource
	d.SetId(time.Now().UTC().String())

	req := &kms.EncryptInput{
		KeyId:     aws.String(d.Get("key_id").(string)),
		Plaintext: []byte(d.Get("plaintext").(string)),
	}

	if ec := d.Get("context"); ec != nil {
		req.EncryptionContext = flex.ExpandStringMap(ec.(map[string]interface{}))
	}

	log.Printf("[DEBUG] KMS encrypt for key: %s", d.Get("key_id").(string))
	resp, err := conn.Encrypt(req)
	if err != nil {
		return err
	}

	d.Set("ciphertext_blob", base64.StdEncoding.EncodeToString(resp.CiphertextBlob))

	return nil
}
