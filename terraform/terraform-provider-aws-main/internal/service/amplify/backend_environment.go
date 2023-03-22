package amplify

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceBackendEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackendEnvironmentCreate,
		Read:   resourceBackendEnvironmentRead,
		Delete: resourceBackendEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"deployment_artifacts": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z-]{1,100}$`), "should be not be more than 100 alphanumeric or hyphen characters"),
			},

			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z]{2,10}$`), "should be between 2 and 10 characters (only lowercase alphabetic)"),
			},

			"stack_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z-]{1,100}$`), "should be not be more than 100 alphanumeric or hyphen characters"),
			},
		},
	}
}

func resourceBackendEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	appID := d.Get("app_id").(string)
	environmentName := d.Get("environment_name").(string)
	id := BackendEnvironmentCreateResourceID(appID, environmentName)

	input := &amplify.CreateBackendEnvironmentInput{
		AppId:           aws.String(appID),
		EnvironmentName: aws.String(environmentName),
	}

	if v, ok := d.GetOk("deployment_artifacts"); ok {
		input.DeploymentArtifacts = aws.String(v.(string))
	}

	if v, ok := d.GetOk("stack_name"); ok {
		input.StackName = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Amplify Backend Environment: %s", input)
	_, err := conn.CreateBackendEnvironment(input)

	if err != nil {
		return fmt.Errorf("error creating Amplify Backend Environment (%s): %w", id, err)
	}

	d.SetId(id)

	return resourceBackendEnvironmentRead(d, meta)
}

func resourceBackendEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	appID, environmentName, err := BackendEnvironmentParseResourceID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing Amplify Backend Environment ID: %w", err)
	}

	backendEnvironment, err := FindBackendEnvironmentByAppIDAndEnvironmentName(conn, appID, environmentName)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Amplify Backend Environment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Amplify Backend Environment (%s): %w", d.Id(), err)
	}

	d.Set("app_id", appID)
	d.Set("arn", backendEnvironment.BackendEnvironmentArn)
	d.Set("deployment_artifacts", backendEnvironment.DeploymentArtifacts)
	d.Set("environment_name", backendEnvironment.EnvironmentName)
	d.Set("stack_name", backendEnvironment.StackName)

	return nil
}

func resourceBackendEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AmplifyConn

	appID, environmentName, err := BackendEnvironmentParseResourceID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing Amplify Backend Environment ID: %w", err)
	}

	log.Printf("[DEBUG] Deleting Amplify Backend Environment: %s", d.Id())
	_, err = conn.DeleteBackendEnvironment(&amplify.DeleteBackendEnvironmentInput{
		AppId:           aws.String(appID),
		EnvironmentName: aws.String(environmentName),
	})

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Amplify Backend Environment (%s): %w", d.Id(), err)
	}

	return nil
}
