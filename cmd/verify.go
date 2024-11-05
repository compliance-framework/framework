package cmd

import (
	"container-solutions.com/continuous-compliance/internal"
	"context"
	"github.com/spf13/cobra"
	"log"
	"maps"
	"os"
	"slices"
)

func VerifyCmd() *cobra.Command {
	var VerifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify policies for cf",
		Long: `Policies in Compliance Framework require some metadata such as controls or policy_ids. 
Verify ensures that all required metadata is set before the policies are uploaded to an OCI registry`,
		Run: VerifyPolicies,
	}

	var PolicyPath string
	VerifyCmd.Flags().StringVarP(&PolicyPath, "policy-path", "p", "", "Directory where policies are stored")
	err := VerifyCmd.MarkFlagRequired("policy-path")
	internal.OnError(err, func(err error) { panic(err) })
	err = VerifyCmd.MarkFlagDirname("policy-path")
	internal.OnError(err, func(err error) { panic(err) })

	return VerifyCmd
}

func VerifyPolicies(cmd *cobra.Command, args []string) {
	policyPath, err := cmd.Flags().GetString("policy-path")
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(policyPath)
	internal.OnError(err, func(err error) {
		if os.IsNotExist(err) {
			log.Fatal("Policy path does not exist at specified policy-path")
		}
		log.Fatal(err)
	})

	ctx := context.TODO()
	compiler := internal.PolicyCompiler(ctx, policyPath)

	valid := true
	for _, module := range compiler.Modules {
		annotations := internal.ExtractAnnotations(module.Comments)
		if annotations["cf_enabled"] == nil {
		}
		if _, exists := annotations["cf_enabled"]; !exists {
			continue
		}
		if annotations["cf_enabled"] != "true" {
			continue
		}
		missingAnnotations := internal.SubtractSlice(internal.RequiredAnnotations, slices.Collect(maps.Keys(annotations)))
		if len(missingAnnotations) > 0 {
			log.Println(module.Package.Location.File, "is missing required annotations", missingAnnotations)
			valid = false
		}
	}

	if !valid {
		log.Fatal("Validation for Compliance Framework Policies failed")
	}

	log.Print("Validation for Compliance Framework Policies successful")
}
