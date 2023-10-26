package bencode

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

// EncodeBencode encodes data to bencode format
func EncodeBencode(data interface{}) string {
	switch v := data.(type) {
	case int:
		return encodeInt(v)
	case string:
		return encodeString(v)
	case []interface{}:
		return encodeList(v)
	case map[string]interface{}:
		return encodeDictionary(v)
	default:
		return ""
	}
}

func encodeInt(val int) string {
	return "i" + strconv.Itoa(val) + "e"
}

func encodeString(val string) string {
	return strconv.Itoa(len(val)) + ":" + val
}

func encodeList(val []interface{}) string {
	var result string
	for _, v := range val {
		result += EncodeBencode(v)
	}
	return "l" + result + "e"
}

func encodeDictionary(val map[string]interface{}) string {
	keys := make([]string, 0, len(val))
	for k := range val {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result string
	for _, k := range keys {
		result += EncodeBencode(k) + EncodeBencode(val[k])
	}

	return "d" + result + "e"
}

// DecodeBencode decodes bencoded string to data
func DecodeBencode(bencodedString string) (interface{}, error) {
	val, _, err := decodeAllBencode(bencodedString, 0)
	return val, err
}

func decodeAllBencode(bencodedString string, startIndex int) (interface{}, int, error) {
	index := startIndex

	for index < len(bencodedString) {
		switch bencodedString[index] {
		case 'e':
			return nil, index + 1, nil
		case 'i':
			return decodeIntBencode(bencodedString, index+1)
		case 'l':
			return decodeListBencode(bencodedString, index+1)
		case 'd':
			return decodeDictBencode(bencodedString, index+1)
		default:
			return decodeStringBencode(bencodedString, index)
		}
	}

	return nil, index, fmt.Errorf("unexpected end of input")
}

func decodeDictBencode(bencodedString string, startIndex int) (map[string]interface{}, int, error) {
	index := startIndex
	ret := make(map[string]interface{})
	var key = ""
	for bencodedString[index] != 'e' {
		if unicode.IsDigit(rune(bencodedString[index])) {
			val, newIndex, err := decodeStringBencode(bencodedString, index)
			if err != nil {
				return ret, 0, err
			}
			index = newIndex
			if key == "" {
				key = val
			} else {
				ret[key] = val
				key = ""
			}

		} else if rune(bencodedString[index]) == 'i' {
			val, newIndex, err := decodeIntBencode(bencodedString, index+1)
			if err != nil {
				return ret, 0, err
			}
			index = newIndex
			ret[key] = val
			key = ""
		} else if rune(bencodedString[index]) == 'l' {
			val, newIndex, err := decodeAllBencode(bencodedString, index)
			if err != nil {
				return ret, 0, err
			}
			index = newIndex
			ret[key] = val
			key = ""
		} else if rune(bencodedString[index]) == 'd' {
			val, newIndex, err := decodeDictBencode(bencodedString, index+1)
			if err != nil {
				return ret, 0, err
			}
			index = newIndex
			ret[key] = val
			key = ""
		}
	}
	return ret, index + 1, nil
}

func decodeIntBencode(bencodedString string, startIndex int) (int, int, error) {
	endIndex := startIndex

	for endIndex < len(bencodedString) && bencodedString[endIndex] != 'e' {
		endIndex++
	}

	numberStr := bencodedString[startIndex:endIndex]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, 0, err
	}

	return number, endIndex + 1, nil
}

func decodeStringBencode(bencodedString string, startIndex int) (string, int, error) {
	colonIndex := startIndex

	for colonIndex < len(bencodedString) && bencodedString[colonIndex] != ':' {
		colonIndex++
	}

	lengthStr := bencodedString[startIndex:colonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, err
	}

	dataStart := colonIndex + 1
	dataEnd := dataStart + length

	if dataEnd > len(bencodedString) {
		return "", 0, fmt.Errorf("string length exceeds available data")
	}

	return bencodedString[dataStart:dataEnd], dataEnd, nil
}

func decodeListBencode(bencodedString string, startIndex int) ([]interface{}, int, error) {
	index := startIndex
	var retList []interface{}

	for bencodedString[index] != 'e' {
		val, newIndex, err := decodeAllBencode(bencodedString, index)
		if err != nil {
			return nil, 0, err
		}
		index = newIndex
		retList = append(retList, val)
	}

	return retList, index + 1, nil
}
