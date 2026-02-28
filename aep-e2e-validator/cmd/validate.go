package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aep-dev/aep-e2e-validator/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	configPath     string
	collection     string
	allCollections bool
	testNames      []string
	headerFlags    []string
)

func parseHeaders(raw []string) ([]validator.Header, error) {
	headers := make([]validator.Header, 0, len(raw))
	for _, h := range raw {
		key, value, ok := strings.Cut(h, "=")
		if !ok {
			return nil, fmt.Errorf("invalid header %q: must be in key=value format", h)
		}
		headers = append(headers, validator.Header{Key: strings.TrimSpace(key), Value: strings.TrimSpace(value)})
	}
	return headers, nil
}

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
		if collection != "" && allCollections {
			return fmt.Errorf("cannot specify both collection and all-collections")
		}

		headers, err := parseHeaders(headerFlags)
		if err != nil {
			return err
		}

		fmt.Printf("Validating with config: %s\n", configPath)
		if allCollections {
			fmt.Println("Validating all collections")
		} else {
			fmt.Printf("Validating collection: %s\n", collection)
		}

		v := validator.NewValidator(configPath, collection, allCollections, testNames, headers)
		exitCode := v.Run()
		if exitCode != validator.ExitCodeSuccess {
			os.Exit(exitCode)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&configPath, "config", "", "Path to the OpenAPI spec file")
	validateCmd.Flags().StringVar(&collection, "collection", "", "Name of the collection to validate")
	validateCmd.Flags().BoolVar(&allCollections, "all-collections", false, "Validate all collections in the spec")
	validateCmd.Flags().StringSliceVar(&testNames, "tests", []string{}, "Comma-separated list of tests to run (e.g. aep-133-create)")
	validateCmd.Flags().StringArrayVarP(&headerFlags, "header", "H", []string{}, "Headers to include in every request (format: key=value, repeatable)")

	validateCmd.MarkFlagRequired("config")
}
