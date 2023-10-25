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


// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		value, _, err := decodeStringBencode(bencodedString, 0)
		return value, err
	} else if rune(bencodedString[0]) == 'i' {
	    value, _, err := decodeIntBencode(bencodedString, 1)
	    return value, err
	} else if rune(bencodedString[0]) == 'l' {
	    value, _, err := decodeListBencode(bencodedString, 1)
	    return value, err
	} else if rune(bencodedString[0]) == 'd' {
	    value, _, err := decodeDictBencode(bencodedString, 1)
	    return value, err
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
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
            if (key == "") {
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
            val, newIndex, err := decodeListBencode(bencodedString, index+1)
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

func decodeListBencode(bencodedString string, startIndex int) ([]interface{}, int, error) {
    index := startIndex
    var ret []interface{}
    for bencodedString[index] != 'e' {
        if unicode.IsDigit(rune(bencodedString[index])) {
            val, newIndex, err := decodeStringBencode(bencodedString, index)
            if err != nil {
                return ret, 0, err
            }
            index = newIndex
            ret = append(ret, val)
        } else if rune(bencodedString[index]) == 'i' {
            val, newIndex, err := decodeIntBencode(bencodedString, index+1)
            if err != nil {
                return ret, 0, err
            }
            index = newIndex
            ret = append(ret, val)
        }
    }
    return ret, index + 1, nil
}

func decodeIntBencode(bencodedString string, startIndex int) (int, int, error) {
    index := startIndex
    isNegative := rune(bencodedString[startIndex]) == '-'
    if (isNegative) {
        index += 1
    }
    var lastIndex int

    for i := index; i < len(bencodedString); i++ {
        if bencodedString[i] == 'e' {
            lastIndex = i
            break
        }
    }

    numberStr := bencodedString[index:lastIndex]

    number, err := strconv.Atoi(numberStr)
    if err != nil {
        return 0, 0, err
    }
    if (isNegative) {
        return number * -1, lastIndex+1, nil
    }
    return number, lastIndex + 1, nil
}

func decodeStringBencode(bencodedString string, startIndex int) (string, int, error) {
    var firstColonIndex int

    for i := startIndex; i < len(bencodedString); i++ {
        if bencodedString[i] == ':' {
            firstColonIndex = i
            break
        }
    }
    lengthStr := bencodedString[startIndex:firstColonIndex]

    length, err := strconv.Atoi(lengthStr)
    if err != nil {
        return "", 0, err
    }

    return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex+1+length, nil
}