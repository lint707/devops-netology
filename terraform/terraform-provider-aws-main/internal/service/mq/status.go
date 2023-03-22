package mq

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mq"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func StatusBroker(conn *mq.MQ, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := conn.DescribeBroker(&mq.DescribeBrokerInput{
			BrokerId: aws.String(id),
		})

		if tfawserr.ErrCodeEquals(err, mq.ErrCodeNotFoundException) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if output == nil {
			return nil, "", nil
		}

		return output, aws.StringValue(output.BrokerState), nil
	}
}
