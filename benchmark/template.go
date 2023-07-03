package benchmark

import _ "embed"

//go:embed templates/dat.tmpl
var DatTmpl string

//go:embed templates/job.tmpl
var JobTmpl string
