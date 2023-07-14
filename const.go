package main

const (
	USE_CHAR_CODE = "USE_CHAR_CODE"
	MIN           = 32
	MAX           = 126
	GLOBAL        = `Function("return this")()`
	EVAL          = `[]["filter"]["constructor"](%v)()`
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
		`Function`: `[]["flat"]`,
		`RegExp`:   `Function("return/"+false+"/")()`,
		`Object`:   `[]["entries"]()`,
	}

	MAPPING = map[string]string{
		`a`: `(false+"")[1]`,
		`b`: `([]["entries"]()+"")[2]`,
		`c`: `([]["flat"]+"")[3]`,
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
		`o`: `(true+[]["flat"])[10]`,
		`p`: `(+(211))["to"+String["name"]](31)[1]`,
		`q`: `("")["fontcolor"]([0]+false+")[20]`,
		`r`: `(true+"")[1]`,
		`s`: `(false+"")[3]`,
		`t`: `(true+"")[0]`,
		`u`: `(undefined+"")[0]`,
		`v`: `(+(31))["to"+String["name"]](32)`,
		`w`: `(+(32))["to"+String["name"]](33)`,
		`x`: `(+(101))["to"+String["name"]](34)[1]`,
		`y`: `(NaN+[Infinity])[10]`,
		`z`: `(+(35))["to"+String["name"]](36)`,

		`A`:  `(NaN+[]["entries"]())[11]`,
		`B`:  `(+[]+Boolean)[10]`,
		`C`:  `Function("return escape")()(("")["italics"]())[2]`,
		`D`:  `Function("return escape")()([]["flat"])["slice"]("-1")`,
		`E`:  `(RegExp+"")[12]`,
		`F`:  `(+[]+Function)[10]`,
		`G`:  `(false+Function("return Date")()())[30]`,
		`H`:  USE_CHAR_CODE,
		`I`:  `(Infinity+"")[0]`,
		`J`:  USE_CHAR_CODE,
		`K`:  USE_CHAR_CODE,
		`L`:  USE_CHAR_CODE,
		`M`:  `(true+Function("return Date")()())[30]`,
		`N`:  `(NaN+"")[0]`,
		`O`:  `(+[]+Object)[10]`,
		`P`:  USE_CHAR_CODE,
		`Q`:  USE_CHAR_CODE,
		`R`:  `(+[]+RegExp)[10]`,
		`S`:  `(+[]+String)[10]`,
		`T`:  `(NaN+Function("return Date")()())[30]`,
		`U`:  `(NaN+Object()["to"+String["name"]]["call"]())[11]`,
		`V`:  USE_CHAR_CODE,
		`W`:  USE_CHAR_CODE,
		`X`:  USE_CHAR_CODE,
		`Y`:  USE_CHAR_CODE,
		`Z`:  USE_CHAR_CODE,
		` `:  `(NaN+[]["flat"])[11]`,
		`!`:  USE_CHAR_CODE,
		`"`:  `("")["fontcolor"]()[12]`,
		`#`:  USE_CHAR_CODE,
		`$`:  USE_CHAR_CODE,
		`%`:  `Function("return escape")()([]["flat"])[21]`,
		`&`:  `("")["fontcolor"](")[13]`,
		`'`:  USE_CHAR_CODE,
		`(`:  `([]["flat"]+"")[13]`,
		`)`:  `([0]+false+[]["flat"])[20]`,
		`*`:  USE_CHAR_CODE,
		`+`:  `(+(+!+[]+(!+[]+[])[!+[]+!+[]+!+[]]+[+!+[]]+[+[]]+[+[]])+[])[2]`,
		`,`:  `[[]]["concat"]([[]])+""`,
		`-`:  `(+(.+[0000001])+"")[2]`,
		`.`:  `(+(+!+[]+[+!+[]]+(!![]+[])[!+[]+!+[]+!+[]]+[!+[]+!+[]]+[+[]])+[])[+!+[]]`,
		`/`:  `(false+[0])["italics"]()[10]`,
		`:`:  `(RegExp()+"")[3]`,
		`;`:  `("")["fontcolor"](NaN+")[21]`,
		`<`:  `("")["italics"]()[0]`,
		`=`:  `("")["fontcolor"]()[11]`,
		`>`:  `("")["italics"]()[2]`,
		`?`:  `(RegExp()+"")[2]`,
		`@`:  USE_CHAR_CODE,
		`[`:  `([]["entries"]()+"")[0]`,
		"\\": `(RegExp("/")+"")[1]`,
		`]`:  `([]["entries"]()+"")[22]`,
		`^`:  USE_CHAR_CODE,
		`_`:  USE_CHAR_CODE,
		"`":  USE_CHAR_CODE,
		`{`:  `(true+[]["flat"])[20]`,
		`|`:  USE_CHAR_CODE,
		`}`:  `([]["flat"]+"")["slice"]("-1")`,
		`~`:  USE_CHAR_CODE,
	}
)
