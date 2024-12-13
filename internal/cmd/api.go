package cmd

import (
	"context"
	"fmt"
	"go-project-template/cmdutil"
	"go-project-template/internal/api"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 5000
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := cmdutil.NewLogger("nuop_api")
			defer func() { _ = logger.Sync() }()

			db, err := cmdutil.NewDatabasePool(ctx)
			if err != nil {
				fmt.Printf("DB Conn error %v", err)
				return err
			}
			defer db.Close()

			api := api.NewAPI(ctx, logger, db)
			srv := api.Server(port)

			go func() { _ = srv.ListenAndServe() }()

			logger.Info("started api", zap.Int("port", port))

			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}

	return cmd
}
