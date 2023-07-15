package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsFuck        = New()
	isEval        = false
	isParentScope = false
)

var rootCMD = &cobra.Command{
	Use:          "unjsfuck",
	Short:        "UnJSFuck",
	Long:         `Encode/Decode JSFuck (0.5.0) obfuscated Javascript`,
	Version:      "1.0.0",
	SilenceUsage: true,
}

var encodeCMD = &cobra.Command{
	Use:   "encode",
	Short: "Encode Javascript code",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Errof reading file: %v", err)
			return
		}

		js := string(data)
		encodedJS := jsFuck.Encode(js)
		wrappedJS := jsFuck.Wrap(encodedJS, isEval, isParentScope)
		fmt.Println(wrappedJS)
	},
}

var decodeCMD = &cobra.Command{
	Use:   "decode",
	Short: "Decode JSfuck obfuscated Javascript code",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs, cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Errof reading file: %v", err)
			return
		}

		encodedJS := string(data)
		decodedJS := jsFuck.Decode(encodedJS)
		fmt.Println(decodedJS)
	},
}

func init() {
	rootCMD.PersistentFlags().BoolVarP(&isEval, "eval", "e", false, "Wrap/Unwrap JS code in eval()")
	rootCMD.PersistentFlags().BoolVarP(&isParentScope, "parent", "p", false, "Run JS in parent scope")
	rootCMD.AddCommand(encodeCMD)
	rootCMD.AddCommand(decodeCMD)
}

func main() {
	jsFuck.Init()

	if err := rootCMD.Execute(); err != nil {
		fmt.Printf("Errof occurred during the execution: %v", err)
		os.Exit(1)
	}
}
