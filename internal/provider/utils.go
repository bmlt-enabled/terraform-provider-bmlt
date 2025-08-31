package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// Helper function to convert string to nullable string value
// Returns types.StringNull() for empty strings, types.StringValue(s) otherwise
func nullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
