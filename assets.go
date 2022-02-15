package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets4548c4de5a42a759758868d96443b68ee4bacf04 = "<!DOCTYPE html>\n<html lang=\"zh-cn\">\n\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n    <title>{{.title}}</title>\n</head>\n\n<body>\n    <h1>{{.contentForH1}}</h1>\n    <main>{{.contentForMain}}</main>\n</body>\n\n</html>"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"index.tmpl"}}, map[string]*assets.File{
	"/index.tmpl": &assets.File{
		Path:     "/index.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1644134751, 1644134751282000000),
		Data:     []byte(_Assets4548c4de5a42a759758868d96443b68ee4bacf04),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1644134806, 1644134806900730039),
		Data:     nil,
	}}, "")
