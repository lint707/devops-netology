package servicequotas

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func findServiceQuotaDefaultByID(conn *servicequotas.ServiceQuotas, serviceCode, quotaCode string) (*servicequotas.ServiceQuota, error) {
	input := &servicequotas.GetAWSDefaultServiceQuotaInput{
		ServiceCode: aws.String(serviceCode),
		QuotaCode:   aws.String(quotaCode),
	}

	output, err := conn.GetAWSDefaultServiceQuota(input)

	if err != nil {
		return nil, err
	}
	if output == nil || output.Quota == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Quota, nil
}

func findServiceQuotaDefaultByName(conn *servicequotas.ServiceQuotas, serviceCode, quotaName string) (*servicequotas.ServiceQuota, error) {
	input := &servicequotas.ListAWSDefaultServiceQuotasInput{
		ServiceCode: aws.String(serviceCode),
	}

	var defaultQuota *servicequotas.ServiceQuota
	err := conn.ListAWSDefaultServiceQuotasPages(input, func(page *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, q := range page.Quotas {
			if aws.StringValue(q.QuotaName) == quotaName {
				defaultQuota = q
				return false
			}
		}

		return !lastPage
	})
	if err != nil {
		return nil, err
	}
	if defaultQuota == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return defaultQuota, nil
}

func findServiceQuotaByID(conn *servicequotas.ServiceQuotas, serviceCode, quotaCode string) (*servicequotas.ServiceQuota, error) {
	input := &servicequotas.GetServiceQuotaInput{
		ServiceCode: aws.String(serviceCode),
		QuotaCode:   aws.String(quotaCode),
	}

	output, err := conn.GetServiceQuota(input)

	if tfawserr.ErrCodeEquals(err, servicequotas.ErrCodeNoSuchResourceException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}
	if err != nil {
		return nil, err
	}

	if output == nil || output.Quota == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if output.Quota.ErrorReason != nil {
		return nil, &resource.NotFoundError{
			Message:     fmt.Sprintf("%s: %s", aws.StringValue(output.Quota.ErrorReason.ErrorCode), aws.StringValue(output.Quota.ErrorReason.ErrorMessage)),
			LastRequest: input,
		}
	}

	if output.Quota.Value == nil {
		return nil, &resource.NotFoundError{
			Message:     "empty value",
			LastRequest: input,
		}
	}

	return output.Quota, nil
}
