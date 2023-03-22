package lambda

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func DataSourceLayerVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayerVersionRead,

		Schema: map[string]*schema.Schema{
			"layer_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"compatible_runtimes"},
			},
			"compatible_runtime": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(lambda.Runtime_Values(), false),
				ConflictsWith: []string{"version"},
			},
			"compatible_runtimes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layer_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_code_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_code_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"signing_profile_version_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_job_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compatible_architecture": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(lambda.Architecture_Values(), false),
				ConflictsWith: []string{"version"},
			},
			"compatible_architectures": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceLayerVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LambdaConn
	layerName := d.Get("layer_name").(string)

	var version int64

	if v, ok := d.GetOk("version"); ok {
		version = int64(v.(int))
	} else {
		listInput := &lambda.ListLayerVersionsInput{
			LayerName: aws.String(layerName),
		}
		if v, ok := d.GetOk("compatible_runtime"); ok {
			listInput.CompatibleRuntime = aws.String(v.(string))
		}

		if v, ok := d.GetOk("compatible_architecture"); ok {
			listInput.CompatibleArchitecture = aws.String(v.(string))
		}

		log.Printf("[DEBUG] Looking up latest version for lambda layer %s", layerName)
		listOutput, err := conn.ListLayerVersions(listInput)
		if err != nil {
			return fmt.Errorf("error listing Lambda Layer Versions (%s): %w", layerName, err)
		}

		if len(listOutput.LayerVersions) == 0 {
			return fmt.Errorf("error listing Lambda Layer Versions (%s): empty response", layerName)
		}

		version = aws.Int64Value(listOutput.LayerVersions[0].Version)
	}

	input := &lambda.GetLayerVersionInput{
		LayerName:     aws.String(layerName),
		VersionNumber: aws.Int64(version),
	}

	log.Printf("[DEBUG] Getting Lambda Layer Version: %s, version %d", layerName, version)
	output, err := conn.GetLayerVersion(input)

	if err != nil {
		return fmt.Errorf("error getting Lambda Layer Version (%s, version %d): %w", layerName, version, err)
	}

	if output == nil {
		return fmt.Errorf("error getting Lambda Layer Version (%s, version %d): empty response", layerName, version)
	}

	if err := d.Set("version", output.Version); err != nil {
		return fmt.Errorf("error setting lambda layer version: %w", err)
	}
	if err := d.Set("compatible_runtimes", flex.FlattenStringList(output.CompatibleRuntimes)); err != nil {
		return fmt.Errorf("error setting lambda layer compatible runtimes: %w", err)
	}
	if err := d.Set("compatible_architectures", flex.FlattenStringList(output.CompatibleArchitectures)); err != nil {
		return fmt.Errorf("Error setting lambda layer compatible architectures: %w", err)
	}
	if err := d.Set("description", output.Description); err != nil {
		return fmt.Errorf("error setting lambda layer description: %w", err)
	}
	if err := d.Set("license_info", output.LicenseInfo); err != nil {
		return fmt.Errorf("error setting lambda layer license info: %w", err)
	}
	if err := d.Set("arn", output.LayerVersionArn); err != nil {
		return fmt.Errorf("error setting lambda layer version arn: %w", err)
	}
	if err := d.Set("layer_arn", output.LayerArn); err != nil {
		return fmt.Errorf("error setting lambda layer arn: %w", err)
	}
	if err := d.Set("created_date", output.CreatedDate); err != nil {
		return fmt.Errorf("error setting lambda layer created date: %w", err)
	}
	if err := d.Set("source_code_hash", output.Content.CodeSha256); err != nil {
		return fmt.Errorf("error setting lambda layer source code hash: %w", err)
	}
	if err := d.Set("source_code_size", output.Content.CodeSize); err != nil {
		return fmt.Errorf("error setting lambda layer source code size: %w", err)
	}
	if err := d.Set("signing_profile_version_arn", output.Content.SigningProfileVersionArn); err != nil {
		return fmt.Errorf("Error setting lambda layer signing profile arn: %w", err)
	}
	if err := d.Set("signing_job_arn", output.Content.SigningJobArn); err != nil {
		return fmt.Errorf("Error setting lambda layer signing job arn: %w", err)
	}

	d.SetId(aws.StringValue(output.LayerVersionArn))

	return nil
}
