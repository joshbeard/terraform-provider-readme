package readme_unit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestAPIRequestOptions(t *testing.T) {
	// Test case 1: Version is set.
	version := basetypes.NewStringValue("1.0")
	options := apiRequestOptions(version)
	if options.Version != "1.0" {
		t.Errorf("Expected version '1.0', but got: %s", options.Version)
	}

	// Test case 2: Version is not set.
	version = basetypes.NewStringValue("")
	options = apiRequestOptions(version)
	if options.Version != "" {
		t.Errorf("Expected an empty version, but got: %s", options.Version)
	}
}
