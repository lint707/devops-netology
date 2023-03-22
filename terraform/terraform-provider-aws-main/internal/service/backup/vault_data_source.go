package backup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceVault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVaultRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_points": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceVaultRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get("name").(string)
	input := &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(name),
	}

	resp, err := conn.DescribeBackupVault(input)
	if err != nil {
		return fmt.Errorf("Error getting Backup Vault: %w", err)
	}

	d.SetId(aws.StringValue(resp.BackupVaultName))
	d.Set("arn", resp.BackupVaultArn)
	d.Set("kms_key_arn", resp.EncryptionKeyArn)
	d.Set("name", resp.BackupVaultName)
	d.Set("recovery_points", resp.NumberOfRecoveryPoints)

	tags, err := ListTags(conn, aws.StringValue(resp.BackupVaultArn))
	if err != nil {
		return fmt.Errorf("error listing tags for Backup Vault (%s): %w", name, err)
	}
	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
