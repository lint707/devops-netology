package glacier

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceVaultLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultLockCreate,
		Read:   resourceVaultLockRead,
		// Allow ignore_deletion_error update
		Update: schema.Noop,
		Delete: resourceVaultLockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"complete_lock": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"ignore_deletion_error": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				ValidateFunc:     verify.ValidIAMPolicyJSON,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			"vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceVaultLockCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GlacierConn
	vaultName := d.Get("vault_name").(string)

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
	}

	input := &glacier.InitiateVaultLockInput{
		AccountId: aws.String("-"),
		Policy: &glacier.VaultLockPolicy{
			Policy: aws.String(policy),
		},
		VaultName: aws.String(vaultName),
	}

	log.Printf("[DEBUG] Initiating Glacier Vault Lock: %s", input)
	output, err := conn.InitiateVaultLock(input)
	if err != nil {
		return fmt.Errorf("error initiating Glacier Vault Lock: %s", err)
	}

	d.SetId(vaultName)

	if !d.Get("complete_lock").(bool) {
		return resourceVaultLockRead(d, meta)
	}

	completeLockInput := &glacier.CompleteVaultLockInput{
		LockId:    output.LockId,
		VaultName: aws.String(vaultName),
	}

	log.Printf("[DEBUG] Completing Glacier Vault (%s) Lock: %s", vaultName, completeLockInput)
	if _, err := conn.CompleteVaultLock(completeLockInput); err != nil {
		return fmt.Errorf("error completing Glacier Vault (%s) Lock: %s", vaultName, err)
	}

	if err := waitVaultLockCompletion(conn, vaultName); err != nil {
		return fmt.Errorf("error waiting for Glacier Vault Lock (%s) completion: %s", d.Id(), err)
	}

	return resourceVaultLockRead(d, meta)
}

func resourceVaultLockRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GlacierConn

	input := &glacier.GetVaultLockInput{
		AccountId: aws.String("-"),
		VaultName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Reading Glacier Vault Lock (%s): %s", d.Id(), input)
	output, err := conn.GetVaultLock(input)

	if tfawserr.ErrCodeEquals(err, glacier.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Glacier Vault Lock (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Glacier Vault Lock (%s): %s", d.Id(), err)
	}

	if output == nil {
		log.Printf("[WARN] Glacier Vault Lock (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("complete_lock", aws.StringValue(output.State) == "Locked")
	d.Set("vault_name", d.Id())

	policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), aws.StringValue(output.Policy))

	if err != nil {
		return err
	}

	d.Set("policy", policyToSet)

	return nil
}

func resourceVaultLockDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GlacierConn

	input := &glacier.AbortVaultLockInput{
		VaultName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Aborting Glacier Vault Lock (%s): %s", d.Id(), input)
	_, err := conn.AbortVaultLock(input)

	if tfawserr.ErrCodeEquals(err, glacier.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil && !d.Get("ignore_deletion_error").(bool) {
		return fmt.Errorf("error aborting Glacier Vault Lock (%s): %s", d.Id(), err)
	}

	return nil
}

func vaultLockRefreshFunc(conn *glacier.Glacier, vaultName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &glacier.GetVaultLockInput{
			AccountId: aws.String("-"),
			VaultName: aws.String(vaultName),
		}

		log.Printf("[DEBUG] Reading Glacier Vault Lock (%s): %s", vaultName, input)
		output, err := conn.GetVaultLock(input)

		if tfawserr.ErrCodeEquals(err, glacier.ErrCodeResourceNotFoundException) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("error reading Glacier Vault Lock (%s): %s", vaultName, err)
		}

		if output == nil {
			return nil, "", nil
		}

		return output, aws.StringValue(output.State), nil
	}
}

func waitVaultLockCompletion(conn *glacier.Glacier, vaultName string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"InProgress"},
		Target:  []string{"Locked"},
		Refresh: vaultLockRefreshFunc(conn, vaultName),
		Timeout: 5 * time.Minute,
	}

	log.Printf("[DEBUG] Waiting for Glacier Vault Lock (%s) completion", vaultName)
	_, err := stateConf.WaitForState()

	return err
}
