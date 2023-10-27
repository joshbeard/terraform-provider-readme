package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Image_Resource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "readme_image" "test" {
					source = "../../examples/resources/readme_image/example.png"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"readme_image.test",
						"url",
						regexp.MustCompile(`^https://files.readme.io/[a-z0-9]{6,}-example\.png$`),
					),
					resource.TestMatchResourceAttr(
						"readme_image.test",
						"filename",
						regexp.MustCompile(`^[a-z0-9]{6,}-example\.png$`),
					),
					resource.TestCheckResourceAttr(
						"readme_image.test",
						"width",
						"1",
					),
					resource.TestCheckResourceAttr(
						"readme_image.test",
						"height",
						"1",
					),
					resource.TestCheckResourceAttr(
						"readme_image.test",
						"color",
						"#371ca1",
					),
				),
			},
		},
	})
}
