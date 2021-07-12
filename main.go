package gojsonparser

import (
	"bufio"
	"bytes"
	"io"
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

func detectSimpleValue(r *bufio.Reader) (string, error) {
	var i int
	var tryParse bool
	var isDigit bool
	var isString bool
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
			return "", err
		}
		if rn == twoColon {
			i++
			tryParse = true
			continue
		}
		if unicode.IsSpace(rn) && isDigit {
			return string(digit), nil
		}
		if i > 0 && tryParse && i != lastIndexAppended {
			if unicode.IsDigit(rn) {
				isDigit = true
				digit = append(digit, rn)
				lastIndexAppended = i
			}
			if rn == semiColon && isDigit {
				return string(digit), nil
			}
			if unicode.IsLetter(rn) && !isDigit {
				isString = true
				str = append(str, rn)
			}
			if rn == semiColon && isString && !isDigit {
				return string(str), nil
			}
			i++
		}
	}
	return "", nil
}

func detectObject(r *bufio.Reader) (string, error) {
	var foundBefore bool
	var open int
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
			i++
			continue
		}
		if foundBefore && i > 0 {
			if rn == braceClose {
				open--
			}
			if rn == braceClose && open == 0 {
				break
			} else {
				s = append(s, rn)
			}
			i++
		}
	}
	return string(s), nil
}

func parseJSON(r *bufio.Reader) (map[string]string, error) {
	m := make(map[string]string, 1)
	s, err := detectObject(r)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte(s))
	rd := bufio.NewReader(buf)
	for {
		key, err := detect(rd, '"')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		m[key], err = detectSimpleValue(rd)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return m, nil
}
