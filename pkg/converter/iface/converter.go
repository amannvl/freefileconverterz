// Package iface defines the interfaces for the converter package
package iface

import (
	"context"
	"fmt"
)

// Converter defines the interface for all converters
type Converter interface {
	// Convert converts the input file to the output file with the target format
	Convert(ctx context.Context, inputPath, outputPath string) error

	// SupportsConversion checks if the converter supports the given conversion
	SupportsConversion(sourceFormat, targetFormat string) bool

	// Cleanup removes any temporary files
	Cleanup(files ...string) error
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
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
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
