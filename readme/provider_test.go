package readme

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/liveoaklabs/readme-api-go-client/readme"
)

func TestAPIRequestOptions(t *testing.T) {
	// Test case 1: Version is set.
	version := basetypes.NewStringValue("1.2")
	options := apiRequestOptions(version)
	if options.Version != "1.2" {
		t.Errorf("Expected version '1.2', but got: %s", options.Version)
	}

	// Test case 2: Version is not set.
	version = basetypes.NewStringValue("")
	options = apiRequestOptions(version)
	if options.Version != "" {
		t.Errorf("Expected an empty version, but got: %s", options.Version)
	}
}

func TestClientError(t *testing.T) {
	// Test case 1: APIResponse with an error message.
	apiResponse := &readme.APIResponse{
		APIErrorResponse: readme.APIErrorResponse{
			Message: "API Error Message",
		},
	}
	err := fmt.Errorf("Some error")
	errorMsg := clientError(err, apiResponse)
	expectedErrorMsg := "API Error Message\nAPI Error Response: {Message:API Error Message}\n"
	if errorMsg != expectedErrorMsg {
		t.Errorf("Expected error message:\n%s\nbut got:\n%s", expectedErrorMsg, errorMsg)
	}

	// Test case 2: APIResponse without an error message.
	apiResponse = &readme.APIResponse{}
	err = fmt.Errorf("Some error")
	errorMsg = clientError(err, apiResponse)
	if errorMsg != "Some error" {
		t.Errorf("Expected error message 'Some error', but got: %s", errorMsg)
	}
}
