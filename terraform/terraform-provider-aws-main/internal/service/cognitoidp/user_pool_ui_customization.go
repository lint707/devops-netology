package cognitoidp

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceUserPoolUICustomization() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserPoolUICustomizationPut,
		Read:   resourceUserPoolUICustomizationRead,
		Update: resourceUserPoolUICustomizationPut,
		Delete: resourceUserPoolUICustomizationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALL",
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"css": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"css", "image_file"},
			},

			"css_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_file": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_file", "css"},
			},

			"image_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserPoolUICustomizationPut(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CognitoIDPConn

	clientId := d.Get("client_id").(string)
	userPoolId := d.Get("user_pool_id").(string)

	input := &cognitoidentityprovider.SetUICustomizationInput{
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String(userPoolId),
	}

	if v, ok := d.GetOk("css"); ok {
		input.CSS = aws.String(v.(string))
	}

	if v, ok := d.GetOk("image_file"); ok {
		imgFile, err := base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			return fmt.Errorf("error Base64 decoding image file for Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %w", userPoolId, clientId, err)
		}

		input.ImageFile = imgFile
	}

	_, err := conn.SetUICustomization(input)

	if err != nil {
		return fmt.Errorf("error setting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %w", userPoolId, clientId, err)
	}

	d.SetId(fmt.Sprintf("%s,%s", userPoolId, clientId))

	return resourceUserPoolUICustomizationRead(d, meta)
}

func resourceUserPoolUICustomizationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CognitoIDPConn

	userPoolId, clientId, err := ParseUserPoolUICustomizationID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing Cognito User Pool UI customization ID (%s): %w", d.Id(), err)
	}

	uiCustomization, err := FindCognitoUserPoolUICustomization(conn, userPoolId, clientId)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, cognitoidentityprovider.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s) not found, removing from state", userPoolId, clientId)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %w", userPoolId, clientId, err)
	}

	if uiCustomization == nil {
		if d.IsNewResource() {
			return fmt.Errorf("error getting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): not found", userPoolId, clientId)
		}

		log.Printf("[WARN] Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s) not found, removing from state", userPoolId, clientId)
		d.SetId("")
		return nil
	}

	d.Set("client_id", uiCustomization.ClientId)
	d.Set("creation_date", aws.TimeValue(uiCustomization.CreationDate).Format(time.RFC3339))
	d.Set("css", uiCustomization.CSS)
	d.Set("css_version", uiCustomization.CSSVersion)
	d.Set("image_url", uiCustomization.ImageUrl)
	d.Set("last_modified_date", aws.TimeValue(uiCustomization.LastModifiedDate).Format(time.RFC3339))
	d.Set("user_pool_id", uiCustomization.UserPoolId)

	return nil
}

func resourceUserPoolUICustomizationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CognitoIDPConn

	userPoolId, clientId, err := ParseUserPoolUICustomizationID(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing Cognito User Pool UI customization ID (%s): %w", d.Id(), err)
	}

	input := &cognitoidentityprovider.SetUICustomizationInput{
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String(userPoolId),
	}

	_, err = conn.SetUICustomization(input)

	if tfawserr.ErrCodeEquals(err, cognitoidentityprovider.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %w", userPoolId, clientId, err)
	}

	return nil
}

func ParseUserPoolUICustomizationID(id string) (string, string, error) {
	idParts := strings.SplitN(id, ",", 2)

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("please make sure ID is in format USER_POOL_ID,CLIENT_ID")
	}

	return idParts[0], idParts[1], nil
}
