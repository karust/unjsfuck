package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

var vm = otto.New()

func contains(s []string, e string) bool {
	for _, v := range s {
		if strings.Contains(e, v) {
			return true
		}
	}
	return false
}

func TestFillMissingDigits(t *testing.T) {
	jsFuck := New()
	jsFuck.fillMissingDigits()

	for i := 0; i <= 9; i++ {
		want := strconv.Itoa(i)

		encoding, ok := jsFuck.MAPPING[want]
		if !ok {
			t.Fatalf("'%v' digit not present in the MAPPING", want)
		}

		result, err := vm.Run(encoding)
		if err != nil {
			t.Fatalf("Failed to interpret the encoding: %v", err)
		}

		got, _ := result.ToString()
		if got != want {
			t.Fatalf("Wanted value '%v' mismatched to the resulted one '%v'", want, got)
		}
	}
}

func TestFillMissingChars(t *testing.T) {
	jsFuck := New()
	jsFuck.fillMissingDigits()
	jsFuck.fillMissingChars()

	for want, encoding := range jsFuck.MAPPING {
		if encoding == USE_CHAR_CODE {
			t.Fatalf("%v key value is unchanged", want)
		}

		// // Interpreter can't check this
		// result, err := vm.Run(encoding)
		// if err != nil {
		// 	t.Fatalf("Failed to interpret the encoding %v: Err %v", encoding, err)
		// }
		// got := fmt.Sprintf("%v", result)
		// if got != want {
		// 	t.Fatalf("Wanted value '%v' mismatched to the resulted one '%v'", want, got)
		// }
	}
}

func TestReplaceMap(t *testing.T) {
	jsFuck := New()

	// Collect plain JS constructions
	constructions := []string{}
	for k := range CONSTRUCTORS {
		constructions = append(constructions, k)
	}
	for k := range SIMPLE {
		constructions = append(constructions, k)
	}

	jsFuck.fillMissingDigits()
	jsFuck.fillMissingChars()
	jsFuck.replaceMap()

	for plain, encoded := range jsFuck.MAPPING {
		if contains(constructions, encoded) {
			t.Fatalf("`%v` - `%v` contains %v", plain, encoded, constructions)
		}

		// If there are any digits, GLOBAL or empty quotes left
		if regexp.MustCompile(`(\d\d+)|\((\d)\)|\[(\d)\]|GLOBAL|\+""|""`).MatchString(encoded) {
			t.Fatalf("%v contains unmaped symbol", encoded)
		}
	}
}

func TestReplaceStrings(t *testing.T) {
	jsFuck := New()
	jsFuck.fillMissingDigits()
	jsFuck.fillMissingChars()
	jsFuck.replaceMap()
	jsFuck.replaceStrings()

	for _, v := range jsFuck.MAPPING {
		if found, _ := regexp.MatchString(`[^\[\]\(\)\!\+]`, v); found {
			t.Fatalf("`%v` contains non-encoded characters", v)
		}
	}
}

func TestSteppedEncode(t *testing.T) {
	jsFuck := New()
	// Test with encoded digits
	jsFuck.fillMissingDigits()
	want := `(false+"")[1]+(false+"")[2]+(true+"")[3]+(true+"")[1]+(true+"")[0]+([]["flat"]+"")[13]+[+!+[]]+[!+[]+!+[]]+[!+[]+!+[]+!+[]]+([0]+false+[]["flat"])[20]`
	got := jsFuck.Encode("alert(123)")
	if want != got {
		t.Fatalf("Test encoded digits. Want %v != Got %v", want, got)
	}
	_, err := vm.Run(got)
	if err != nil {
		t.Fatalf("Failed to interpret the encoding: %v. Err: %v", got, err)
	}

	// Test with missing chars
	jsFuck.fillMissingChars()
	want = `(false+"")[1]+(false+"")[2]+(true+"")[3]+(true+"")[1]+(true+"")[0]+([]["flat"]+"")[13]+Function("return unescape")()("%"+(27)+"")+Function("return unescape")()("%"+(4)+"a")+Function("return unescape")()("%"+(4)+"b")+Function("return unescape")()("%"+(4)+"c")+Function("return unescape")()("%"+(27)+"")+([0]+false+[]["flat"])[20]`
	got = jsFuck.Encode("alert('JKL')")
	if want != got {
		t.Fatalf("Test encoded missing chars. Want %v != Got %v", want, got)
	}
	_, err = vm.Run(got)
	if err != nil {
		t.Fatalf("Failed to interpret the encoding: %v. Err: %v", got, err)
	}

	// Test with mappings
	jsFuck.replaceMap()
	want = `(![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([]["flat"]+[])[+!+[]+[!+[]+!+[]+!+[]]]+(+[![]]+[])[+[]]+([][[]]+[])[+[]]+((+[])["constructor"]+[])[+!+[]+[+!+[]]]+([]["entries"]()+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+([]["flat"]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[!![]]+[]+([+[]]+![]+[]["flat"])[!+[]+!+[]+[+[]]]+([+[]]+![]+[]["flat"])[!+[]+!+[]+[+[]]]`
	got = jsFuck.Encode("alert(Number(true))")
	if want != got {
		t.Fatalf("Test encoded mappings. Want %v != Got %v", want, got)
	}

	// Test with encoded strings
	jsFuck.replaceStrings()
	want = `(![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+(+[![]]+[])[+[]]+([][[]]+[])[+[]]+((+[])[([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+([][[]]+[])[+!+[]]+(![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[])[+!+[]]+([][[]]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(!![]+[])[+!+[]]]+[])[+!+[]+[+!+[]]]+([][(!![]+[])[!+[]+!+[]+!+[]]+([][[]]+[])[+!+[]]+(!![]+[])[+[]]+(!![]+[])[+!+[]]+([![]]+[][[]])[+!+[]+[+[]]]+(!![]+[])[!+[]+!+[]+!+[]]+(![]+[])[!+[]+!+[]+!+[]]]()+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[!![]]+[]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]`
	got = jsFuck.Encode("alert(Number(true))")
	if want != got {
		t.Fatalf("Test encoded mappings. Want %v != Got %v", want, got)
	}
}

func TestEncode(t *testing.T) {
	jsFuck := New()
	jsFuck.Init()

	want, _ := os.ReadFile("./test/test_encoded.js")
	plain, _ := os.ReadFile("./test/test_plain.js")

	encoded := jsFuck.Encode(string(plain))
	got := jsFuck.Wrap(encoded, true, false)

	if got != string(want) {
		t.Fatalf("Wanted value != Got: %v", got)
	}
}

func TestDecode(t *testing.T) {
	jsFuck := New()
	jsFuck.Init()

	want := `(    function(){        var Some="Hallo Welt!";        alert(Some);    }());`
	data, _ := os.ReadFile("./test/test_encoded.js")
	got := jsFuck.Decode(string(data))

	if got != want {
		t.Fatalf("Wanted value: `%v` != Got: `%v`", want, got)
	}
}
