package main

import (
	"context"
	"go-project-template/internal/cmd"
	"os"
	"os/signal"
)

// @title           Project
// @version         0.0.1
// @description     Primary backend services the Project.

// @contact.name   API Support
// @contact.url    https://greycanary.io
// @contact.email  dev@greycanary.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
