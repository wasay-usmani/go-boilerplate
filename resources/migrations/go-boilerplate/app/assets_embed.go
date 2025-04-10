package app

import "embed"

//go:embed migrations
var MigrationAssets embed.FS
