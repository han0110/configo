package util

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/tabwriter"
	"text/template"
)

const (
	// CharDot represents character ".".
	CharDot = "."
	// CharDash represents character "-".
	CharDash = "-"
	// CharUnderscore represents character "_".
	CharUnderscore = "_"
)

var (
	delimeterRegexp = regexp.MustCompile(`\.|-|_`)
	camelCaseRegexp = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// TabExecutor executes template with tab writter.
func TabExecutor(w io.Writer, tmpl *template.Template, data interface{}) error {
	tw := tabwriter.NewWriter(w, 4, 0, 4, ' ', 0)
	if err := tmpl.Execute(tw, data); err != nil {
		return err
	}
	return tw.Flush()
}

// ToDotCase transform str into dot case (e.g. log.level).
func ToDotCase(str string) string {
	return strings.ToLower(ReplaceDelimiter(str, CharDot))
}

// ToDashCase transform str into dash case (e.g. log-level).
func ToDashCase(str string) string {
	return strings.ToLower(ReplaceDelimiter(str, CharDash))
}

// ToScreamingCase transform str into screaming case (e.g. log.level to LOG_LEVEL).
func ToScreamingCase(str string) string {
	return strings.ToUpper(ReplaceDelimiter(str, CharUnderscore))
}

// ReplaceDelimiter changes delimiter of string to new and also connect camel case with it.
func ReplaceDelimiter(str, delimeter string) string {
	str = delimeterRegexp.ReplaceAllString(str, delimeter)
	str = camelCaseRegexp.ReplaceAllString(str, fmt.Sprintf("${1}%s${2}", delimeter))
	return str
}

// UniqueStrings returns unique slice.
func UniqueStrings(slice []string) []string {
	seen := make(map[string]struct{}, len(slice))
	count := 0
	for _, str := range slice {
		if _, yes := seen[str]; !yes {
			seen[str] = struct{}{}
			slice[count] = str
			count++
		}
	}
	return slice[:count]
}

// StringsToSet turns slice into set
func StringsToSet(slice []string) map[string]bool {
	set := make(map[string]bool)
	for _, str := range slice {
		set[str] = true
	}
	return set
}
