package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/hyloblog/hyloblog-ssg/internal/ast/area"
)

var serveCmd = &cobra.Command{
	Use:   "serve [source] [theme]",
	Short: "Serve a blog from Markdown files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf(
				"must provide source and theme directories",
			)
		}
		src, theme := args[0], args[1]

		h, err := choosehandler(src, theme, livereload)
		if err != nil {
			return fmt.Errorf("cannot choose handler: %w", err)
		}
		s := &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        h,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Printf("listening on %s...", s.Addr)
		return s.ListenAndServe()
	},
}

func choosehandler(src, theme string, livereload bool) (http.Handler, error) {
	if livereload {
		return area.CreateLiveHandler(src, theme, chromastyle), nil
	}
	blog, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	h, err := blog.Handler(theme)
	if err != nil {
		return nil, fmt.Errorf("cannot make http handler: %w", err)
	}
	return h, nil
}

var (
	port       int
	livereload bool
)

func init() {
	serveCmd.Flags().IntVarP(
		&port, "port", "p", 8000, "Port to serve on",
	)
	serveCmd.Flags().BoolVarP(
		&livereload, "livereload", "D", false, "Enable live reloading",
	)
	serveCmd.Flags().StringVarP(
		&chromastyle, "style", "s", "based", "Chroma style to use",
	)
	rootCmd.AddCommand(serveCmd)
}
