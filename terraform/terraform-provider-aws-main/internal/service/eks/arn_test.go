/*
This file is a hard copy of:
https://github.com/kubernetes-sigs/aws-iam-authenticator/blob/7547c74e660f8d34d9980f2c69aa008eed1f48d0/pkg/arn/arn_test.go

With the following modifications:
 - Rename package eks
*/

package eks_test

import (
	"fmt"
	"testing"

	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
)

var arnTests = []struct {
	arn      string // input arn
	expected string // canonacalized arn
	err      error  // expected error value
}{
	{"NOT AN ARN", "", fmt.Errorf("Not an arn")},
	//lintignore:AWSAT005
	{"arn:aws:iam::123456789012:user/Alice", "arn:aws:iam::123456789012:user/Alice", nil},
	//lintignore:AWSAT005
	{"arn:aws:iam::123456789012:role/Users", "arn:aws:iam::123456789012:role/Users", nil},
	//lintignore:AWSAT005
	{"arn:aws:sts::123456789012:assumed-role/Admin/Session", "arn:aws:iam::123456789012:role/Admin", nil},
	//lintignore:AWSAT005
	{"arn:aws:sts::123456789012:federated-user/Bob", "arn:aws:sts::123456789012:federated-user/Bob", nil},
	//lintignore:AWSAT005
	{"arn:aws:iam::123456789012:root", "arn:aws:iam::123456789012:root", nil},
	//lintignore:AWSAT005
	{"arn:aws:sts::123456789012:assumed-role/Org/Team/Admin/Session", "arn:aws:iam::123456789012:role/Org/Team/Admin", nil},
}

func TestUserARN(t *testing.T) {
	for _, tc := range arnTests {
		actual, err := tfeks.Canonicalize(tc.arn)
		if err != nil && tc.err == nil || err == nil && tc.err != nil {
			t.Errorf("Canoncialize(%s) expected err: %v, actual err: %v", tc.arn, tc.err, err)
			continue
		}
		if actual != tc.expected {
			t.Errorf("tfeks.Canonicalize(%s) expected: %s, actual: %s", tc.arn, tc.expected, actual)
		}
	}
}
