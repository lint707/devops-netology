package route53resolver

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validResolverName,
			},

			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(route53resolver.Action_Values(), false),
			},

			"block_override_dns_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(route53resolver.BlockOverrideDnsType_Values(), false),
			},

			"block_override_domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},

			"block_override_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 604800),
			},

			"block_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(route53resolver.BlockResponse_Values(), false),
			},

			"firewall_domain_list_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"firewall_rule_group_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	firewallRuleGroupId := d.Get("firewall_rule_group_id").(string)
	firewallDomainListId := d.Get("firewall_domain_list_id").(string)
	input := &route53resolver.CreateFirewallRuleInput{
		CreatorRequestId:     aws.String(resource.PrefixedUniqueId("tf-r53-resolver-firewall-rule-")),
		Name:                 aws.String(d.Get("name").(string)),
		Action:               aws.String(d.Get("action").(string)),
		FirewallRuleGroupId:  aws.String(firewallRuleGroupId),
		FirewallDomainListId: aws.String(firewallDomainListId),
		Priority:             aws.Int64(int64(d.Get("priority").(int))),
	}

	if v, ok := d.GetOk("block_override_dns_type"); ok {
		input.BlockOverrideDnsType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("block_override_domain"); ok {
		input.BlockOverrideDomain = aws.String(v.(string))
	}

	if v, ok := d.GetOk("block_override_ttl"); ok {
		input.BlockOverrideTtl = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("block_response"); ok {
		input.BlockResponse = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Route 53 Resolver DNS Firewall rule: %#v", input)
	_, err := conn.CreateFirewallRule(input)
	if err != nil {
		return fmt.Errorf("error creating Route 53 Resolver DNS Firewall rule: %w", err)
	}

	d.SetId(FirewallRuleCreateID(firewallRuleGroupId, firewallDomainListId))

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	rule, err := FindFirewallRuleByID(conn, d.Id())

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Route53 Resolver DNS Firewall rule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting Route 53 Resolver DNS Firewall rule (%s): %w", d.Id(), err)
	}

	if rule == nil {
		log.Printf("[WARN] Route 53 Resolver DNS Firewall rule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("action", rule.Action)
	d.Set("block_override_dns_type", rule.BlockOverrideDnsType)
	d.Set("block_override_domain", rule.BlockOverrideDomain)
	d.Set("block_override_ttl", rule.BlockOverrideTtl)
	d.Set("block_response", rule.BlockResponse)
	d.Set("firewall_rule_group_id", rule.FirewallRuleGroupId)
	d.Set("firewall_domain_list_id", rule.FirewallDomainListId)
	d.Set("priority", rule.Priority)

	return nil
}

func resourceFirewallRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	input := &route53resolver.UpdateFirewallRuleInput{
		Name:                 aws.String(d.Get("name").(string)),
		Action:               aws.String(d.Get("action").(string)),
		FirewallRuleGroupId:  aws.String(d.Get("firewall_rule_group_id").(string)),
		FirewallDomainListId: aws.String(d.Get("firewall_domain_list_id").(string)),
		Priority:             aws.Int64(int64(d.Get("priority").(int))),
	}

	if v, ok := d.GetOk("block_override_dns_type"); ok {
		input.BlockOverrideDnsType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("block_override_domain"); ok {
		input.BlockOverrideDomain = aws.String(v.(string))
	}

	if v, ok := d.GetOk("block_override_ttl"); ok {
		input.BlockOverrideTtl = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("block_response"); ok {
		input.BlockResponse = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Updating Route 53 Resolver DNS Firewall rule: %#v", input)
	_, err := conn.UpdateFirewallRule(input)
	if err != nil {
		return fmt.Errorf("error updating Route 53 Resolver DNS Firewall rule (%s): %w", d.Id(), err)
	}

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53ResolverConn

	_, err := conn.DeleteFirewallRule(&route53resolver.DeleteFirewallRuleInput{
		FirewallRuleGroupId:  aws.String(d.Get("firewall_rule_group_id").(string)),
		FirewallDomainListId: aws.String(d.Get("firewall_domain_list_id").(string)),
	})

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Route 53 Resolver DNS Firewall rule (%s): %w", d.Id(), err)
	}

	return nil
}
