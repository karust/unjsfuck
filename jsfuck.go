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
	JS string
}

// Checks if passed some Javascript and if so assigns an instance variable
// to that of the pass Javascript.
// Populates MAPPING dictionary with the keys corresponding encoded value.
// Keyword arguments:
// js -- string containing the encoded Javascript to be
//
//	decoded (defualt None)
func New(js string) (*JSFuck, error) {

	jsFck := JSFuck{
		JS: js,
	}

	return &jsFck, nil
}

func (jsFck *JSFuck) Init() {
	jsFck.fillMissingDigits()
	jsFck.fillMissingChars()
	jsFck.replaceMap()
	jsFck.replaceStrings()
}

// Calculates 0-9`s encoded value and adds it to MAPPING
func (jsFck *JSFuck) fillMissingDigits() {
	for number := 0; number < 10; number++ {
		output := `+[]`

		if number > 0 {
			output = `+!` + output
		}

		for i := 0; i < number-1; i++ {
			output = `+!+[]` + output
		}

		if number > 1 {
			output = output[1:]
		}

		MAPPING[strconv.Itoa(number)] = `[` + output + `]`
	}
}

// Iterates over MAPPING and fills missing character values with a string
// containing their ascii value represented in hex
func (jsFck *JSFuck) fillMissingChars() {
	// TODO: Remove ?
	digitRgxp, _ := regexp.Compile(`\d+`)
	letterRgxp, _ := regexp.Compile(`[^\d+]`)

	for key, value := range MAPPING {
		if value == USE_CHAR_CODE {
			hexidec := hex.EncodeToString([]byte{key[0]})

			// TODO: Remove ?
			digitSearch := digitRgxp.FindString(hexidec)
			letterSearch := letterRgxp.FindString(hexidec)

			// digit, letter := "", ""
			// if digit != ""
			// digit = digit_search[0] if digit_search else ``
			// letter = letter_search[0] if letter_search else ``

			encodedString := fmt.Sprintf(`Function("return unescape")()("%%"+(%s)+"%s")`, digitSearch, letterSearch)
			MAPPING[key] = encodedString
		}
	}
}

// Iterates over MAPPING from MIN to MAX and replaces value with values
// found in CONSTRUCTORS and SIMPLE, as well as using digitalReplacer and
// numberReplacer to replace numeric values
func (jsFck *JSFuck) replaceMap() {
	replace := func(pattern, replacement, value string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(value, replacement)
	}

	digitReplacer := func(x string) string {
		re := regexp.MustCompile(`\d`)
		return MAPPING[re.FindString(x)]
	}

	numberReplacer := func(number string) string {
		values := strings.Split(number, "")
		head, _ := strconv.Atoi(values[0])
		output := `+[]`

		values = values[1:]

		if head > 0 {
			output = `+!` + output
		}

		for i := 1; i < head; i++ {
			output = `+!+[]` + output
		}

		if head > 1 {
			output = output[1:]
		}

		re := regexp.MustCompile(`(\d)`)
		return re.ReplaceAllStringFunc(strings.Join(append([]string{output}, values...), "+"), digitReplacer)
	}

	for i := MIN; i <= MAX; i++ {
		char := fmt.Sprintf("%c", i)
		value, ok := MAPPING[char]
		if !ok || value == "" {
			continue
		}

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

		re := regexp.MustCompile(`(\d\d+)`)
		value = re.ReplaceAllStringFunc(value, numberReplacer)

		re = regexp.MustCompile(`\((\d)\)`)
		value = re.ReplaceAllStringFunc(value, digitReplacer)

		re = regexp.MustCompile(`\[(\d)\]`)
		value = re.ReplaceAllStringFunc(value, digitReplacer)

		value = replace(`GLOBAL`, GLOBAL, value)
		value = replace(`\+""`, `+[]`, value)
		value = replace(`""`, `[]+[]`, value)

		MAPPING[char] = value
	}
}

// Replaces strings added in `replaceMap` with encoded there values
func (jsFck *JSFuck) replaceStrings() {

	// Function to split characters to correctly encode later
	mappingReplacer := func(b string) string {
		b = b[1 : len(b)-1] // Remove quotes, can't get groups here
		splitted := strings.Split(b, "")
		return strings.Join(splitted, "+")
	}

	// Replace content between quotes
	for key, value := range MAPPING {
		re := regexp.MustCompile(`\"([^\"]+)\"`)
		MAPPING[key] = re.ReplaceAllStringFunc(value, mappingReplacer)
	}

	pattern := `[^\[\]\(\)\!\+]`
	findNonEncodedRegexp := regexp.MustCompile(pattern)

	missing := make(map[string]string)
	valueReplacer := func(c string) string {
		if _, ok := missing[c]; ok {
			return c
		} else {
			return MAPPING[c]
		}
	}

	found := true
	for found {
		found = false
		missing = make(map[string]string)

		for key, value := range MAPPING {
			if findNonEncodedRegexp.MatchString(value) {
				missing[key] = value
				found = true
			}
		}

		for key := range missing {
			value := MAPPING[key]
			value = findNonEncodedRegexp.ReplaceAllStringFunc(value, valueReplacer)

			MAPPING[key] = value
			missing[key] = value
		}
	}
}

// Iterates over MAPPING and replaces every value found with
// its corresponding key
// Keyword arguments:
// js -- string containing Javascript encoded with JSFuck
// Returns:
// js -- string of decoded Javascript
func (jsFck *JSFuck) mapping(js string) string {
	// Reverse sort MAP
	keys := make([]string, 0, len(MAPPING))
	for key := range MAPPING {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return len(MAPPING[keys[i]]) > len(MAPPING[keys[j]]) })

	for _, key := range keys {
		value := MAPPING[key]
		//fmt.Printf("%v%v\n", key, value)
		js = strings.ReplaceAll(js, value, key)
	}
	return js
}

// Unevals a piece of Javascript wrapped with an encoded eval
// Keyword arguments:
// js -- string containing an eval wrapped string of Javascript
// Returns:
// js -- string with eval removed
func (jsFck *JSFuck) uneval(js string) string {
	js = strings.ReplaceAll(js, `[][flat][constructor](`, ``)
	return js[:len(js)-3] // remove )()
}

// Decodes JSFuck`d Javascript
// Keyword arguments:
// js -- string containing the JSFuck to be decoded (defualt None)
// Returns:
// js -- string of decoded Javascript
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

// Encodes vanilla Javascript to JSFuck obfuscated Javascript
// Keyword arguments:
// js                            -- string of unobfuscated Javascript
// wrapWithEval        -- boolean determines whether to wrap with an eval
// runInParentScope -- boolean determines whether to run in parents scope
func (jsFck *JSFuck) Encode(js string) string {

	output := []string{}
	regex := ""

	for key := range SIMPLE {
		regex += key + "|"
	}
	regex += "."

	inputReplacer := func(c string) string {
		replacement, ok := SIMPLE[c]
		if ok {
			output = append(output, "["+replacement+"]+[]")
		} else {
			replacement, ok = MAPPING[c]
			if ok {
				output = append(output, replacement)
			} else {
				replacement = fmt.Sprintf(
					"([]+[])[%v][%v](%v)",
					jsFck.Encode("constructor"), jsFck.Encode("fromCharCode"), jsFck.Encode(fmt.Sprintf("%c", c[0])),
				)

				output = append(output, replacement)
				MAPPING[c] = replacement
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

func (jsFck *JSFuck) Wrap(js string, wrapWithEval, runInParentScope bool) string {
	if wrapWithEval {
		if runInParentScope {
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
