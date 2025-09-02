package provider

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HTTP status code constants
const (
	HTTPStatusOK        = 200
	HTTPStatusCreated   = 201
	HTTPStatusNoContent = 204
	HTTPStatusNotFound  = 404
)

// safeInt64ToInt32 safely converts int64 to int32, clamping to int32 limits
func safeInt64ToInt32(value int64) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}

// Helper function for client type error messages
func clientTypeError(actualType interface{}) string {
	return fmt.Sprintf("Expected *BMTLClientData, got: %T. "+
		"Please report this issue to the provider developers.", actualType)
}

// Helper function to convert string to nullable string value
// Returns types.StringNull() for empty strings, types.StringValue(s) otherwise
func nullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
