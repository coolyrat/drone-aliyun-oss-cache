package main

import (
	"drone-alicloud-oss/log"
	"drone-alicloud-oss/storage/oss"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var (
	gitVersion = "unknown"
	goVersion  = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "oss cache plugin"
	app.Version = gitVersion
	app.Action = run
	app.Flags = []cli.Flag{
		// Cache information

		cli.StringFlag{
			Name:   "filename",
			Usage:  "Filename for the cache",
			EnvVar: "PLUGIN_FILENAME",
		},
		cli.StringFlag{
			Name:   "root",
			Usage:  "root",
			EnvVar: "PLUGIN_ROOT",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "path",
			EnvVar: "PLUGIN_PATH",
		},
		cli.StringFlag{
			Name:   "fallback_path",
			Usage:  "fallback_path",
			EnvVar: "PLUGIN_FALLBACK_PATH",
		},
		cli.StringSliceFlag{
			Name:   "mount",
			Usage:  "cache directories",
			EnvVar: "PLUGIN_MOUNT",
		},
		cli.BoolFlag{
			Name:   "rebuild",
			Usage:  "rebuild the cache directories",
			EnvVar: "PLUGIN_REBUILD",
		},
		cli.BoolFlag{
			Name:   "restore",
			Usage:  "restore the cache directories",
			EnvVar: "PLUGIN_RESTORE",
		},
		cli.BoolFlag{
			Name:   "flush",
			Usage:  "flush the cache",
			EnvVar: "PLUGIN_FLUSH",
		},
		cli.StringFlag{
			Name:   "flush_age",
			Usage:  "flush cache files older than # days",
			EnvVar: "PLUGIN_FLUSH_AGE",
			Value:  "30",
		},
		cli.StringFlag{
			Name:   "flush_path",
			Usage:  "path to search for flushable cache files",
			EnvVar: "PLUGIN_FLUSH_PATH",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug plugin output",
			EnvVar: "PLUGIN_DEBUG",
		},

		// Aliyun OSS information

		cli.StringFlag{
			Name:   "ak",
			EnvVar: "PLUGIN_ACCESS_KEY",
			Usage:  "Aliyun OSS access key",
		},
		cli.StringFlag{
			Name:   "sk",
			EnvVar: "PLUGIN_ACCESS_KEY_SECRET",
			Usage:  "Aliyun OSS access key secret",
		},
		cli.StringFlag{
			Name:   "endpoint",
			EnvVar: "PLUGIN_ENDPOINT",
			Usage:  "Aliyun OSS endpoint",
		},
		cli.StringFlag{
			Name:   "bucket",
			EnvVar: "PLUGIN_BUCKET",
			Usage:  "Aliyun OSS bucket name",
		},

		// Build information (for setting defaults)

		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Value:  "master",
			Usage:  "repository default branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Logger.Fatal("fatal run drone-oss-cache",
			zap.String("gitVersion", gitVersion),
			zap.String("goVersion", goVersion),
			zap.Error(err))
	}
}

func run(c *cli.Context) error {

	// Determine the mode for the plugin
	rebuild := c.Bool("rebuild")
	restore := c.Bool("restore")

	if isMultipleModes(rebuild, restore) {
		return errors.New("must use a single mode: rebuild or restore")
	} else if !rebuild && !restore {
		return errors.New("no action specified")
	}

	var mode string
	var mount []string

	if rebuild {
		// Look for the mount points to rebuild
		mount = c.StringSlice("mount")

		if len(mount) == 0 {
			return errors.New("no mounts specified")
		}

		mode = RebuildMode
	} else {
		mode = RestoreMode
	}

	// Get the path to place the cache files
	path := c.GlobalString("path")

	// Defaults to <owner>/<repo>/<branch>/
	if len(path) == 0 {
		log.Logger.Info("no path specified. Creating default")

		path = fmt.Sprintf(
			"%s/%s/%s",
			c.String("repo.owner"),
			c.String("repo.name"),
			c.String("commit.branch"),
		)
	}

	// Get the filename
	filename := c.GlobalString("filename")

	if len(filename) == 0 {
		log.Logger.Info("no filename specified. Creating default")

		filename = "archive.tar"
	}

	s := oss.NewStorage(&oss.Config{
		Endpoint: c.String("endpoint"),
		AK:       c.String("ak"),
		SK:       c.String("sk"),
		Bucket:   c.String("bucket"),
	})

	p := &Plugin{
		Filename: filename,
		Path:     path,
		Mode:     mode,
		Mount:    mount,
		Storage:  s,
	}

	return p.Exec()
}

func isMultipleModes(bools ...bool) bool {
	var b bool
	for _, v := range bools {
		if b && b == v {
			return true
		}

		if v {
			b = true
		}
	}

	return false
}
