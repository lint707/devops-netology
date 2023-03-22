package events_test

import (
	"context"
	"reflect"
	"testing"

	tfevents "github.com/hashicorp/terraform-provider-aws/internal/service/events"
)

func testResourceTargetStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"arn":       "arn:aws:test:us-east-1:123456789012:test", //lintignore:AWSAT003,AWSAT005
		"rule":      "testrule",
		"target_id": "testtargetid",
	}
}

func testResourceTargetStateDataV0EventBusName() map[string]interface{} {
	return map[string]interface{}{
		"arn":            "arn:aws:test:us-east-1:123456789012:test", //lintignore:AWSAT003,AWSAT005
		"event_bus_name": "testbus",
		"rule":           "testrule",
		"target_id":      "testtargetid",
	}
}

func testResourceTargetStateDataV1() map[string]interface{} {
	v0 := testResourceTargetStateDataV0()
	return map[string]interface{}{
		"arn":            v0["arn"],
		"event_bus_name": "default",
		"rule":           v0["rule"],
		"target_id":      v0["target_id"],
	}
}

func testResourceTargetStateDataV1EventBusName() map[string]interface{} {
	v0 := testResourceTargetStateDataV0EventBusName()
	return map[string]interface{}{
		"arn":            v0["arn"],
		"event_bus_name": v0["event_bus_name"],
		"rule":           v0["rule"],
		"target_id":      v0["target_id"],
	}
}

func TestTargetStateUpgradeV0(t *testing.T) {
	expected := testResourceTargetStateDataV1()
	actual, err := tfevents.TargetStateUpgradeV0(context.Background(), testResourceTargetStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestTargetStateUpgradeV0EventBusName(t *testing.T) {
	expected := testResourceTargetStateDataV1EventBusName()
	actual, err := tfevents.TargetStateUpgradeV0(context.Background(), testResourceTargetStateDataV0EventBusName(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
