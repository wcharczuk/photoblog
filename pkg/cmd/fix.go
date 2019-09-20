package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/uuid"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/config"
)

// Fix returns the fix tree of commands.
func Fix(flags config.Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fix",
		Short: "Fix contains commands used to modify the posts database",
	}

	slugify := &cobra.Command{
		Use:   "slugify",
		Short: "Rename all subdirectories in the posts according to the slugify rules",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, configPath, err := config.ReadConfig(flags)
			sh.Fatal(err)

			log := Logger(flags, "slugify")
			slant.Print(log.Output, "BLOGCTL")

			if configPath != "" {
				log.Infof("using config path: %s", configPath)
			}

			postsPath := cfg.PostsPathOrDefault()
			contents, err := ioutil.ReadDir(postsPath)
			sh.Fatal(err)
			var from, temp, to string
			for _, object := range contents {
				if !object.IsDir() {
					continue
				}
				from = filepath.Join(postsPath, object.Name())
				to = filepath.Join(postsPath, stringutil.Slugify(object.Name()))

				if from != to {
					if flags.DryRun != nil && *flags.DryRun {
						log.Infof("(dry run) would rename from %s to %s", from, to)
					} else {
						temp = filepath.Join(postsPath, uuid.V4().String())
						err = os.Rename(from, temp)
						sh.Fatal(err)
						err = os.Rename(temp, to)
						sh.Fatal(err)
						log.Infof("renamed from %s to %s", from, to)
					}
				}
			}
		},
	}
	cmd.AddCommand(slugify)
	return cmd
}
