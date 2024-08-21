package cli

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Yaroher2442/FamilySyncHub/build"
)

func Main() { //nolint:funlen // it's ok
	root := &cobra.Command{
		Version: build.Version,
		Run: func(_ *cobra.Command, _ []string) {
			if err := run(); err != nil {
				log.Fatalf("Fail start app: %v", err)
			}
		},
	}

	if err := root.Execute(); err != nil {
		log.Fatalf("root.Execute: %v", err)
	}
}
