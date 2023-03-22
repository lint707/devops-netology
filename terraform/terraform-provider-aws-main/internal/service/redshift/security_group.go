package redshift

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringNotInSlice([]string{"default"}, false),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z-]+$`), "must contain only lowercase alphanumeric characters and hyphens"),
				),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
			},

			"ingress": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"security_group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"security_group_owner_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceSecurityGroupIngressHash,
			},
		},

		DeprecationMessage: `With the retirement of EC2-Classic the aws_redshift_security_group resource has been deprecated and will be removed in a future version.`,
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftConn

	var err error
	var errs []error

	name := d.Get("name").(string)
	desc := d.Get("description").(string)
	sgInput := &redshift.CreateClusterSecurityGroupInput{
		ClusterSecurityGroupName: aws.String(name),
		Description:              aws.String(desc),
	}
	log.Printf("[DEBUG] Redshift security group create: name: %s, description: %s", name, desc)
	_, err = conn.CreateClusterSecurityGroup(sgInput)
	if err != nil {
		return fmt.Errorf("Error creating RedshiftSecurityGroup: %s", err)
	}

	d.SetId(d.Get("name").(string))

	log.Printf("[INFO] Redshift Security Group ID: %s", d.Id())
	sg, err := resourceSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	ingresses := d.Get("ingress").(*schema.Set)
	for _, ing := range ingresses.List() {
		err := resourceSecurityGroupAuthorizeRule(ing, *sg.ClusterSecurityGroupName, conn)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}

	log.Println("[INFO] Waiting for Redshift Security Group Ingress Authorizations to be authorized")
	stateConf := &resource.StateChangeConf{
		Pending: []string{"authorizing"},
		Target:  []string{"authorized"},
		Refresh: resourceSecurityGroupStateRefreshFunc(d, meta),
		Timeout: 10 * time.Minute,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	sg, err := resourceSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	rules := &schema.Set{
		F: resourceSecurityGroupIngressHash,
	}

	for _, v := range sg.IPRanges {
		rule := map[string]interface{}{"cidr": aws.StringValue(v.CIDRIP)}
		rules.Add(rule)
	}

	for _, g := range sg.EC2SecurityGroups {
		rule := map[string]interface{}{
			"security_group_name":     aws.StringValue(g.EC2SecurityGroupName),
			"security_group_owner_id": aws.StringValue(g.EC2SecurityGroupOwnerId),
		}
		rules.Add(rule)
	}

	d.Set("ingress", rules)
	d.Set("name", sg.ClusterSecurityGroupName)
	d.Set("description", sg.Description)

	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftConn

	if d.HasChange("ingress") {
		o, n := d.GetChange("ingress")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removeIngressRules := expandSGRevokeIngress(os.Difference(ns).List())
		if len(removeIngressRules) > 0 {
			for _, r := range removeIngressRules {
				r.ClusterSecurityGroupName = aws.String(d.Id())

				_, err := conn.RevokeClusterSecurityGroupIngress(&r)
				if err != nil {
					return err
				}
			}
		}

		addIngressRules := expandSGAuthorizeIngress(ns.Difference(os).List())
		if len(addIngressRules) > 0 {
			for _, r := range addIngressRules {
				r.ClusterSecurityGroupName = aws.String(d.Id())

				_, err := conn.AuthorizeClusterSecurityGroupIngress(&r)
				if err != nil {
					return err
				}
			}
		}

	}
	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftConn

	log.Printf("[DEBUG] Redshift Security Group destroy: %v", d.Id())
	opts := redshift.DeleteClusterSecurityGroupInput{
		ClusterSecurityGroupName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Redshift Security Group destroy configuration: %v", opts)
	_, err := conn.DeleteClusterSecurityGroup(&opts)

	if err != nil {
		newerr, ok := err.(awserr.Error)
		if ok && newerr.Code() == "InvalidRedshiftSecurityGroup.NotFound" {
			return nil
		}
		return err
	}

	return nil
}

func resourceSecurityGroupRetrieve(d *schema.ResourceData, meta interface{}) (*redshift.ClusterSecurityGroup, error) {
	conn := meta.(*conns.AWSClient).RedshiftConn

	opts := redshift.DescribeClusterSecurityGroupsInput{
		ClusterSecurityGroupName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Redshift Security Group describe configuration: %#v", opts)

	resp, err := conn.DescribeClusterSecurityGroups(&opts)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving Redshift Security Groups: %s", err)
	}

	if len(resp.ClusterSecurityGroups) != 1 ||
		aws.StringValue(resp.ClusterSecurityGroups[0].ClusterSecurityGroupName) != d.Id() {
		return nil, fmt.Errorf("Unable to find Redshift Security Group: %#v", resp.ClusterSecurityGroups)
	}

	return resp.ClusterSecurityGroups[0], nil
}

func resourceSecurityGroupIngressHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["cidr"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_owner_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return create.StringHashcode(buf.String())
}

func resourceSecurityGroupAuthorizeRule(ingress interface{}, redshiftSecurityGroupName string, conn *redshift.Redshift) error {
	ing := ingress.(map[string]interface{})

	opts := redshift.AuthorizeClusterSecurityGroupIngressInput{
		ClusterSecurityGroupName: aws.String(redshiftSecurityGroupName),
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.CIDRIP = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = aws.String(attr.(string))
	}

	if attr, ok := ing["security_group_owner_id"]; ok && attr != "" {
		opts.EC2SecurityGroupOwnerId = aws.String(attr.(string))
	}

	log.Printf("[DEBUG] Authorize ingress rule configuration: %#v", opts)
	_, err := conn.AuthorizeClusterSecurityGroupIngress(&opts)

	if err != nil {
		return fmt.Errorf("Error authorizing security group ingress: %s", err)
	}

	return nil
}

func resourceSecurityGroupStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceSecurityGroupRetrieve(d, meta)

		if err != nil {
			log.Printf("Error on retrieving Redshift Security Group when waiting: %s", err)
			return nil, "", err
		}

		statuses := make([]string, 0, len(v.EC2SecurityGroups)+len(v.IPRanges))
		for _, ec2g := range v.EC2SecurityGroups {
			statuses = append(statuses, *ec2g.Status)
		}
		for _, ips := range v.IPRanges {
			statuses = append(statuses, *ips.Status)
		}

		for _, stat := range statuses {
			// Not done
			if stat != "authorized" {
				return nil, "authorizing", nil
			}
		}

		return v, "authorized", nil
	}
}

func expandSGAuthorizeIngress(configured []interface{}) []redshift.AuthorizeClusterSecurityGroupIngressInput {
	var ingress []redshift.AuthorizeClusterSecurityGroupIngressInput

	// Loop over our configured parameters and create
	// an array of aws-sdk-go compatible objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		i := redshift.AuthorizeClusterSecurityGroupIngressInput{}

		if v, ok := data["cidr"]; ok {
			i.CIDRIP = aws.String(v.(string))
		}

		if v, ok := data["security_group_name"]; ok {
			i.EC2SecurityGroupName = aws.String(v.(string))
		}

		if v, ok := data["security_group_owner_id"]; ok {
			i.EC2SecurityGroupOwnerId = aws.String(v.(string))
		}

		ingress = append(ingress, i)
	}

	return ingress
}

func expandSGRevokeIngress(configured []interface{}) []redshift.RevokeClusterSecurityGroupIngressInput {
	var ingress []redshift.RevokeClusterSecurityGroupIngressInput

	// Loop over our configured parameters and create
	// an array of aws-sdk-go compatible objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		i := redshift.RevokeClusterSecurityGroupIngressInput{}

		if v, ok := data["cidr"]; ok {
			i.CIDRIP = aws.String(v.(string))
		}

		if v, ok := data["security_group_name"]; ok {
			i.EC2SecurityGroupName = aws.String(v.(string))
		}

		if v, ok := data["security_group_owner_id"]; ok {
			i.EC2SecurityGroupOwnerId = aws.String(v.(string))
		}

		ingress = append(ingress, i)
	}

	return ingress
}
