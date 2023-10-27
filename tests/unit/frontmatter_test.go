package readme_unit_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/liveoaklabs/terraform-provider-readme/readme/frontmatter"
)

func TestGetValue(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Attribute exists in the front matter.
	// Arrange
	body := `---
    title: Sample Title
    categorySlug: sample-category
    ---`

	// Act
	for attribute, expectedValue := range map[string]string{
		"Title":        "Sample Title",
		"CategorySlug": "sample-category",
	} {
		value, err := frontmatter.GetValue(ctx, body, attribute)
		if err != "" {
			t.Errorf("Expected no error, but got: %s", err)
		}

		// Assert
		if !value.IsValid() {
			t.Errorf("Expected a valid value, but got an invalid value: %v", value)
		} else if value.Kind() != reflect.String || value.String() != expectedValue {
			t.Errorf("Expected '%s', but got: %v", expectedValue, value)
		}
	}

	// Test case 2: Attribute does not exist in the front matter.
	// Arrange
	body = `---
    categorySlug: sample-category
    ---`

	// Act
	attribute := "Title"
	value, err := frontmatter.GetValue(ctx, body, attribute)
	if err != "" {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Assert
	if value.IsValid() {
		t.Errorf("Expected an invalid value, but got a valid value")
	}

	// Test case 3: Empty front matter.
	// Arrange
	body = `This document`

	// Act
	attribute = "Title"
	value, err = frontmatter.GetValue(ctx, body, attribute)
	if err != "" {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Assert
	if value.IsValid() {
		t.Errorf("Expected an invalid value, but got a valid value")
	}
}
