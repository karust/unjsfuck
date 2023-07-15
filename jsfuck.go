package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type JSFuck struct {
	MAPPING map[string]string
}

func New() *JSFuck {
	jsFck := JSFuck{}
	jsFck.MAPPING = make(map[string]string)

	// Initialize local mapping from precalculated values
	for k, v := range DEFAULT_MAPPING {
		jsFck.MAPPING[k] = v
	}
	return &jsFck
}

func (jsFck *JSFuck) Init() {
	jsFck.fillMissingDigits()
	jsFck.fillMissingChars()
	jsFck.replaceMap()
	jsFck.replaceStrings()
}

func (jsFck *JSFuck) encodeDigit(digit int) string {
	output := `+[]` //0

	if digit > 0 {
		output = `+!` + output
	}

	for i := 0; i < digit-1; i++ {
		output = `+!+[]` + output
	}

	if digit > 1 {
		output = output[1:]
	}
	return output
}

// Calculates encoded values for 0-9 and adds it to MAPPING
func (jsFck *JSFuck) fillMissingDigits() {
	for number := 0; number <= 9; number++ {
		jsFck.MAPPING[strconv.Itoa(number)] = `[` + jsFck.encodeDigit(number) + `]`
	}
}

// Iterates over MAPPING and fills missing character values with a string
// containing their ASCII value represented in hex
func (jsFck *JSFuck) fillMissingChars() {
	digitRgxp := regexp.MustCompile(`\d+`)
	letterRgxp := regexp.MustCompile(`[^\d+]`)

	for key, value := range jsFck.MAPPING {
		if value == USE_CHAR_CODE {
			hexidec := hex.EncodeToString([]byte{key[0]})

			// Separate HEX digit from letter
			digitSearch := digitRgxp.FindString(hexidec)
			letterSearch := letterRgxp.FindString(hexidec)

			encodedString := fmt.Sprintf(`Function("return unescape")()("%%"+(%s)+"%s")`, digitSearch, letterSearch)
			jsFck.MAPPING[key] = encodedString
		}
	}
}

// Iterates over MAPPING from MIN to MAX and replaces value with values
// found in CONSTRUCTORS and SIMPLE, as well as using digitalReplacer and
// numberReplacer to replace numeric values
func (jsFck *JSFuck) replaceMap() {

	// Replaces found patterns in value with replacement
	replace := func(pattern, replacement, value string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(value, replacement)
	}

	// Replace found digit with encoded one
	digitReplacer := func(x string) string {
		re := regexp.MustCompile(`\d`)
		return jsFck.MAPPING[re.FindString(x)]
	}

	// Split number digits and encode
	numberReplacer := func(number string) string {
		digits := strings.Split(number, "")
		firstDigit, _ := strconv.Atoi(digits[0])

		// ?Keep implementation as in JSfuck
		encoded := jsFck.encodeDigit(firstDigit)
		concatenated := strings.Join(append([]string{encoded}, digits[1:]...), "+")

		re := regexp.MustCompile(`\d`)
		return re.ReplaceAllStringFunc(concatenated, digitReplacer)
	}

	// For every declared char
	for i := MIN; i <= MAX; i++ {
		char := fmt.Sprintf("%c", i)
		value, ok := jsFck.MAPPING[char]
		if !ok || value == "" {
			continue
		}

		// Replace all contructors till nothing's left
		original := ""
		for value != original {
			original = value
			for key, val := range CONSTRUCTORS {
				if !strings.Contains(value, key) {
					continue
				}
				value = replace(`\b`+key, val+`["constructor"]`, value)
			}

			for key, val := range SIMPLE {
				if !strings.Contains(value, key) {
					continue
				}
				value = replace(`\b`+key, val, value)
			}
		}

		// Replace all numbers and digits
		re := regexp.MustCompile(`(\d\d+)`)
		value = re.ReplaceAllStringFunc(value, numberReplacer)

		re = regexp.MustCompile(`\((\d)\)`)
		value = re.ReplaceAllStringFunc(value, digitReplacer)

		re = regexp.MustCompile(`\[(\d)\]`)
		value = re.ReplaceAllStringFunc(value, digitReplacer)

		// Replace GLOBAL and empty quotes
		value = replace(`GLOBAL`, GLOBAL, value)
		value = replace(`\+""`, `+[]`, value)
		value = replace(`""`, `[]+[]`, value)

		jsFck.MAPPING[char] = value
	}
}

// Replaces strings added in `replaceMap` with encoded there values
func (jsFck *JSFuck) replaceStrings() {
	// Split characters to correctly encode later
	mappingReplacer := func(b string) string {
		b = b[1 : len(b)-1] // Remove quotes, can't get groups here
		splitted := strings.Split(b, "")
		return strings.Join(splitted, "+")
	}

	// Replace content between quotes
	for key, value := range jsFck.MAPPING {
		re := regexp.MustCompile(`\"([^\"]+)\"`)
		jsFck.MAPPING[key] = re.ReplaceAllStringFunc(value, mappingReplacer)
	}

	pattern := `[^\[\]\(\)\!\+]`
	findNonEncodedRegexp := regexp.MustCompile(pattern)

	missing := make(map[string]string)
	valueReplacer := func(c string) string {
		if _, ok := missing[c]; ok {
			return c
		} else {
			return jsFck.MAPPING[c]
		}
	}

	found := true
	for found {
		found = false
		missing = make(map[string]string)

		for key, value := range jsFck.MAPPING {
			if findNonEncodedRegexp.MatchString(value) {
				missing[key] = value
				found = true
			}
		}

		for key := range missing {
			value := jsFck.MAPPING[key]
			value = findNonEncodedRegexp.ReplaceAllStringFunc(value, valueReplacer)

			jsFck.MAPPING[key] = value
			missing[key] = value
		}
	}
}

// Iterates over jsFck.MAPPING and replaces every value found with its corresponding key
func (jsFck *JSFuck) mapping(js string) string {
	// Reverse sort MAP, so bigger encodings first
	keys := make([]string, 0, len(jsFck.MAPPING))
	for key := range jsFck.MAPPING {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return len(jsFck.MAPPING[keys[i]]) > len(jsFck.MAPPING[keys[j]]) })

	// Replace values from big to small
	for _, key := range keys {
		value := jsFck.MAPPING[key]
		js = strings.ReplaceAll(js, value, key)
	}
	return js
}

// Unevals a piece of Javascript wrapped with an encoded eval
func (jsFck *JSFuck) uneval(js string) string {
	js = strings.ReplaceAll(js, `[][flat][constructor](`, ``)
	return js[:len(js)-3] // remove )()
}

// Decodes JSFuck`d Javascript
func (jsFck *JSFuck) Decode(js string) string {
	js = jsFck.mapping(js)

	//removes concatenation operators
	re := regexp.MustCompile(`\+`)
	js = re.ReplaceAllString(js, ``)
	js = strings.ReplaceAll(js, "++", "+")

	// check to see if source js is eval`d
	if strings.Contains(js, `[][fill][constructor]`) {
		js = jsFck.uneval(js)
	}

	return js
}

// Encodes plain Javascript to JSFuck obfuscated Javascript
func (jsFck *JSFuck) Encode(js string) string {

	output := []string{}
	regex := ""

	for key := range SIMPLE {
		regex += key + "|"
	}
	regex += "."

	inputReplacer := func(c string) string {
		if c == "\n" || c == "\r" {
			return ""
		}
		replacement, ok := SIMPLE[c]
		if ok {
			output = append(output, "["+replacement+"]+[]")
		} else {
			replacement, ok = jsFck.MAPPING[c]
			if ok {
				output = append(output, replacement)
			} else {
				replacement = fmt.Sprintf(
					"([]+[])[%v][%v](%v)",
					jsFck.Encode("constructor"), jsFck.Encode("fromCharCode"), jsFck.Encode(fmt.Sprintf("%c", c[0])),
				)

				output = append(output, replacement)
				jsFck.MAPPING[c] = replacement
			}
		}
		return replacement
	}

	re := regexp.MustCompile(regex)
	re.ReplaceAllStringFunc(js, inputReplacer)

	result := strings.Join(output, "+")
	if match, _ := regexp.MatchString(`^\d$`, js); match {
		result += "+[]"
	}

	return result
}

func (jsFck *JSFuck) Wrap(js string, isWrapEval, isParentScope bool) string {
	if isWrapEval {
		if isParentScope {
			js = "[][" + jsFck.Encode("flat") + "]" +
				"[" + jsFck.Encode("constructor") + "]" +
				"(" + jsFck.Encode("return eval") + ")()" +
				"(" + js + ")"
		} else {
			js = "[][" + jsFck.Encode("flat") + "]" +
				"[" + jsFck.Encode("constructor") + "]" +
				"(" + js + ")()"
		}
	}
	return js
}
