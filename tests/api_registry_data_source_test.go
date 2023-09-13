package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAPIRegistryDataSource(t *testing.T) {
	tfConfig := `
	resource "readme_api_specification" "test" {
						definition      = file("testdata/example.json")
						delete_category = true
	}
	data "readme_api_registry" "test" { uuid = readme_api_specification.test.uuid }
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tfConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test the UUID (dynamic)
					// TODO: Test the definition
					// resource.TestCheckResourceAttr(
					// 	"data.readme_api_registry.test",
					// 	"uuid",
					// 	// "somethingUnique",
					// 	"readme_api_specification.test.uuid",
					// ),

					resource.TestCheckResourceAttr("data.readme_api_registry.test", "id", "readme"),
					// resource.TestCheckResourceAttr(
					// 	"data.readme_api_registry.test",
					// 	"definition",
					// 	`{"one": "two"}`,
					// ),
				),
			},
		},
	})
}

func TestAPIRegistryDataSource_GetError(t *testing.T) {
	expectError, _ := regexp.Compile(
		`Unable to retrieve API registry metadata\.`,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "readme_api_registry" "test" { uuid = "somethingUnique" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_api_registry.test",
						"uuid",
						"somethingUnique",
					),
				),

				ExpectError: expectError,
			},
		},
	})
}

// func testAccCheckValues(widget *example.Widget, name string) resource.TestCheckFunc {
//     return func(s *terraform.State) error {
//         if *widget.Active != true {
//             return fmt.Errorf("bad active state, expected \"true\", got: %#v", *widget.Active)
//         }
//         if *widget.Name != name {
//             return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, *widget.Name)
//         }
//         return nil
//     }
// }
//
// // testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
// func testAccCheckExampleResourceExists(n string, widget *example.Widget) resource.TestCheckFunc {
//     return func(s *terraform.State) error {
//         // find the corresponding state object
//         rs, ok := s.RootModule().Resources[n]
//         if !ok {
//             return fmt.Errorf("Not found: %s", n)
//         }
//
//         // retrieve the configured client from the test setup
//         conn := testAccProvider.Meta().(*ExampleClient)
//         resp, err := conn.DescribeWidget(&example.DescribeWidgetsInput{
//             WidgetIdentifier: rs.Primary.ID,
//         })
//
//         if err != nil {
//             return err
//         }
//
//         if resp.Widget == nil {
//             return fmt.Errorf("Widget (%s) not found", rs.Primary.ID)
//         }
//
//         // assign the response Widget attribute to the widget pointer
//         *widget = *resp.Widget
//
//         return nil
//     }
// }
