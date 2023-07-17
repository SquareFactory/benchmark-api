package benchmark

import _ "embed"

//go:embed templates/dat.tmpl
var DatTmpl string

//go:embed templates/multinode.tmpl
var MultiNodeTmpl string

//go:embed templates/singlenode.tmpl
var SingleNodeTmpl string
