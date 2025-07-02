package converter

import (
	"context"
	"io"
)

// Converter defines the interface for all converters
type Converter interface {
	// Convert converts the input file to the target format
	Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error)
	
	// SupportedFormats returns a map of supported source formats to target formats
	SupportedFormats() map[string][]string
	
	// ValidateOptions validates the conversion options
	ValidateOptions(options map[string]interface{}) error
}

// ConversionError represents a conversion error
type ConversionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *ConversionError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ConversionError) Unwrap() error {
	return e.Err
}

// NewConversionError creates a new ConversionError
func NewConversionError(code, message string, err error) *ConversionError {
	return &ConversionError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
