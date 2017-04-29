package lexer

import (
	"reflect"
	"testing"
)

type parseTest struct {
	name     string
	input    string
	err      bool
	errVal   string
	sections section
}

var parseTests = []parseTest{
	{"empty", "", true, "parse error: nothing to parse", nil},
	{"brokensection", "[foo\nfoo=bar\n", true, "parse error: no right bracket found", nil},
	{"brokenkey", "[foo]\nfoo\n", true, "parse error: no equal sign found", nil},
	{"brokenvalue", "[foo]\nfoo=bar", true, "parse error: no new line found", nil},
	{"success", "[foo]\nfoo=FOO\n", false, "", section{"foo": kp{"foo": "FOO"}}},
	{"successcomplex", "[foo]\nfoo=FOO\n\n[bar]\nbar=BAR\n", false, "", section{
		"foo": kp{"foo": "FOO"},
		"bar": kp{"bar": "BAR"},
	}},
}

func TestParser(t *testing.T) {
	for _, pt := range parseTests {
		p, err := parse(pt.input)

		if (pt.err && (p != nil || err == nil || err.Error() != pt.errVal)) ||
			(!pt.err && (p == nil || err != nil || !reflect.DeepEqual(pt.sections, p))) {
			t.Errorf("test %s did not pass\n", pt.name)
		}
	}
}
