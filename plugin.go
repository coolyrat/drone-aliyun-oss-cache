package main

import (
	"drone-alicloud-oss/log"
	"github.com/drone/drone-cache-lib/archive/util"
	"github.com/drone/drone-cache-lib/cache"
	"github.com/drone/drone-cache-lib/storage"
	"go.uber.org/zap"
	pathutil "path"
)

// Plugin structure
type Plugin struct {
	Filename string
	Path     string
	Mode     string
	Mount    []string

	Storage storage.Storage
}

const (
	// RestoreMode for resotre mode string
	RestoreMode = "restore"
	// RebuildMode for rebuild mode string
	RebuildMode = "rebuild"
)

// Exec runs the plugin
func (p *Plugin) Exec() error {
	var err error

	at, err := util.FromFilename(p.Filename)

	if err != nil {
		return err
	}

	c := cache.New(p.Storage, at)

	path := pathutil.Join(p.Path, p.Filename)

	if p.Mode == RebuildMode {
		log.Logger.Info("Rebuilding cache", zap.String("path", path))
		err = c.Rebuild(p.Mount, path)

		if err == nil {
			log.Logger.Info("Cache rebuilt")
		}
	}

	if p.Mode == RestoreMode {
		log.Logger.Info("Restoring cache", zap.String("path", path))
		err = c.Restore(path, "")

		if err == nil {
			log.Logger.Info("Cache restored")
		}
	}

	return err
}
