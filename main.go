package main

import (
	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "ht",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewVariableLocationRule(),
				rules.NewVariableOrderRule(),
				rules.NewVariableFieldOrderRule(),
				rules.NewKeyAttributesRule(),
			},
		},
	})
}
