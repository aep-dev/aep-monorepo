package cmd

import (
	"fmt"

	"github.com/aep-dev/aep-e2e-validator/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	configPath    string
	collection    string
	allCollections bool
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an AEP API",
	Long:  `Run end-to-end validation against an AEP API defined by an OpenAPI spec.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			return fmt.Errorf("config path is required")
		}
		if collection == "" && !allCollections {
			return fmt.Errorf("either collection or all-collections must be specified")
		}
		fmt.Printf("Validating with config: %s\n", configPath)
		if allCollections {
			fmt.Println("Validating all collections")
		} else {
			fmt.Printf("Validating collection: %s\n", collection)
		}

		v := validator.NewValidator(configPath, collection, allCollections)
		return v.Run()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&configPath, "config", "", "Path to the OpenAPI spec file")
	validateCmd.Flags().StringVar(&collection, "collection", "", "Name of the collection to validate")
	validateCmd.Flags().BoolVar(&allCollections, "all-collections", false, "Validate all collections in the spec")

	validateCmd.MarkFlagRequired("config")
}
