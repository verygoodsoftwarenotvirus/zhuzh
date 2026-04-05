package main

import (
	"testing"
)

func TestSchemaHelpers(t *testing.T) {
	t.Run("optionalFloatRangeSchema", func(t *testing.T) {
		result := optionalFloatRangeSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("uint16RangeWithOptionalMaxSchema", func(t *testing.T) {
		result := uint16RangeWithOptionalMaxSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("float32RangeWithOptionalMaxSchema", func(t *testing.T) {
		result := float32RangeWithOptionalMaxSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("uint32RangeWithOptionalMaxSchema", func(t *testing.T) {
		result := uint32RangeWithOptionalMaxSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("optionalUint32RangeSchema", func(t *testing.T) {
		result := optionalUint32RangeSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("optionalFloat32RangeSchema", func(t *testing.T) {
		result := optionalFloat32RangeSchema()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("floatField", func(t *testing.T) {
		result := floatField("test description")
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result["type"] != numberType {
			t.Errorf("expected type %q, got %q", numberType, result["type"])
		}
	})

	t.Run("uintField", func(t *testing.T) {
		result := uintField("test description")
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result["type"] != intType {
			t.Errorf("expected type %q, got %q", intType, result["type"])
		}
		if result["minimum"] != 0 {
			t.Errorf("expected minimum 0, got %v", result["minimum"])
		}
	})
}
