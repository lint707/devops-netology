package route53resolver

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceQueryLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceQueryLogConfigCreate,
		Read:   resourceQueryLogConfigRead,
		Update: resourceQueryLogConfigUpdate,
		Delete: resourceQueryLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"destination_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validResolverName,
			},

			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"share_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceQueryLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &route53resolver.CreateResolverQueryLogConfigInput{
		CreatorRequestId: aws.String(resource.PrefixedUniqueId("tf-r53-resolver-query-log-config-")),
		DestinationArn:   aws.String(d.Get("destination_arn").(string)),
		Name:             aws.String(d.Get("name").(string)),
	}
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating Route53 Resolver Query Log Config: %s", input)
	output, err := conn.CreateResolverQueryLogConfig(input)

	if err != nil {
		return fmt.Errorf("error creating Route53 Resolver Query Log Config: %w", err)
	}

	d.SetId(aws.StringValue(output.ResolverQueryLogConfig.Id))

	_, err = WaitQueryLogConfigCreated(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error waiting for Route53 Resolver Query Log Config (%s) to become available: %w", d.Id(), err)
	}

	return resourceQueryLogConfigRead(d, meta)
}

func resourceQueryLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	queryLogConfig, err := FindResolverQueryLogConfigByID(conn, d.Id())

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Route53 Resolver Query Log Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Route53 Resolver Query Log Config (%s): %w", d.Id(), err)
	}

	if queryLogConfig == nil {
		log.Printf("[WARN] Route53 Resolver Query Log Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	arn := aws.StringValue(queryLogConfig.Arn)
	d.Set("arn", arn)
	d.Set("destination_arn", queryLogConfig.DestinationArn)
	d.Set("name", queryLogConfig.Name)
	d.Set("owner_id", queryLogConfig.OwnerId)
	d.Set("share_status", queryLogConfig.ShareStatus)

	tags, err := ListTags(conn, arn)
	if err != nil {
		return fmt.Errorf("error listing tags for Route53 Resolver Query Log Config (%s): %w", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceQueryLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating Route53 Resolver Query Log Config (%s) tags: %s", d.Get("arn").(string), err)
		}
	}

	return resourceQueryLogConfigRead(d, meta)
}

func resourceQueryLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	log.Printf("[DEBUG] Deleting Route53 Resolver Query Log Config (%s)", d.Id())
	_, err := conn.DeleteResolverQueryLogConfig(&route53resolver.DeleteResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Route53 Resolver Query Log Config (%s): %w", d.Id(), err)
	}

	_, err = WaitQueryLogConfigDeleted(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error waiting for Route53 Resolver Query Log Config (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}
