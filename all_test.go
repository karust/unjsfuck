package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

var (
	vm            = goja.New()
	jsunfuck, err = New("")
	test1         = `
(![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[+!+[]]+[!+[]+!+[]]+[!+[]+!+[]+!+[]]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]`
)

func TestFillMissingDigits(t *testing.T) {
	jsunfuck.fillMissingDigits()

	for i := 0; i < 10; i++ {
		want := fmt.Sprintf("%d", i)

		val, ok := MAPPING[want]
		if !ok {
			t.Fatalf("'%v' digit not present in the MAPPING", want)
		}

		result, err := vm.RunString(val)
		if err != nil {
			t.Fatalf("Failed to interpret the encoding: %v", err)
		}

		got := fmt.Sprintf("%v", result)
		if got != want {
			t.Fatalf("Wanted value '%v' mismatched to the resulted one '%v'", want, got)
		}
	}
}

func TestFillMissingChars(t *testing.T) {
	jsunfuck.fillMissingChars()

	for want, val := range MAPPING {
		if val == USE_CHAR_CODE {
			t.Fatalf("%v key value is unchanged", want)
		}

		result, err := vm.RunString(val)
		if err != nil {
			t.Fatalf("Failed to interpret the encoding: %v", err)
		}

		got := fmt.Sprintf("%v", result)
		if got != want {
			t.Fatalf("Wanted value '%v' mismatched to the resulted one '%v'", want, got)
		}
	}
}

func TestReplaceMap(t *testing.T) {
	contains := func(s []string, e string) bool {
		for _, v := range s {
			if strings.Contains(e, v) {
				return true
			}
		}
		return false
	}

	keys := []string{}
	for k := range CONSTRUCTORS {
		keys = append(keys, k)
	}
	for k := range SIMPLE {
		keys = append(keys, k)
	}

	jsunfuck.fillMissingDigits()
	//jsunfuck.fillMissingChars()

	mappingCache := make(map[string]string)
	for k, v := range MAPPING {
		mappingCache[k] = v
	}
	jsunfuck.replaceMap()

	for k, v := range MAPPING {
		if contains(keys, v) {
			t.Fatalf("%v. %v contains %v", k, v, keys)
		}

		if regexp.MustCompile(`(\d\d+)|\((\d)\)|\[(\d)\]|GLOBAL|\+""|""`).MatchString(v) {
			t.Fatalf("%v contains unmaped symbol", v)
		}
	}
}

func TestReplaceStrings(t *testing.T) {
	jsunfuck.fillMissingDigits()
	jsunfuck.fillMissingChars()
	jsunfuck.replaceMap()
	jsunfuck.replaceStrings()

	for _, v := range MAPPING {
		if found, _ := regexp.MatchString(`[^\[\]\(\)\!\+]`, v); found {
			t.Fatalf("`%v` contains non-encoded characters", v)
		}
	}
}

func TestInitMapChars(t *testing.T) {
	jsunfuck.Init()

	for want, v := range MAPPING {
		got, err := vm.RunString(v)
		if err != nil {
			t.Fatalf("Failed to interpret `%v` char encoding %v. Err: %v", want, v, err)
		}

		if got.String() != want {
			t.Fatalf("Got %v != Want %v. Err: %v", got, want, err)
		}
	}
}

func TestDecode(t *testing.T) {
	jsunfuck.Init()
	res := jsunfuck.Decode(`[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]][([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+([][[]]+[])[+!+[]]+(![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[])[+!+[]]+([][[]]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(!![]+[])[+!+[]]]((![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[+!+[]]+[!+[]+!+[]]+[!+[]+!+[]+!+[]]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]+([]+[])[(![]+[])[+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+([][[]]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(![]+[])[!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(!![]+[])[+!+[]]](+[![]]+([]+[])[(![]+[])[+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+([][[]]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(![]+[])[!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(!![]+[])[+!+[]]]()[+!+[]+[!+[]+!+[]]])[!+[]+!+[]+[+!+[]]])()`)
	fmt.Println(res, jsunfuck.uneval(res))
}

func TestEncode(t *testing.T) {
	// Test with encoded digits
	jsunfuck.fillMissingDigits()
	want := `(false+"")[1]+(false+"")[2]+(true+"")[3]+(true+"")[1]+(true+"")[0]+([]["flat"]+"")[13]+[+!+[]]+[!+[]+!+[]]+[!+[]+!+[]+!+[]]+([0]+false+[]["flat"])[20]`
	got := jsunfuck.Encode("alert(123)")
	if want != got {
		t.Fatalf("Test encoded digits. Want %v != Got %v", want, got)
	}
	_, err := vm.RunString(got)
	if err != nil {
		t.Fatalf("Failed to interpret the encoding: %v. Err: %v", got, err)
	}

	// Test with missing chars
	jsunfuck.fillMissingChars()
	want = `(false+"")[1]+(false+"")[2]+(true+"")[3]+(true+"")[1]+(true+"")[0]+([]["flat"]+"")[13]+Function("return unescape")()("%"+(27)+"")+Function("return unescape")()("%"+(4)+"a")+Function("return unescape")()("%"+(4)+"b")+Function("return unescape")()("%"+(4)+"c")+Function("return unescape")()("%"+(27)+"")+([0]+false+[]["flat"])[20]`
	got = jsunfuck.Encode("alert('JKL')")
	if want != got {
		t.Fatalf("Test encoded missing chars. Want %v != Got %v", want, got)
	}
	_, err = vm.RunString(got)
	if err != nil {
		t.Fatalf("Failed to interpret the encoding: %v. Err: %v", got, err)
	}

	// Test with mappings
	jsunfuck.replaceMap()
	want = `(![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([]["flat"]+[])[+!+[]+[!+[]+!+[]+!+[]]]+(+[![]]+[])[+[]]+([][[]]+[])[+[]]+((+[])["constructor"]+[])[+!+[]+[+!+[]]]+([]["entries"]()+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+([]["flat"]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[!![]]+[]+([+[]]+![]+[]["flat"])[!+[]+!+[]+[+[]]]+([+[]]+![]+[]["flat"])[!+[]+!+[]+[+[]]]`
	got = jsunfuck.Encode("alert(Number(true))")
	if want != got {
		t.Fatalf("Test encoded mappings. Want %v != Got %v", want, got)
	}

	// Test with encoded strings
	jsunfuck.replaceStrings()
	want = `(![]+[])[+!+[]]+(![]+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+(!![]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+(+[![]]+[])[+[]]+([][[]]+[])[+[]]+((+[])[([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+([][[]]+[])[+!+[]]+(![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[])[+!+[]]+([][[]]+[])[+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+[]]+(!![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[+!+[]+[+[]]]+(!![]+[])[+!+[]]]+[])[+!+[]+[+!+[]]]+([][(!![]+[])[!+[]+!+[]+!+[]]+([][[]]+[])[+!+[]]+(!![]+[])[+[]]+(!![]+[])[+!+[]]+([![]]+[][[]])[+!+[]+[+[]]]+(!![]+[])[!+[]+!+[]+!+[]]+(![]+[])[!+[]+!+[]+!+[]]]()+[])[!+[]+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+(!![]+[])[+!+[]]+([][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]]+[])[+!+[]+[!+[]+!+[]+!+[]]]+[!![]]+[]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]+([+[]]+![]+[][(![]+[])[+[]]+(![]+[])[!+[]+!+[]]+(![]+[])[+!+[]]+(!![]+[])[+[]]])[!+[]+!+[]+[+[]]]`
	got = jsunfuck.Encode("alert(Number(true))")
	if want != got {
		t.Fatalf("Test encoded mappings. Want %v != Got %v", want, got)
	}

	//(function(){var SOME=123;alert(SOME);}());
	res := jsunfuck.Encode(`alert(123);`)
	fmt.Println(res)
	//fmt.Println(jsunfuck.Wrap(res, true, true))
}
