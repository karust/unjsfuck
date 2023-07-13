package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	USE_CHAR_CODE = "USE_CHAR_CODE"
	MIN           = 32
	MAX           = 126
	GLOBAL        = `Function("return this")()`
)

var (
	SIMPLE = map[string]string{
		`false`:     `![]`,
		`true`:      `!![]`,
		`undefined`: `[][[]]`,
		`NaN`:       `+[![]]`,
		`Infinity`:  `+(+!+[]+(!+[]+[])[!+[]+!+[]+!+[]]+[+!+[]]+[+[]]+[+[]]+[+[]])`, // +"1e1000"
	}

	CONSTRUCTORS = map[string]string{
		`Array`:    `[]`,
		`Number`:   `(+[])`,
		`String`:   `([]+[])`,
		`Boolean`:  `(![])`,
		`Function`: `[]["fill"]`,
		`RegExp`:   `[]["fill"]("return/"+false+"/")()`,
	}

	MAPPING = map[string]string{
		`a`: `(false+"")[1]`,
		`b`: `([]["entries"]()+"")[2]`,
		`c`: `([]["fill"]+"")[3]`,
		`d`: `(undefined+"")[2]`,
		`e`: `(true+"")[3]`,
		`f`: `(false+"")[0]`,
		`g`: `(false+[0]+String)[20]`,
		`h`: `(+(101))["to"+String["name"]](21)[1]`,
		`i`: `([false]+undefined)[10]`,
		`j`: `([]["entries"]()+"")[3]`,
		`k`: `(+(20))["to"+String["name"]](21)`,
		`l`: `(false+"")[2]`,
		`m`: `(Number+"")[11]`,
		`n`: `(undefined+"")[1]`,
		`o`: `(true+[]["fill"])[10]`,
		`p`: `(+(211))["to"+String["name"]](31)[1]`,
		`q`: `(+(212))["to"+String["name"]](31)[1]`,
		`r`: `(true+"")[1]`,
		`s`: `(false+"")[3]`,
		`t`: `(true+"")[0]`,
		`u`: `(undefined+"")[0]`,
		`v`: `(+(31))["to"+String["name"]](32)`,
		`w`: `(+(32))["to"+String["name"]](33)`,
		`x`: `(+(101))["to"+String["name"]](34)[1]`,
		`y`: `(NaN+[Infinity])[10]`,
		`z`: `(+(35))["to"+String["name"]](36)`,
		`A`: `(+[]+Array)[10]`,
		`B`: `(+[]+Boolean)[10]`,
		`C`: `Function("return escape")()(("")["italics"]())[2]`,
		`D`: `Function("return escape")()([]["fill"])["slice"]("-1")`,
		`E`: `(RegExp+"")[12]`,
		`F`: `(+[]+Function)[10]`,
		`G`: `(false+Function("return Date")()())[30]`,
		`H`: USE_CHAR_CODE,
		`I`: `(Infinity+"")[0]`,
		`J`: USE_CHAR_CODE,
		`K`: USE_CHAR_CODE,
		`L`: USE_CHAR_CODE,
		`M`: `(true+Function("return Date")()())[30]`,
		`N`: `(NaN+"")[0]`,
		`O`: `(NaN+Function("return{}")())[11]`,
		`P`: USE_CHAR_CODE,
		`Q`: USE_CHAR_CODE,
		`R`: `(+[]+RegExp)[10]`,
		`S`: `(+[]+String)[10]`,
		`T`: `(NaN+Function("return Date")()())[30]`,
		`U`: `(NaN+Function("return{}")()["to"+String["name"]]["call"]())[11]`,
		`V`: USE_CHAR_CODE,
		`W`: USE_CHAR_CODE,
		`X`: USE_CHAR_CODE,
		`Y`: USE_CHAR_CODE,
		`Z`: USE_CHAR_CODE,
		` `: `(NaN+[]["fill"])[11]`,
		`!`: USE_CHAR_CODE,
		`"`: `("")["fontcolor"]()[12]`,
		`#`: USE_CHAR_CODE,
		`$`: USE_CHAR_CODE,
		`%`: `Function("return escape")()([]["fill"])[21]`,
		`&`: `("")["link"](0+")[10]`,
		`\`: USE_CHAR_CODE,
		`(`: `(undefined+[]["fill"])[22]`,
		`)`: `([0]+false+[]["fill"])[20]`,
		`*`: USE_CHAR_CODE,
		`+`: `(+(+!+[]+(!+[]+[])[!+[]+!+[]+!+[]]+[+!+[]]+[+[]]+[+[]])+[])[2]`,
		`,`: `([]["slice"]["call"](false+"")+"")[1]`,
		`-`: `(+(.+[0000000001])+"")[2]`,
		`.`: `(+(+!+[]+[+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+[!+[]+!+[]]+[+[]])+[])[+!+[]]`,
		`/`: `(false+[0])["italics"]()[10]`,
		`:`: `(RegExp()+"")[3]`,
		`;`: `("")["link"](")[14]`,
		`<`: `("")["italics"]()[0]`,
		`=`: `("")["fontcolor"]()[11]`,
		`>`: `("")["italics"]()[2]`,
		`?`: `(RegExp()+"")[2]`,
		`@`: USE_CHAR_CODE,
		`[`: `([]["entries"]()+"")[0]`,
		//`\\`: USE_CHAR_CODE,
		`]`: `([]["entries"]()+"")[22]`,
		`^`: USE_CHAR_CODE,
		`_`: USE_CHAR_CODE,
		`'`: USE_CHAR_CODE,
		`{`: `(true+[]["fill"])[20]`,
		`|`: USE_CHAR_CODE,
		`}`: `([]["fill"]+"")["slice"]("-1")`,
		`~`: USE_CHAR_CODE,
	}
)

type JSUnfuckIt struct {
	JS string
}

// Checks if passed some Javascript and if so assigns an instance variable
// to that of the pass Javascript.
// Populates MAPPING dictionary with the keys corresponding encoded value.
// Keyword arguments:
// js -- string containing the encoded Javascript to be
//
//	decoded (defualt None)
func New(js string) (*JSUnfuckIt, error) {

	jsFck := JSUnfuckIt{
		JS: js,
	}

	return &jsFck, nil
}

func (jsFck *JSUnfuckIt) Init() {
	jsFck.fillMissingDigits()
	//jsFck.fillMissingChars()
	jsFck.replaceMap()
	jsFck.replaceStrings()
}

// Calculates 0-9's encoded value and adds it to MAPPING
func (jsFck *JSUnfuckIt) fillMissingDigits() {
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
func (jsFck *JSUnfuckIt) fillMissingChars() {
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
			// digit = digit_search[0] if digit_search else ''
			// letter = letter_search[0] if letter_search else ''

			encodedString := fmt.Sprintf(`Function("return unescape")()("%%"+(%s)+"%s")`, digitSearch, letterSearch)
			MAPPING[key] = encodedString
		}
	}
}

// Iterates over MAPPING from MIN to MAX and replaces value with values
// found in CONSTRUCTORS and SIMPLE, as well as using digitalReplacer and
// numberReplacer to replace numeric values
func (jsFck *JSUnfuckIt) replaceMap() {
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

// Replaces strings added in __replaceMap with there encoded values
func (jsFck *JSUnfuckIt) replaceStrings() {
	pattern := `[^\[\]\(\)\!\+]`
	missing := make(map[string]string)

	// determines if there are still characters to replace
	findMissing := func() bool {
		done := false

		missing = make(map[string]string)

		for key, value := range MAPPING {
			re := regexp.MustCompile(pattern)
			matches := re.FindAllString(value, -1)

			if value != "" && len(matches) > 0 {
				missing[key] = value
				done = true
			}
		}

		return done
	}

	mappingReplacer := func(b string) string {
		splitted := strings.Split(b, "")
		joined := strings.Join(splitted, "+")
		return joined
	}

	valueReplacer := func(c string) string {
		fmt.Println("VAL:", c, missing[c], MAPPING[c])
		if val, ok := missing[c]; ok && val != "" {
			return c
		} else {
			return MAPPING[c]
		}
	}

	for key, value := range MAPPING {
		re := regexp.MustCompile(`\"([^\"]+)\"`)
		MAPPING[key] = re.ReplaceAllStringFunc(value, mappingReplacer)
	}

	for findMissing() {
		for key, value := range missing {
			re := regexp.MustCompile(pattern)
			value = re.ReplaceAllStringFunc(value, valueReplacer)

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
func (jsFck *JSUnfuckIt) mapping(js string) string {
	for key, value := range MAPPING {
		js = strings.ReplaceAll(js, value, key)
	}
	return js
}

// Unevals a piece of Javascript wrapped with an encoded eval

// Keyword arguments:
// js -- string containing an eval wrapped string of Javascript

// Returns:
// js -- string with eval removed
func (jsFck *JSUnfuckIt) uneval(js string) string {
	js = strings.ReplaceAll(js, `[][fill][constructor](`, ``)

	ev := `return eval)()(`

	if strings.Contains(js, ev) {
		js = js[(strings.Index(js, ev) + len(ev)):]
	}
	return js
}

// Decodes JSFuck'd Javascript
// Keyword arguments:
// js -- string containing the JSFuck to be decoded (defualt None)
// Returns:
// js -- string of decoded Javascript
func (jsFck *JSUnfuckIt) Decode(js string) string {
	js = jsFck.mapping(js)

	// removes concatenation operators
	// re := regexp.MustCompile(`\+`)
	// js = re.ReplaceAllString(js, ``)
	// js = strings.ReplaceAll(js, "++", "+")

	// // check to see if source js is eval'd
	// if strings.Contains(js, `[][fill][constructor]`) {
	// 	js = jsFck.uneval(js)
	// }

	return js
}
