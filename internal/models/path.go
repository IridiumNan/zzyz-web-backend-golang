package models

import (
	"path"
)

// NOTE:
// Just support Linux deployment
// Use / as the spelitor of dirs

// For more details about the structure about datadir
// see docs/architecture_diagram

const (
	dataPath = "zzyz-data"

	idDir    = "By-ID"
	titleDir = "By-Title"

	buildDir       = "build"
	htmlPendingDir = "html_pending"
	htmlReleaseDIR = "html_release"

	rawDir = "raw"

	zipCacheDir = "zip_cache"
)

type PathSet struct {
	BuildDIr string

	IDDir string

	TitleDir string

	HTMLPendingDir string

	HTMLReleaseDir string

	RawDir string

	ZipCacheDir string
}

func DefaultPathSet() *PathSet {
	return &PathSet{
		BuildDIr:       path.Join(dataPath, buildDir),
		IDDir:          path.Join(dataPath, buildDir, idDir),
		TitleDir:       path.Join(dataPath, buildDir, titleDir),
		HTMLPendingDir: path.Join(dataPath, buildDir, htmlPendingDir),
		HTMLReleaseDir: path.Join(dataPath, buildDir, htmlReleaseDIR),
		RawDir:         path.Join(dataPath, rawDir),
		ZipCacheDir:    path.Join(dataPath, zipCacheDir),
	}
}
