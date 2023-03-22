package accessanalyzer_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// AccessAnalyzer is limited to one per region, so run serially
// locally and in TeamCity.
func TestAccAccessAnalyzer_serial(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"Analyzer": {
			"basic":             testAccAnalyzer_basic,
			"disappears":        testAccAnalyzer_disappears,
			"Tags":              testAccAnalyzer_Tags,
			"Type_Organization": testAccAnalyzer_Type_Organization,
		},
		"ArchiveRule": {
			"basic":          testAccAnalyzerArchiveRule_basic,
			"disappears":     testAccAnalyzerArchiveRule_disappears,
			"update_filters": testAccAnalyzerArchiveRule_updateFilters,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AccessAnalyzerConn

	input := &accessanalyzer.ListAnalyzersInput{}

	_, err := conn.ListAnalyzers(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
