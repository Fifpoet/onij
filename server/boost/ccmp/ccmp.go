package ccmp

import (
	"cmp"
	"hash/crc32"
	"onij/boost/collection/collext"
	"onij/boost/conv"
	"onij/boost/exp"
	"slices"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/net/html"
)

func ArraysEqual[T cmp.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}

func ArraysEqualFunc[T any, K cmp.Ordered](a, b []T, orderFn func(T) K) bool {
	if len(a) != len(b) {
		return false
	}

	return ArraysEqual(collext.Pick(a, orderFn), collext.Pick(b, orderFn))
}

func StringEqual(o, a string) bool {
	if len(o) != len(a) {
		return false
	}

	if len(o) <= 1024*2 {
		return o == a
	}

	crc32O := crc32.ChecksumIEEE(conv.StringToBytes(o))
	crc32A := crc32.ChecksumIEEE(conv.StringToBytes(a))
	return crc32O == crc32A
}

func PtrEqual[T comparable](o, a *T) bool {
	if o == nil && a == nil {
		return true
	}
	if o == nil || a == nil {
		return false
	}
	switch v := (any)(o).(type) {
	case *string:
		return StringEqual(*v, *(any)(a).(*string))
	default:
		return *o == *a
	}
}

func PtrValueOrZeroEqual[T comparable](o, a *T) bool {
	switch vo := (any)(o).(type) {
	case *string:
		return StringEqual(exp.ValueOrZero(vo), exp.ValueOrZero((any)(a).(*string)))
	default:
		return exp.ValueOrZero(o) == exp.ValueOrZero(a)
	}
}

func Min[T constraints.Integer | constraints.Float](o, a T) T {
	if o < a {
		return o
	}
	return a
}

func Max[T constraints.Integer | constraints.Float](o, a T) T {
	if o > a {
		return o
	}
	return a
}

func RichTextBasicEqual(o, a string) bool {
	if len(o) == 0 && len(o) == len(a) {
		return true
	}

	basicO := basicRichText(o)
	basicA := basicRichText(a)
	return StringEqual(basicO, basicA)
}

func basicRichText(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return s
	}

	const separator = "\u200C"

	var extract func(*html.Node)
	var builder strings.Builder

	extract = func(node *html.Node) {
		switch node.Type {
		case html.TextNode:
			if text := strings.TrimSpace(node.Data); len(text) != 0 {
				builder.WriteString(text)
				builder.WriteString(separator)
			}
		case html.ElementNode:
			if node.Data != "img" {
				break
			}
			var src, width, height string
			for _, attr := range node.Attr {
				switch attr.Key {
				case "src":
					src = attr.Val
				case "width":
					width = attr.Val
				case "height":
					height = attr.Val
				}
			}
			builder.WriteString(src)
			builder.WriteString(separator)
			builder.WriteString(width)
			builder.WriteString(separator)
			builder.WriteString(height)
			builder.WriteString(separator)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	extract(doc)
	return builder.String()
}
