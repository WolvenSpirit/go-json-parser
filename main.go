package gojsonparser

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unicode"
)

const (
	twoColon    = ':'
	doubleQuote = '"'
	singleQuote = '\''
	braceOpen   = '{'
	braceClose  = '}'
	semiColon   = ','
)

// Value ...
type Value struct {
	String string
	Map    map[string]string
}

func detect(r *bufio.Reader, delim rune) (string, error) {
	var foundBefore bool
	var s []rune
	var i int
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if rn == delim && !foundBefore {
			foundBefore = true
			i++
			continue
		}
		if foundBefore && i > 0 {
			if rn == delim {
				break
			} else {
				s = append(s, rn)
			}
			i++
		}
	}
	return string(s), nil
}

func detectSimpleValue(r *bufio.Reader) (s string, err error) {
	var i int
	var tryParse bool
	var isDigit bool
	var isString bool
	var isObject bool
	var digit []rune
	var str []rune
	var lastIndexAppended int
	for j := 0; j < r.Size(); j++ {
		rn, n, err := r.ReadRune()
		if err != nil || n == -1 {
			if err == io.EOF {
				if isDigit {
					return string(digit), nil
				}
				if isString {
					return string(str), nil
				}
			}
			return s, err
		}
		if rn == twoColon && !isObject {
			i++
			tryParse = true
			continue
		}
		if i > 0 && tryParse && i != lastIndexAppended {
			if rn == braceOpen {
				isObject = true
				r.UnreadRune()
				return detectObject(r)
			}
			if unicode.IsDigit(rn) && !isObject {
				isDigit = true
				digit = append(digit, rn)
				lastIndexAppended = i
			}
			if rn == semiColon && isDigit && !isObject {
				return string(digit), nil
			}
			if unicode.IsLetter(rn) && !isDigit && !isObject {
				isString = true
				str = append(str, rn)
			}
			if rn == semiColon && isString && !isDigit && !isObject {
				return string(str), nil
			}
			i++
		}
	}
	return
}

func detectObject(r *bufio.Reader) (string, error) {
	var foundBefore bool
	var open int = 0
	var s []rune
	var i int
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return string(s), err
			}
			return "", err
		}
		if rn == braceOpen {
			open++
		}
		if rn == braceOpen && !foundBefore {
			foundBefore = true
			s = append(s, rn)
			i++
			continue
		}
		if foundBefore && i > 0 {
			if rn == braceClose {
				open--
			}
			if rn == braceClose && open == 0 {
				s = append(s, rn)
				return string(s), nil
			}
			s = append(s, rn)
			i++
		}
	}
}

// Parse takes a *bufio.Reader containing a json string and parses it into a map[string]string.
//
// Note: All values are returned as string, please use strconv to parse numeric values.
//
// Parse does only a one dimension parse, if a key has an object then that object will be returned as a string
// you can additionally call Parse on that object string subsequently to get it in a map[string]string.
func Parse(r *bufio.Reader) (map[string]string, error) {
	m := make(map[string]string, 1)
	s, err := detectObject(r)
	if err != nil {
		return m, err
	}
	buf := bytes.NewBuffer([]byte(s))
	rd := bufio.NewReader(buf)
	for {
		key, err := detect(rd, '"')
		if err != nil {
			if err == io.EOF {
				return m, err
			}
		}
		m[key], err = detectSimpleValue(rd)
		if err != nil {
			if err == io.EOF {
				return m, err
			}
		}
	}
}

// Parse2Dimensional calls Parse also on first child objects.
func Parse2Dimensional(r *bufio.Reader) (map[string]Value, error) {
	result := make(map[string]Value, 1)
	m, err := Parse(r)
	for k := range m {
		result[k] = Value{String: m[k]}
	}
	if err != nil && err != io.EOF {
		return result, err
	}
	for k := range m {
		if m[k][0] == braceOpen {
			b := bytes.NewBuffer([]byte(m[k]))
			rd := bufio.NewReader(b)
			mv, err := Parse(rd)
			if err != nil && err != io.EOF {
				return result, err
			}
			v := Value{String: result[k].String, Map: mv}
			result[k] = v
		}
	}
	return result, nil
}

func stripSpaces(str string) string {
	re := regexp.MustCompile(`\s`)
	return string(re.ReplaceAll([]byte(str), []byte("")))
}
