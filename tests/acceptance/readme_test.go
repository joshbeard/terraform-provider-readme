package readme_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type testCase struct {
	name        string
	config      string
	expectError string
	checks      []resource.TestCheckFunc // List of TestCheckFuncs to run against the data source response.
}

type testCheck struct {
	resource string
	attr     string
	value    string
}

func TestMain(m *testing.M) {
	// Setup
	code := m.Run()

	// Teardown
	fmt.Println("Cleaning up...")
	teardown()
	os.Exit(code)
}

func GenerateTests(checks []testCheck, testCases []testCase) (string, string, []resource.TestCheckFunc) {
	cfg, errCfg := "", ""
	testChecks := []resource.TestCheckFunc{}

	for i, tc := range testCases { //nolint:varnamelen
		// Replaces spaces with underscores in the test case name.
		tc.name = regexp.MustCompile(`\s+`).ReplaceAllString(tc.name, "_")
		// Replace commas
		tc.name = regexp.MustCompile(`,`).ReplaceAllString(tc.name, "_")

		// Replace $TEST$ with the test case name.
		rscName := fmt.Sprintf("_%d_%s", i, tc.name)

		if len(tc.checks) > 0 {
			testChecks = append(testChecks, tc.checks...)
		} else {
			for _, check := range checks {
				rsc := check.resource + rscName
				testChecks = append(testChecks, resource.TestCheckResourceAttr(rsc, check.attr, check.value))
			}
		}

		tc.config = regexp.MustCompile(`\$TEST\$`).ReplaceAllString(tc.config, rscName)
		if tc.expectError == "" {
			cfg += tc.config + "\n"
		} else {
			errCfg += tc.config + "\n"
		}
	}

	return cfg, errCfg, testChecks
}
