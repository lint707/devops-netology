package lambda

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	arn2 "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

const mutexLayerKey = `aws_lambda_layer_version`

func ResourceLayerVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayerVersionPublish,
		Read:   resourceLayerVersionRead,
		Delete: resourceLayerVersionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compatible_architectures": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(lambda.Architecture_Values(), false),
				},
			},
			"compatible_runtimes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MinItems: 0,
				MaxItems: 15,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(lambda.Runtime_Values(), false),
				},
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"filename": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"s3_bucket", "s3_key", "s3_object_version"},
			},
			"layer_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"license_info": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			"s3_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"filename"},
			},
			"s3_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"filename"},
			},
			"s3_object_version": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"filename"},
			},
			"signing_job_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_profile_version_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skip_destroy": {
				Type:     schema.TypeBool,
				Default:  false,
				ForceNew: true,
				Optional: true,
			},
			"source_code_hash": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"source_code_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayerVersionPublish(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LambdaConn

	layerName := d.Get("layer_name").(string)
	filename, hasFilename := d.GetOk("filename")
	s3Bucket, bucketOk := d.GetOk("s3_bucket")
	s3Key, keyOk := d.GetOk("s3_key")
	s3ObjectVersion, versionOk := d.GetOk("s3_object_version")

	if !hasFilename && !bucketOk && !keyOk && !versionOk {
		return errors.New("filename or s3_* attributes must be set")
	}

	var layerContent *lambda.LayerVersionContentInput
	if hasFilename {
		conns.GlobalMutexKV.Lock(mutexLayerKey)
		defer conns.GlobalMutexKV.Unlock(mutexLayerKey)
		file, err := loadFileContent(filename.(string))
		if err != nil {
			return fmt.Errorf("Unable to load %q: %s", filename.(string), err)
		}
		layerContent = &lambda.LayerVersionContentInput{
			ZipFile: file,
		}
	} else {
		if !bucketOk || !keyOk {
			return errors.New("s3_bucket and s3_key must all be set while using s3 code source")
		}
		layerContent = &lambda.LayerVersionContentInput{
			S3Bucket: aws.String(s3Bucket.(string)),
			S3Key:    aws.String(s3Key.(string)),
		}
		if versionOk {
			layerContent.S3ObjectVersion = aws.String(s3ObjectVersion.(string))
		}
	}

	params := &lambda.PublishLayerVersionInput{
		Content:     layerContent,
		Description: aws.String(d.Get("description").(string)),
		LayerName:   aws.String(layerName),
		LicenseInfo: aws.String(d.Get("license_info").(string)),
	}

	if v, ok := d.GetOk("compatible_runtimes"); ok && v.(*schema.Set).Len() > 0 {
		params.CompatibleRuntimes = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("compatible_architectures"); ok && v.(*schema.Set).Len() > 0 {
		params.CompatibleArchitectures = flex.ExpandStringSet(v.(*schema.Set))
	}

	log.Printf("[DEBUG] Publishing Lambda layer: %s", params)
	result, err := conn.PublishLayerVersion(params)
	if err != nil {
		return fmt.Errorf("Error creating lambda layer: %s", err)
	}

	d.SetId(aws.StringValue(result.LayerVersionArn))
	return resourceLayerVersionRead(d, meta)
}

func resourceLayerVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LambdaConn

	layerName, version, err := LayerVersionParseID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing lambda layer ID: %s", err)
	}

	layerVersion, err := conn.GetLayerVersion(&lambda.GetLayerVersionInput{
		LayerName:     aws.String(layerName),
		VersionNumber: aws.Int64(version),
	})

	if tfawserr.ErrCodeEquals(err, lambda.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Lambda Layer Version (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Lambda Layer version (%s): %s", d.Id(), err)
	}

	if err := d.Set("layer_name", layerName); err != nil {
		return fmt.Errorf("Error setting lambda layer name: %s", err)
	}
	if err := d.Set("version", strconv.FormatInt(version, 10)); err != nil {
		return fmt.Errorf("Error setting lambda layer version: %s", err)
	}
	if err := d.Set("arn", layerVersion.LayerVersionArn); err != nil {
		return fmt.Errorf("Error setting lambda layer version arn: %s", err)
	}
	if err := d.Set("layer_arn", layerVersion.LayerArn); err != nil {
		return fmt.Errorf("Error setting lambda layer arn: %s", err)
	}
	if err := d.Set("description", layerVersion.Description); err != nil {
		return fmt.Errorf("Error setting lambda layer description: %s", err)
	}
	if err := d.Set("license_info", layerVersion.LicenseInfo); err != nil {
		return fmt.Errorf("Error setting lambda layer license info: %s", err)
	}
	if err := d.Set("created_date", layerVersion.CreatedDate); err != nil {
		return fmt.Errorf("Error setting lambda layer created date: %s", err)
	}
	if err := d.Set("source_code_hash", layerVersion.Content.CodeSha256); err != nil {
		return fmt.Errorf("Error setting lambda layer source code hash: %s", err)
	}
	if err := d.Set("signing_profile_version_arn", layerVersion.Content.SigningProfileVersionArn); err != nil {
		return fmt.Errorf("Error setting lambda layer signing profile arn: %s", err)
	}
	if err := d.Set("signing_job_arn", layerVersion.Content.SigningJobArn); err != nil {
		return fmt.Errorf("Error setting lambda layer signing job arn: %s", err)
	}
	if err := d.Set("source_code_size", layerVersion.Content.CodeSize); err != nil {
		return fmt.Errorf("Error setting lambda layer source code size: %s", err)
	}
	if err := d.Set("compatible_runtimes", flex.FlattenStringList(layerVersion.CompatibleRuntimes)); err != nil {
		return fmt.Errorf("Error setting lambda layer compatible runtimes: %s", err)
	}

	if err := d.Set("compatible_architectures", flex.FlattenStringList(layerVersion.CompatibleArchitectures)); err != nil {
		return fmt.Errorf("Error setting lambda layer compatible architectures: %s", err)
	}

	return nil
}

func resourceLayerVersionDelete(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("skip_destroy"); ok && v.(bool) {
		log.Printf("[DEBUG] Retaining Lambda Layer Version %q", d.Id())
		return nil
	}

	conn := meta.(*conns.AWSClient).LambdaConn

	version, err := strconv.ParseInt(d.Get("version").(string), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing lambda layer version: %s", err)
	}

	_, err = conn.DeleteLayerVersion(&lambda.DeleteLayerVersionInput{
		LayerName:     aws.String(d.Get("layer_name").(string)),
		VersionNumber: aws.Int64(version),
	})
	if err != nil {
		return fmt.Errorf("Error deleting Lambda Layer Version (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Lambda layer %q deleted", d.Get("arn").(string))
	return nil
}

func LayerVersionParseID(id string) (layerName string, version int64, err error) {
	arn, err := arn2.Parse(id)
	if err != nil {
		return
	}
	parts := strings.Split(arn.Resource, ":")
	if len(parts) != 3 || parts[0] != "layer" {
		err = fmt.Errorf("lambda_layer ID must be a valid Layer ARN")
		return
	}

	layerName = parts[1]
	version, err = strconv.ParseInt(parts[2], 10, 64)
	return
}
