package mwaa

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mwaa"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	environmentStatusNotFound = "NotFound"
	environmentStatusUnknown  = "Unknown"
)

// statusEnvironment fetches the Environment and its Status
func statusEnvironment(conn *mwaa.MWAA, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		environment, err := findEnvironmentByName(conn, name)

		if tfawserr.ErrCodeEquals(err, mwaa.ErrCodeResourceNotFoundException) {
			return nil, environmentStatusNotFound, nil
		}

		if err != nil {
			return nil, environmentStatusUnknown, err
		}

		if environment == nil {
			return nil, environmentStatusNotFound, nil
		}

		return environment, aws.StringValue(environment.Status), nil
	}
}
