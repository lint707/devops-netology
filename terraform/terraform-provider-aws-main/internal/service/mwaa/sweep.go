//go:build sweep
// +build sweep

package mwaa

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mwaa"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_mwaa_environment", &resource.Sweeper{
		Name: "aws_mwaa_environment",
		F:    sweepEnvironment,
	})
}

func sweepEnvironment(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).MWAAConn

	listOutput, err := conn.ListEnvironments(&mwaa.ListEnvironmentsInput{})
	if err != nil {
		if sweep.SkipSweepError(err) || tfawserr.ErrCodeEquals(err, "InternalFailure") {
			log.Printf("[WARN] Skipping MWAA Environment sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("Error retrieving MWAA Environment: %s", err)
	}
	for _, environment := range listOutput.Environments {
		name := aws.StringValue(environment)
		r := ResourceEnvironment()
		d := r.Data(nil)
		d.SetId(name)

		err := r.Delete(d, client)
		if err != nil {
			log.Printf("[ERROR] Failed to delete MWAA Environment %s: %s", name, err)
		}
	}
	return nil
}
