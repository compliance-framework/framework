package policy_manager

import "testing"

var sshConfiguration = map[string][]string{
	"authorizedkeysfile": {
		".ssh/authorized_keys",
		".ssh/authorized_keys2",
	},
	"listenaddress": {
		"[::]:22",
		"0.0.0.0:22",
	},
	"passwordauthentication": {
		"yes",
	},
	"permitrootlogin": {
		//"without-password",
		"with-password",
	},
	"port": {
		"22",
	},
	"pubkeyauthentication": {
		"yes",
	},
}

func TestPolicyManager_New(t *testing.T) {
	t.Run("Policy Manager understands bundles", func(t *testing.T) {

	})
	t.Run("Policy Manager understands directories", func(t *testing.T) {

	})
}

//func TestSomething(t *testing.T) {
//	ctx := context.TODO()
//
//	r := rego.New(
//		rego.Query("data"),
//		rego.LoadBundle("./bundle_ssh.tar.gz"),
//		rego.Package("compliance_framework.local_ssh"),
//		rego.Input(sshConfiguration),
//	)
//
//	query, err := r.PrepareForEval(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, module := range query.Modules() {
//		// Exclude any test files for this compilation
//		if strings.HasSuffix(module.Package.Location.File, "_test.rego") {
//			continue
//		}
//
//		finalResult := Result{
//			Policy: Policy{
//				File:        module.Package.Location.File,
//				Package:     Package(module.Package.Path.String()),
//				Annotations: module.Annotations,
//			},
//			AdditionalVariables: map[string]interface{}{},
//			Violations:          nil,
//		}
//
//		sub := rego.New(
//			rego.Query(module.Package.Path.String()),
//			rego.LoadBundle("./bundle_ssh.tar.gz"),
//			rego.Package(module.Package.Path.String()),
//			rego.Input(sshConfiguration),
//		)
//
//		results, err := sub.Eval(ctx)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Println(results)
//
//		for _, result := range results {
//			for _, expression := range result.Expressions {
//				moduleOutputs := expression.Value.(map[string]interface{})
//
//				for key, value := range moduleOutputs {
//					if !slices.Contains([]string{"violation"}, key) {
//						finalResult.AdditionalVariables[key] = value
//					}
//				}
//
//				for _, tester := range moduleOutputs["violation"].([]interface{}) {
//					finalResult.Violations = append(finalResult.Violations, tester.(map[string]interface{}))
//				}
//
//				fmt.Println(finalResult.String())
//				//for _,violate := range violations["violation"] {
//				//
//				//}
//				//fmt.Println(expression)
//			}
//			//fmt.Println(result.Expressions)
//			//fmt.Println(result.Bindings)
//		}
//
//		//fmt.Println(result.([]map[string]interface{}))
//
//		//break
//		//fmt.Println(module.Package.Path)
//		fmt.Println("#########################################")
//		fmt.Println(finalResult.String())
//	}
//
//}

//func TestSomething(t *testing.T) {
//
//	ctx := context.TODO()
//
//	r := rego.New(
//		rego.Query("data.compliance_framework.local_ssh"),
//		rego.LoadBundle("./bundle_ssh.tar.gz"),
//		rego.Package("compliance_framework"),
//	)
//
//	query, err := r.PrepareForEval(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	results, err := query.Eval(ctx, rego.EvalInput(sshConfiguration))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for _, module := range query.Modules() {
//		fmt.Println(module.Package.Path.String())
//
//	}
//
//	var resultCompiled []struct {
//		Namespace           string
//		Violations          []map[string]interface{}
//		AdditionalVariables map[string]interface{}
//	}
//
//	fmt.Println(results)
//
//	for _, result := range results {
//		for _, expression := range result.Expressions {
//			for namespace, outcome := range expression.Value.(map[string]interface{}) {
//
//				fmt.Println("###########", outcome)
//
//				variables := outcome.(map[string]interface{})
//
//				var violations []map[string]interface{}
//
//				for _, violate := range variables["violation"].([]interface{}) {
//					violations = append(violations, violate.(map[string]interface{}))
//				}
//
//				//for name, value := range variables {
//				//
//				//}
//
//				result := struct {
//					Namespace           string
//					Violations          []map[string]interface{}
//					AdditionalVariables map[string]interface{}
//				}{
//					Namespace:           namespace,
//					Violations:          violations,
//					AdditionalVariables: variables,
//				}
//
//				resultCompiled = append(resultCompiled, result)
//
//				//fmt.Println()
//			}
//		}
//	}
//
//	fmt.Println(resultCompiled)
//	//fmt.Println(results)
//
//}
