// SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"context"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

var vanityCfg = []byte(`
host: go.example.com
paths:
  /mycoolproject:
    repo: https://github.com/user/mycoolproject
    packages:
      - pkg/package1
      - pkg/package2

  /myothercoolproject:
    repo: https://github.com/user/myothercoolproject

`)

var indexTemplate = `
<!DOCTYPE html>
<html>
<body>
<h1>Welcome to {{.Host}}</h1>
<ul>
{{range .Vanity}}<li><a href="https://pkg.go.dev/{{.Path}}">{{.Path}}</a></li>{{end}}
</ul>
</body>
</html>
`

var indexOut = `
<!DOCTYPE html>
<html>
<body>
<h1>Welcome to go.example.com</h1>
<ul>
<li><a href="https://pkg.go.dev/go.example.com/mycoolproject">go.example.com/mycoolproject</a></li><li><a href="https://pkg.go.dev/go.example.com/myothercoolproject">go.example.com/myothercoolproject</a></li>
</ul>
</body>
</html>
`

var projectTemplate = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Import}} {{.VCS}} {{.Repo}}">
<meta name="go-source" content="{{.Import}} {{.Display}}">
</head>
<body>
Nothing to see here folks!
</body>
</html>
`

var mycoolprojectOut = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="go.example.com/mycoolproject git https://github.com/user/mycoolproject">
<meta name="go-source" content="go.example.com/mycoolproject https://github.com/user/mycoolproject https://github.com/user/mycoolproject/tree/master{/dir} https://github.com/user/mycoolproject/blob/master{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here folks!
</body>
</html>
`

var myOtherCoolProjectOut = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="go.example.com/myothercoolproject git https://github.com/user/myothercoolproject">
<meta name="go-source" content="go.example.com/myothercoolproject https://github.com/user/myothercoolproject https://github.com/user/myothercoolproject/tree/master{/dir} https://github.com/user/myothercoolproject/blob/master{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here folks!
</body>
</html>
`

func TestGenerate(t *testing.T) {
	t.Parallel()
	a := assertions.New(t)
	ctx := context.Background()
	gen, err := New(ctx, vanityCfg)
	a.So(err, should.BeNil)
	a.So(len(gen.paths), should.Equal, 2)
	index, err := gen.Index(ctx, indexTemplate)
	a.So(err, should.BeNil)
	a.So(string(index), should.Equal, indexOut)
	vanity, err := gen.Project(ctx, projectTemplate)
	a.So(err, should.BeNil)
	a.So(len(vanity.items), should.Equal, 2)
	mcp := vanity.items["/mycoolproject"]
	a.So(string(mcp.Content), should.Equal, mycoolprojectOut)
	a.So(len(mcp.PkgNames), should.Equal, 2)
	mocp := vanity.items["/myothercoolproject"]
	a.So(string(mocp.Content), should.Equal, myOtherCoolProjectOut)
	a.So(len(mocp.PkgNames), should.Equal, 0)
	invalid := vanity.items["/mynonexistantproject"]
	a.So(invalid.Content, should.BeNil)
	a.So(invalid.PkgNames, should.BeNil)
}
