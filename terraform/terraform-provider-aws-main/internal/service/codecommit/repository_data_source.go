package codecommit

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"repository_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"repository_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"clone_url_http": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"clone_url_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CodeCommitConn

	repositoryName := d.Get("repository_name").(string)
	input := &codecommit.GetRepositoryInput{
		RepositoryName: aws.String(repositoryName),
	}

	out, err := conn.GetRepository(input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, codecommit.ErrCodeRepositoryDoesNotExistException) {
			log.Printf("[WARN] CodeCommit Repository (%s) not found, removing from state", d.Id())
			d.SetId("")
			return fmt.Errorf("Resource codecommit repository not found for %s", repositoryName)
		} else {
			return fmt.Errorf("Error reading CodeCommit Repository: %w", err)
		}
	}

	if out.RepositoryMetadata == nil {
		return fmt.Errorf("no matches found for repository name: %s", repositoryName)
	}

	d.SetId(aws.StringValue(out.RepositoryMetadata.RepositoryName))
	d.Set("arn", out.RepositoryMetadata.Arn)
	d.Set("clone_url_http", out.RepositoryMetadata.CloneUrlHttp)
	d.Set("clone_url_ssh", out.RepositoryMetadata.CloneUrlSsh)
	d.Set("repository_name", out.RepositoryMetadata.RepositoryName)
	d.Set("repository_id", out.RepositoryMetadata.RepositoryId)

	return nil
}
