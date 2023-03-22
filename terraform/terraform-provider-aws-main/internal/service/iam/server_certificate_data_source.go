package iam

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceServerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerCertificateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validation.StringLenBetween(0, 128),
			},

			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(0, 128-resource.UniqueIDSuffixLength),
			},

			"path_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"latest": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"upload_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_body": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_chain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type CertificateByExpiration []*iam.ServerCertificateMetadata

func (m CertificateByExpiration) Len() int {
	return len(m)
}

func (m CertificateByExpiration) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m CertificateByExpiration) Less(i, j int) bool {
	return m[i].Expiration.After(*m[j].Expiration)
}

func dataSourceServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	var matcher = func(cert *iam.ServerCertificateMetadata) bool {
		return strings.HasPrefix(aws.StringValue(cert.ServerCertificateName), d.Get("name_prefix").(string))
	}
	if v, ok := d.GetOk("name"); ok {
		matcher = func(cert *iam.ServerCertificateMetadata) bool {
			return aws.StringValue(cert.ServerCertificateName) == v.(string)
		}
	}

	var metadatas []*iam.ServerCertificateMetadata
	input := &iam.ListServerCertificatesInput{}
	if v, ok := d.GetOk("path_prefix"); ok {
		input.PathPrefix = aws.String(v.(string))
	}
	log.Printf("[DEBUG] Reading IAM Server Certificate")
	err := conn.ListServerCertificatesPages(input, func(p *iam.ListServerCertificatesOutput, lastPage bool) bool {
		for _, cert := range p.ServerCertificateMetadataList {
			if matcher(cert) {
				metadatas = append(metadatas, cert)
			}
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("Error describing certificates: %w", err)
	}

	if len(metadatas) == 0 {
		return fmt.Errorf("Search for AWS IAM server certificate returned no results")
	}
	if len(metadatas) > 1 {
		if !d.Get("latest").(bool) {
			return fmt.Errorf("Search for AWS IAM server certificate returned too many results")
		}

		sort.Sort(CertificateByExpiration(metadatas))
	}

	metadata := metadatas[0]
	d.SetId(aws.StringValue(metadata.ServerCertificateId))
	d.Set("arn", metadata.Arn)
	d.Set("path", metadata.Path)
	d.Set("name", metadata.ServerCertificateName)
	if metadata.Expiration != nil {
		d.Set("expiration_date", metadata.Expiration.Format(time.RFC3339))
	}

	log.Printf("[DEBUG] Get Public Key Certificate for %s", *metadata.ServerCertificateName)
	serverCertificateResp, err := conn.GetServerCertificate(&iam.GetServerCertificateInput{
		ServerCertificateName: metadata.ServerCertificateName,
	})
	if err != nil {
		return err
	}
	d.Set("upload_date", serverCertificateResp.ServerCertificate.ServerCertificateMetadata.UploadDate.Format(time.RFC3339))
	d.Set("certificate_body", serverCertificateResp.ServerCertificate.CertificateBody)
	d.Set("certificate_chain", serverCertificateResp.ServerCertificate.CertificateChain)

	return nil
}
