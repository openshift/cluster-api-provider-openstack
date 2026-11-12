// Copyright 2011-2019 Canonical Ltd
// Copyright 2025 The go-yaml Project Contributors
// SPDX-License-Identifier: Apache-2.0

// Tag resolution for YAML scalars.
// Determines implicit types (int, float, bool, null) from untagged scalar
// values.

package libyaml

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// resolveMapItem holds a resolved value and its YAML tag for exact string
// matches in the resolution table.
type resolveMapItem struct {
	value any
	tag   string
}

// Resolver handles tag resolution for YAML nodes.
type Resolver struct {
	opts *Options
}

// NewResolver creates a new Resolver with the given options.
func NewResolver(opts *Options) *Resolver {
	return &Resolver{opts: opts}
}

// Resolve walks the node tree and resolves tags for untagged nodes.
// This is called after composition to:
// - Default quoted scalars to !!str
// - Default sequences to !!seq
// - Default mappings to !!map
// - Resolve plain scalars to implicit types (int, float, bool, null)
func (r *Resolver) Resolve(n *Node) {
	if n == nil {
		return
	}

	switch n.Kind {
	case ScalarNode:
		if n.Tag == "" {
			if n.Style&(SingleQuotedStyle|DoubleQuotedStyle|LiteralStyle|FoldedStyle) != 0 {
				// Quoted scalars default to !!str without value resolution
				n.Tag = strTag
			} else {
				// Plain scalars: resolve type from value
				n.Tag, _ = resolve("", n.Value, r.formatError())
			}
		}

	case SequenceNode:
		if n.Tag == "" {
			n.Tag = seqTag
		}
		for _, child := range n.Content {
			r.Resolve(child)
		}

	case MappingNode:
		if n.Tag == "" {
			n.Tag = mapTag
		}
		for _, child := range n.Content {
			r.Resolve(child)
		}

	case DocumentNode:
		for _, child := range n.Content {
			r.Resolve(child)
		}

	case AliasNode:
		// Alias nodes point to already-resolved nodes
	}
}

// resolve determines the YAML tag and Go value for a scalar string.
// It takes a tag hint and the scalar string value, and returns the resolved
// tag and the corresponding Go value (int, float, bool, etc.).
// If the tag is already specified and non-resolvable, it returns the input
// unchanged.
func resolve(tag string, in string, formatter func(*LoadError) string) (rtag string, out any) {
	tag = shortTag(tag)
	if !resolvableTag(tag) {
		return tag, in
	}

	defer func() {
		switch tag {
		case "", rtag, strTag, binaryTag:
			return
		case floatTag:
			// Numbers are resolved into json.Number, so an explicit
			// !!float tag only changes the tag of an integer.
			if rtag == intTag {
				if _, ok := out.(json.Number); ok {
					rtag = floatTag
					return
				}
			}
		}
		Fail(formatResolverError(
			fmt.Sprintf("cannot construct %s `%s` as a %s", shortTag(rtag), in, shortTag(tag)),
			Mark{},
			formatter,
		))
	}()

	// Any data is accepted as a !!str or !!binary.
	// Otherwise, the prefix is enough of a hint about what it might be.
	hint := byte('N')
	if in != "" {
		hint = resolveTable[in[0]]
	}
	if hint != 0 && tag != strTag && tag != binaryTag {
		// Handle things we can lookup in a map.
		if item, ok := resolveMap[in]; ok {
			return item.tag, item.value
		}

		// Base 60 floats are a bad idea, were dropped in YAML 1.2, and
		// are purposefully unsupported here. They're still quoted on
		// the way out for compatibility with other parser, though.

		switch hint {
		case 'M':
			// We've already checked the map above.

		case '.', 'D', 'S':
			// Int or float.
			if tag == timestampTag {
				// Ignore timestamp tag.
				tag = strTag
				return strTag, in
			}

			plain := strings.ReplaceAll(in, "_", "")
			if isYAMLInt(in) {
				return intTag, json.Number(normalizeNumber(plain))
			}
			if isYAMLBases(in) {
				var neg bool
				switch plain[0] {
				case '-':
					neg = true
					fallthrough
				case '+':
					plain = plain[1:]
				}
				u, err := strconv.ParseUint(plain, 0, 64)
				if err == nil {
					if !neg {
						return intTag, json.Number(strconv.FormatUint(u, 10))
					}
					if u <= -math.MinInt64 {
						return intTag, json.Number(strconv.FormatInt(-int64(u), 10))
					}
				}
			}
			if isYAMLFloat(in) {
				return floatTag, json.Number(normalizeNumber(plain))
			}
		default:
			panic("internal error: missing handler for resolver table: " + string(rune(hint)) + " (with " + in + ")")
		}
	}
	return strTag, in
}

// resolveTable provides a fast lookup table for initial character-based
// classification during tag resolution.
// resolveMap maps specific scalar strings to their resolved values and tags.
var (
	resolveTable = make([]byte, 256)
	resolveMap   = make(map[string]resolveMapItem)
)

// isYAMLInt reports whether s is a decimal integer in YAML style, which is
// an optional sign followed by a zero, or by digits not starting with a zero,
// with underscore separators: ^[-+]?(?:0|[1-9][0-9_]*)$
func isYAMLInt(s string) bool {
	if s = trimSign(s); s == "" {
		return false
	}
	if s[0] == '0' {
		return len(s) == 1
	}
	if s[0] < '1' || s[0] > '9' {
		return false
	}
	return trimDigits(s) == ""
}

// isYAMLBases reports whether s is a binary, an octal, or a hexadecimal
// integer in YAML style, which is an optional sign followed by a zero, a base
// prefix, and the digits of the base, with underscore separators:
// ^[-+]?0(?:b[01_]+|o?[0-7_]+|x[0-9a-fA-F_]+)$
func isYAMLBases(s string) bool {
	if s = trimSign(s); len(s) < 2 || s[0] != '0' {
		return false
	}
	switch s = s[1:]; s[0] {
	case 'b':
		return isBaseChars(s[1:], isBinaryChar)
	case 'o':
		return isBaseChars(s[1:], isOctalChar)
	case 'x':
		return isBaseChars(s[1:], isHexChar)
	default:
		// The octal prefix is optional for compatibility with YAML 1.1.
		return isBaseChars(s, isOctalChar)
	}
}

// isYAMLFloat reports whether s is a float in YAML style, which is an optional
// sign followed by digits with an optional decimal point and an optional
// exponent, with underscore separators:
// ^[-+]?(?:\.[0-9][0-9_]*|[0-9][0-9_]*(?:\.(?:[0-9][0-9_]*)?)?)(?:[eE][-+]?[0-9][0-9_]*)?$
func isYAMLFloat(s string) bool {
	if s = trimSign(s); s == "" {
		return false
	}
	if s[0] == '.' {
		// The fraction requires a digit when there is no integer part.
		if s = s[1:]; s == "" || !isDecimalChar(s[0]) {
			return false
		}
		s = trimDigits(s)
	} else {
		if !isDecimalChar(s[0]) {
			return false
		}
		if s = trimDigits(s); s == "" {
			return true
		}
		if s[0] == '.' {
			if s = s[1:]; s != "" && isDecimalChar(s[0]) {
				s = trimDigits(s)
			}
		}
	}
	if s == "" {
		return true
	}
	if s[0] != 'e' && s[0] != 'E' {
		return false
	}
	s = trimSign(s[1:])
	if s == "" || !isDecimalChar(s[0]) {
		return false
	}
	return trimDigits(s) == ""
}

// trimSign trims the leading sign of s.
func trimSign(s string) string {
	if s != "" && (s[0] == '-' || s[0] == '+') {
		return s[1:]
	}
	return s
}

// trimDigits trims the leading digits of s, which may contain underscore
// separators after the first digit.
func trimDigits(s string) string {
	i := 1
	for ; i < len(s); i++ {
		if !isDecimalChar(s[i]) && s[i] != '_' {
			break
		}
	}
	return s[i:]
}

// isBaseChars reports whether s consists of the digits of a base, which
// requires at least one character, and may contain underscore separators.
func isBaseChars(s string, isBaseChar func(byte) bool) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isBaseChar(s[i]) && s[i] != '_' {
			return false
		}
	}
	return true
}

func isDecimalChar(c byte) bool { return '0' <= c && c <= '9' }

func isBinaryChar(c byte) bool { return c == '0' || c == '1' }

func isOctalChar(c byte) bool { return '0' <= c && c <= '7' }

func isHexChar(c byte) bool {
	return isDecimalChar(c) || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// init initializes the resolveTable with character class mappings for tag resolution.
func init() {
	t := resolveTable
	t[int('+')] = 'S' // Sign
	t[int('-')] = 'S'
	for _, c := range "0123456789" {
		t[int(c)] = 'D' // Digit
	}
	for _, c := range "yYnNtTfFoO~<" { // < for merge key <<
		t[int(c)] = 'M' // In map
	}
	t[int('.')] = '.' // Float (potentially in map)

	resolveMapList := []struct {
		v   any
		tag string
		l   []string
	}{
		{true, boolTag, []string{"true", "True", "TRUE"}},
		{false, boolTag, []string{"false", "False", "FALSE"}},
		{nil, nullTag, []string{"", "~", "null", "Null", "NULL"}},
		{math.NaN(), floatTag, []string{".nan", ".NaN", ".NAN"}},
		{math.Inf(+1), floatTag, []string{".inf", ".Inf", ".INF"}},
		{math.Inf(+1), floatTag, []string{"+.inf", "+.Inf", "+.INF"}},
		{math.Inf(-1), floatTag, []string{"-.inf", "-.Inf", "-.INF"}},
		{"<<", mergeTag, []string{"<<"}},
	}

	m := resolveMap
	for _, item := range resolveMapList {
		for _, s := range item.l {
			m[s] = resolveMapItem{item.v, item.tag}
		}
	}
}

// resolvableTag checks if a tag can be automatically resolved from a scalar
// value.
func resolvableTag(tag string) bool {
	switch tag {
	case "", strTag, boolTag, intTag, floatTag, nullTag, timestampTag:
		return true
	}
	return false
}

// normalizeNumber rewrites a YAML number literal in JSON number syntax, so
// that the resulting [json.Number] can be marshaled as JSON.
func normalizeNumber(s string) string {
	if isJSONNumber(s) {
		return s
	}
	var sign string
	if s[0] == '-' || s[0] == '+' {
		if s[0] == '-' {
			sign = "-"
		}
		s = s[1:]
	}
	var exp string
	if i := strings.IndexAny(s, "eE"); i >= 0 {
		exp, s = s[i:], s[:i]
	}
	var frac string
	if i := strings.IndexByte(s, '.'); i >= 0 {
		if frac, s = s[i:], s[:i]; frac == "." {
			frac = ""
		}
	}
	if s = strings.TrimLeft(s, "0"); s == "" {
		s = "0"
	}
	return sign + s + frac + exp
}

// isJSONNumber reports whether s is already in JSON number syntax.
func isJSONNumber(s string) bool {
	i := 0
	if s[0] == '-' {
		i++
	}
	// A leading plus sign, or no digit before the decimal point.
	if i == len(s) || !isDecimalChar(s[i]) {
		return false
	}
	// A redundant leading zero.
	if s[i] == '0' && i+1 < len(s) && isDecimalChar(s[i+1]) {
		return false
	}
	// No digit after the decimal point.
	if i := strings.IndexByte(s, '.'); i >= 0 &&
		(i+1 == len(s) || !isDecimalChar(s[i+1])) {
		return false
	}
	return true
}

// encodeBase64 encodes s as base64 that is broken up into multiple lines
// as appropriate for the resulting length.
func encodeBase64(s string) string {
	const lineLen = 70
	encLen := base64.StdEncoding.EncodedLen(len(s))
	lines := encLen/lineLen + 1
	buf := make([]byte, encLen*2+lines)
	in := buf[0:encLen]
	out := buf[encLen:]
	base64.StdEncoding.Encode(in, []byte(s))
	k := 0
	for i := 0; i < len(in); i += lineLen {
		j := i + lineLen
		if j > len(in) {
			j = len(in)
		}
		k += copy(out[k:], in[i:j])
		if lines > 1 {
			out[k] = '\n'
			k++
		}
	}
	return string(out[:k])
}

// formatResolverError creates a LoadError for resolver-stage errors.
func (r *Resolver) formatError() func(*LoadError) string {
	if r == nil || r.opts == nil {
		return nil
	}
	return r.opts.FormatLoadError
}

func formatResolverError(message string, mark Mark, formatter func(*LoadError) string) *LoadError {
	return &LoadError{
		Stage:     ResolverStage,
		Mark:      mark,
		Message:   message,
		formatter: formatter,
	}
}

// formatResolverErrorContext creates a LoadError with both context and
// problem information for resolver-stage errors.
func formatResolverErrorContext(context string, contextMark Mark, message string, mark Mark) *LoadError {
	return &LoadError{
		Stage:       ResolverStage,
		ContextMark: contextMark,
		ContextMsg:  context,
		Mark:        mark,
		Message:     message,
	}
}
