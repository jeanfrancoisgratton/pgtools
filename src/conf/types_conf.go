// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/15 08:08
// Original filename: src/conf/types_conf.go

package conf

import "regexp"

// Row holds one pg_settings entry we display.
type Row struct {
	Name      string
	Setting   string
	Unit      string
	Source    string
	Category  string
	ShortDesc string
}

// very conservative name validator: letters, digits, _ and .
var ValidName = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
var FullOutput = false
