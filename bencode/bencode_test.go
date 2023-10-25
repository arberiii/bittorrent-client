package bencode

import "testing"

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

