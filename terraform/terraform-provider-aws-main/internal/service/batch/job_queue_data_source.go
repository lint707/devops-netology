package batch

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceJobQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJobQueueRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scheduling_policy_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tftags.TagsSchemaComputed(),

			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"compute_environment_order": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_environment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceJobQueueRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BatchConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	params := &batch.DescribeJobQueuesInput{
		JobQueues: []*string{aws.String(d.Get("name").(string))},
	}
	log.Printf("[DEBUG] Reading Batch Job Queue: %s", params)
	desc, err := conn.DescribeJobQueues(params)

	if err != nil {
		return err
	}

	if len(desc.JobQueues) == 0 {
		return fmt.Errorf("no matches found for name: %s", d.Get("name").(string))
	}

	if len(desc.JobQueues) > 1 {
		return fmt.Errorf("multiple matches found for name: %s", d.Get("name").(string))
	}

	jobQueue := desc.JobQueues[0]
	d.SetId(aws.StringValue(jobQueue.JobQueueArn))
	d.Set("arn", jobQueue.JobQueueArn)
	d.Set("name", jobQueue.JobQueueName)
	d.Set("scheduling_policy_arn", jobQueue.SchedulingPolicyArn)
	d.Set("status", jobQueue.Status)
	d.Set("status_reason", jobQueue.StatusReason)
	d.Set("state", jobQueue.State)
	d.Set("priority", jobQueue.Priority)

	ceos := make([]map[string]interface{}, 0)
	for _, v := range jobQueue.ComputeEnvironmentOrder {
		ceo := map[string]interface{}{}
		ceo["compute_environment"] = aws.StringValue(v.ComputeEnvironment)
		ceo["order"] = int(aws.Int64Value(v.Order))
		ceos = append(ceos, ceo)
	}
	if err := d.Set("compute_environment_order", ceos); err != nil {
		return fmt.Errorf("error setting compute_environment_order: %w", err)
	}

	if err := d.Set("tags", KeyValueTags(jobQueue.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
