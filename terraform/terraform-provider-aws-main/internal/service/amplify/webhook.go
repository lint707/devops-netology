package amplify

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"branch_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z/_.-]{1,255}$`), "should be not be more than 255 letters, numbers, and the symbols /_.-"),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	input := &amplify.CreateWebhookInput{
		AppId:      aws.String(d.Get("app_id").(string)),
		BranchName: aws.String(d.Get("branch_name").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Amplify Webhook: %s", input)
	output, err := conn.CreateWebhook(input)

	if err != nil {
		return fmt.Errorf("error creating Amplify Webhook: %w", err)
	}

	d.SetId(aws.StringValue(output.Webhook.WebhookId))

	return resourceWebhookRead(d, meta)
}

func resourceWebhookRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	webhook, err := FindWebhookByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Amplify Webhook (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Amplify Webhook (%s): %w", d.Id(), err)
	}

	webhookArn := aws.StringValue(webhook.WebhookArn)
	arn, err := arn.Parse(webhookArn)

	if err != nil {
		return fmt.Errorf("error parsing %q: %w", webhookArn, err)
	}

	// arn:${Partition}:amplify:${Region}:${Account}:apps/${AppId}/webhooks/${WebhookId}
	parts := strings.Split(arn.Resource, "/")

	if len(parts) != 4 {
		return fmt.Errorf("unexpected format for ARN resource (%s)", arn.Resource)
	}

	d.Set("app_id", parts[1])
	d.Set("arn", webhookArn)
	d.Set("branch_name", webhook.BranchName)
	d.Set("description", webhook.Description)
	d.Set("url", webhook.WebhookUrl)

	return nil
}

func resourceWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	input := &amplify.UpdateWebhookInput{
		WebhookId: aws.String(d.Id()),
	}

	if d.HasChange("branch_name") {
		input.BranchName = aws.String(d.Get("branch_name").(string))
	}

	if d.HasChange("description") {
		input.Description = aws.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] Updating Amplify Webhook: %s", input)
	_, err := conn.UpdateWebhook(input)

	if err != nil {
		return fmt.Errorf("error updating Amplify Webhook (%s): %w", d.Id(), err)
	}

	return resourceWebhookRead(d, meta)
}

func resourceWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	log.Printf("[DEBUG] Deleting Amplify Webhook: %s", d.Id())
	_, err := conn.DeleteWebhook(&amplify.DeleteWebhookInput{
		WebhookId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Amplify Webhook (%s): %w", d.Id(), err)
	}

	return nil
}
