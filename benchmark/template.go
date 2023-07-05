package benchmark

import _ "embed"

//go:embed templates/dat.tmpl
var DatTmpl string

//go:embed templates/sbatch.tmpl
var SbatchTmpl string
