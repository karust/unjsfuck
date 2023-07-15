# UnJSFuck
[![Go Report Card](https://goreportcard.com/badge/github.com/karust/openserp)](https://goreportcard.com/report/github.com/karust/openserp)
<a href="https://github.com/karust/unjsfuck/actions"><img src="https://github.com/karust/unjsfuck/actions/workflows/build_tests.yml/badge.svg"/></a>
[![codecov](https://codecov.io/gh/karust/unjsfuck/branch/main/graph/badge.svg?token=WJRP98YJCW)](https://codecov.io/gh/karust/unjsfuck)

Encode/Decode [JSFuck](https://github.com/aemkei/jsfuck/) (0.5.0) obfuscated Javascript.

Helpful resources:
* https://jsfuck.com/ - test encoding (results may differ)
* https://enkhee-osiris.github.io/Decoder-JSFuck/ - test decoding

## Usage
Use latest release [binary](https://github.com/karust/unjsfuck/releases) or install the tool with:
```sh
go install github.com/karust/unjsfuck
```

### Encode
```sh
unjsfuck encode ./test/test_plain.js
```

### Decode
```sh
unjsfuck decode ./test/test_enc.js
```

### Test
```sh
go test . -v
```


## Package usage
### Install
```sh
go get github.com/karust/unjsfuck
```
### Decode
```go
yourEncodedJS := "..."

jsFuck := New()
jsFuck.Init()
fmt.Println(jFuck.Decode(yourEncodedJS))
```
### Encode
```go
yourPlainJS := "alert(123);"

jsFuck := New()
jsFuck.Init()

encoded := jsFuck.Encode(yourEncodedJS)

// Wrap in eval and parent scope execution
wrapped := jsFuck.Wrap(true, true) 
fmt.Println(wrapped)
```