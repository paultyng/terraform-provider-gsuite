package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// example spreadsheet from Google's docs: https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
const exampleSpreadsheetID = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"

func TestAccDataSourceSheetsSpreadsheetValues_basic(t *testing.T) {
	const readRange = "Class Data!A2:B"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSheetsSpreadsheetValues(exampleSpreadsheetID, readRange),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("A1", "Alexandra"),
					resource.TestCheckOutput("B1", "Female"),
					resource.TestCheckOutput("B2", "Male"),
				),
			},
		},
	})
}

func testAccDataSourceSheetsSpreadsheetValues(spreadsheetID, readRange string) string {
	return fmt.Sprintf(`
data "gsuite_sheets_spreadsheet_values" "test" {
	spreadsheet_id = "%s"
	range = "%s"
}

output "A1" {
	value = data.gsuite_sheets_spreadsheet_values.test.values[0][0]
}

output "B1" {
	value = data.gsuite_sheets_spreadsheet_values.test.values[0][1]
}

output "B2" {
	value = data.gsuite_sheets_spreadsheet_values.test.values[1][1]
}
`, spreadsheetID, readRange)
}
