package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area/areainfo"
)

var genCmd = &cobra.Command{
	Use:   "gen [source] [target] [theme]",
	Short: "Generate a site from Markdown files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf(
				"must provide source, target and theme directories",
			)
		}
		src, target, theme := args[0], args[1], args[2]

		blog, err := area.ParseArea(src, chromastyle)
		if err != nil {
			return fmt.Errorf("cannot parse: %w", err)
		}
		if err := blog.GenerateSite(
			target, theme, areainfo.PurposeStaticServe,
		); err != nil {
			return fmt.Errorf("cannot generate: %w", err)
		}
		return nil
	},
}

var chromastyle string

func init() {
	rootCmd.AddCommand(genCmd)
	rootCmd.Flags().StringVarP(
		&chromastyle, "style", "s", "based", "Chroma style to use",
	)
}
