package web

import "embed"

//go:embed assets/*
var StaticAssets embed.FS

//go:embed index.tmpl.html
var IndexTmpl string
