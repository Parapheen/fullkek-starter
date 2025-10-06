package templates

import "embed"

// Files exposes the embedded template filesystem.
//
// The filesystem layout is organized as follows:
//
//	base/     -> files shared by every project
//	features/ -> modular overlays that can be composed like LEGO bricks
//
// Additional features can extend the generator by adding new directories below
// internal/templates/features/.
//
//go:embed base/** features/**
var Files embed.FS
