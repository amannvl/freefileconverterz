package base

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
)

// BaseConverter provides common functionality for all converters
type BaseConverter struct {
	name            string
	supportedFormats map[string][]string
}

// NewBaseConverter creates a new BaseConverter
func NewBaseConverter(name string, supportedFormats map[string][]string) *BaseConverter {
	return &BaseConverter{
		name:            name,
		supportedFormats: supportedFormats,
	}
}

// SupportedFormats returns the supported formats for this converter
func (c *BaseConverter) SupportedFormats() map[string][]string {
	return c.supportedFormats
}

// ValidateOptions validates the conversion options
func (c *BaseConverter) ValidateOptions(options map[string]interface{}) error {
	// Default implementation accepts any options
	// Individual converters can override this method
	return nil
}

// SupportsConversion checks if the converter supports the given conversion
func (c *BaseConverter) SupportsConversion(sourceFormat, targetFormat string) bool {
	supported, exists := c.supportedFormats[sourceFormat]
	if !exists {
		return false
	}

	for _, format := range supported {
		if format == targetFormat {
			return true
		}
	}
	return false
}

// Convert is a base implementation that should be overridden by specific converters
func (c *BaseConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
	// Default implementation reads the input and returns it as-is
	// This should be overridden by specific converters
	if _, err := ioutil.ReadAll(input); err != nil {
		return nil, iface.NewConversionError("read_error", "failed to read input", err)
	}
	return nil, fmt.Errorf("conversion not implemented for %s converter", c.name)
}

// Ensure BaseConverter implements iface.Converter
var _ iface.Converter = (*BaseConverter)(nil)
