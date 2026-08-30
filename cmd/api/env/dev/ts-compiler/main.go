package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/resources"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/src"
	"github.com/evanw/esbuild/pkg/api"
)

var (
	watch = false
)

func main() {
	flag.BoolVar(&watch, "watch", false, "Enable watcher mode")
	flag.Parse()

	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("failure", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	sctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	return build(sctx)
}

func build(ctx context.Context) error {

	entryMap := map[string]string{
		src.TsDirectoryPath + "/modules/public/greetings/index.ts":      "modules/public/greetings/index",
		src.TsDirectoryPath + "/modules/public/greetings/hello.ts":      "modules/public/greetings/hello",
		src.TsDirectoryPath + "/modules/public/greetings/hello-name.ts": "modules/public/greetings/hello-name",
	}

	var entryPoints []api.EntryPoint
	for inputPath, outputPath := range entryMap {
		entryPoints = append(entryPoints, api.EntryPoint{
			InputPath:  inputPath,
			OutputPath: outputPath,
		})
	}

	opts := api.BuildOptions{
		Outdir:              resources.JsDirectoryPath,
		EntryPointsAdvanced: entryPoints,
		Bundle:              true,
		Format:              api.FormatESModule,
		LogLevel:            api.LogLevelInfo,
		MinifyIdentifiers:   true,
		MinifySyntax:        true,
		MinifyWhitespace:    true,
		Sourcemap:           api.SourceMapLinked,
		Target:              api.ESNext,
		Write:               true,
	}

	if watch {
		slog.Info("ts-compiler watching...")

		buildCtx, err := api.Context(opts)
		if err != nil {
			return err
		}
		defer buildCtx.Dispose()

		if err := buildCtx.Watch(api.WatchOptions{}); err != nil {
			return err
		}

		<-ctx.Done()
		return nil
	}

	slog.Info("building...")

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		errs := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = errors.New(err.Text)
		}
		return errors.Join(errs...)
	}

	return nil
}
