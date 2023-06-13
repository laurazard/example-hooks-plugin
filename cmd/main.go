package main

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/laurazard/hints-plugin/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := RootCommand(dockerCli)
		originalPreRun := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				//nolint: wrapcheck
				return err
			}
			if originalPreRun != nil {
				return originalPreRun(cmd, args)
			}
			return nil
		}
		return cmd
	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Docker Inc.",
			Version:       "0.1",
			HookCommands:  []string{"build"},
		})
}

func RootCommand(dockerCli command.Cli) *cobra.Command {
	rootCmd := &cobra.Command{
		Short:            "Docker Hints",
		Use:              "hints",
		TraverseChildren: true,
	}

	rootCmd.AddCommand(hookCommand(dockerCli))
	return rootCmd
}

func hookCommand(dockerCli command.Cli) *cobra.Command {
	hookCmd := &cobra.Command{
		Use:   manager.HookSubcommandName,
		Short: "runs the plugins hooks",
		RunE: utils.Adapt(func(ctx context.Context, args []string) error {
			runHooks(dockerCli.Out(), []byte(args[0]))
			return nil
		}),
		Args: cobra.ExactArgs(1),
	}

	return hookCmd
}

func runHooks(out io.Writer, input []byte) {
	// TODO: change example to provide different
	// hint/template based on PluginData
	var c manager.HookPluginData
	_ = json.Unmarshal(input, &c)

	hint := "Run this image with `docker run " + manager.TemplateReplaceFlagValue("tag") + "`"
	returnType := manager.HookMessage{
		Template: hint,
	}
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "     ")
	_ = enc.Encode(returnType)
}
