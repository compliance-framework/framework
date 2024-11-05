package main

import (
	"container-solutions.com/continuous-compliance/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

//func main() {
//
//	ctx := context.TODO()
//
//	// Construct a Rego object that can be prepared or evaluated.
//	r := rego.New(
//		rego.Query(os.Args[2]),
//		rego.Load([]string{os.Args[1]}, nil))
//
//	// Create a prepared query that can be evaluated.
//	query, err := r.PrepareForEval(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Load the input document from stdin.
//	var input interface{}
//	dec := json.NewDecoder(os.Stdin)
//	dec.UseNumber()
//	if err := dec.Decode(&input); err != nil {
//		log.Fatal(err)
//	}
//
//	// Execute the prepared query.
//	rs, err := query.Eval(ctx, rego.EvalInput(input))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(rs.Allowed())
//
//	// Do something with the result.
//	fmt.Println(rs)
//}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "cf",
		Short: "cf manages policies for the compliance framework",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cf called")
			// Do Stuff Here
		},
	}

	rootCmd.AddCommand(cmd.VerifyCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
