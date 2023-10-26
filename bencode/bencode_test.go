package bencode

import (
	"reflect"
	"testing"
)

func TestEncodeBencode(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		result := EncodeBencode(42)
		expected := "i42e"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("String", func(t *testing.T) {
		result := EncodeBencode("hello")
		expected := "5:hello"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("List", func(t *testing.T) {
		result := EncodeBencode([]interface{}{1, "two", []interface{}{3, 4}})
		expected := "li1e3:twoli3ei4eee"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("Dictionary", func(t *testing.T) {
		data := map[string]interface{}{
			"key1": 42,
			"key2": "value2",
		}
		result := EncodeBencode(data)
		expected := "d4:key1i42e4:key26:value2e"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("Invalid Data Type", func(t *testing.T) {
		result := EncodeBencode([]int{1, 2, 3})
		expected := ""
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestDecodeBencode(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		result, _ := DecodeBencode("i42e")
		expected := 42
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("String", func(t *testing.T) {
		result, _ := DecodeBencode("5:hello")
		expected := "hello"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("List", func(t *testing.T) {
		// li1e3:twoli3ei4eee
		result, _ := DecodeBencode("li1e3:twoli3ei4eee")
		expected := []interface{}{1, "two", []interface{}{3, 4}}

		areEqual := reflect.DeepEqual(expected, result)
		if !areEqual {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Dictionary", func(t *testing.T) {
		data := map[string]interface{}{
			"key1": 42,
			"key2": []interface{}{1},
		}
		result, _ := DecodeBencode("d4:key1i42e4:key2li1eee")
		expected := data

		areEqual := reflect.DeepEqual(expected, result)
		if !areEqual {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	//t.Run("Invalid Data Type", func(t *testing.T) {
	//	_, err := DecodeBencode("xyz")
	//	if err == nil {
	//		t.Errorf("Expected an error but got nil")
	//	}
	//})
}
