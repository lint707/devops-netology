package waf

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

type WafRetryer struct {
	Connection *waf.WAF
}

type withTokenFunc func(token *string) (interface{}, error)

func (t *WafRetryer) RetryWithToken(f withTokenFunc) (interface{}, error) {
	conns.GlobalMutexKV.Lock("WafRetryer")
	defer conns.GlobalMutexKV.Unlock("WafRetryer")

	var out interface{}
	var tokenOut *waf.GetChangeTokenOutput
	err := resource.Retry(15*time.Minute, func() *resource.RetryError {
		var err error
		tokenOut, err = t.Connection.GetChangeToken(&waf.GetChangeTokenInput{})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Failed to acquire change token: %w", err))
		}

		out, err = f(tokenOut.ChangeToken)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, waf.ErrCodeStaleDataException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		tokenOut, err = t.Connection.GetChangeToken(&waf.GetChangeTokenInput{})

		if err != nil {
			return nil, fmt.Errorf("error getting WAF change token: %w", err)
		}

		out, err = f(tokenOut.ChangeToken)
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewRetryer(conn *waf.WAF) *WafRetryer {
	return &WafRetryer{Connection: conn}
}
