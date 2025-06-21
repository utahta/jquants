package jquants

import (
	"fmt"
)

// Test helper functions for creating pointers to primitive types

func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}


// Helper functions for comparing pointer values

func compareFloat64Ptr(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func ptrToStr(p *float64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%f", *p)
}