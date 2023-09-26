package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apparentlymart/terraform-provider/tfprovider"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <provider-executable> [provider-args...]\n", args[0])
		os.Exit(1)
	}
	args = args[1:]

	ctx := context.Background()
	provider, err := tfprovider.Start(ctx, args[0], args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer provider.Close()

	schema, diags := provider.Schema(ctx)
	showDiagnosticsMaybeExit(diags, provider)

	if len(schema.ManagedResourceTypes) != 0 {
		fmt.Print("\n# Managed Resource Types\n\n")
		for name := range schema.ManagedResourceTypes {
			fmt.Printf("- %s\n", name)
		}
	}

	if len(schema.DataResourceTypes) != 0 {
		fmt.Print("\n# Data Resource Types\n\n")
		for name := range schema.DataResourceTypes {
			fmt.Printf("- %s\n", name)
		}
	}

	// hashicups_coffees
	var input = make(map[string]cty.Value)
	readDataSource(provider, ctx, "hashicups_coffees", input)

	// hashicups_order
	// For hashicups_order and hashicups_ingredients we need to configure provider
	input = make(map[string]cty.Value)
	input["username"] = cty.StringVal("education")
	input["password"] = cty.StringVal("test123")
	configureProvider(provider, ctx, input)
	input = make(map[string]cty.Value)
	input["id"] = cty.NumberIntVal(2)
	readDataSource(provider, ctx, "hashicups_order", input)

	fmt.Print("\n")
}

func showDiagnosticsMaybeExit(diags tfprovider.Diagnostics, provider tfprovider.Provider) {
	for _, diag := range diags {
		switch diag.Severity {
		case tfprovider.Error:
			fmt.Fprintf(os.Stderr, "Error: %s; %s", diag.Summary, diag.Detail)
		case tfprovider.Warning:
			fmt.Fprintf(os.Stderr, "Warning: %s; %s", diag.Summary, diag.Detail)
		default:
			fmt.Fprintf(os.Stderr, "???: %s; %s", diag.Summary, diag.Detail)
		}
	}
	if diags.HasErrors() {
		provider.Close()
		os.Exit(1)
	}
}

func readDataSource(provider tfprovider.Provider, ctx context.Context, datasourceName string, input map[string]cty.Value) {
	schema, diags := provider.Schema(ctx)
	showDiagnosticsMaybeExit(diags, provider)

	ds := schema.DataResourceTypes[datasourceName]
	var config = make(map[string]cty.Value)
	for name, attr := range ds.Content.Attributes {
		config[name] = attr.EmptyValue()
	}
	for name, attr := range ds.Content.BlockTypes {
		config[name] = attr.EmptyValue()
	}
	for name, val := range input {
		config[name] = val
	}
	drs := provider.DataResourceType(datasourceName)

	conf := cty.ObjectVal(config)
	req := tfprovider.DataResourceReadRequest{Config: conf}
	coffees, diags := drs.Read(ctx, req)
	showDiagnosticsMaybeExit(diags, provider)
	r, _ := ctyjson.Marshal(coffees.State, coffees.State.Type())
	fmt.Printf("Result for %s is %s", datasourceName, string(r))
}

func configureProvider(provider tfprovider.Provider, ctx context.Context, input map[string]cty.Value) {
	schema, diags := provider.Schema(ctx)
	showDiagnosticsMaybeExit(diags, provider)

	var config = make(map[string]cty.Value)
	for name, attr := range schema.ProviderConfig.Attributes {
		config[name] = attr.EmptyValue()
	}
	for name, attr := range schema.ProviderConfig.BlockTypes {
		config[name] = attr.EmptyValue()
	}
	for name, val := range input {
		config[name] = val
	}
	pc := tfprovider.Config{
		Value: cty.ObjectVal(config),
	}
	diags = provider.Configure(ctx, pc)
	fmt.Println("DIAGS")
	for _, d := range diags {
		fmt.Println(d.Detail)
	}
}
