package lexer

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var lexTests = []lexTest{
	{"empty", "", []item{item{typ: itemEOF}}},
	{"brokensection", "[foo\nfoo=bar\n", []item{item{typ: itemLeftBracket}, item{typ: itemError}}},
	{"brokenkey", "[foo]\nfoo\n", []item{
		item{typ: itemLeftBracket},
		item{typ: itemSection},
		item{typ: itemRightBracket},
		item{typ: itemError},
	}},
	{"brokenvalue", "[foo]\nfoo=bar", []item{
		item{typ: itemLeftBracket},
		item{typ: itemSection},
		item{typ: itemRightBracket},
		item{typ: itemKey},
		item{typ: itemEqualSign},
		item{typ: itemError},
	}},
	{"success", "[foo]\nfoo=foo\n", []item{
		item{typ: itemLeftBracket},
		item{typ: itemSection},
		item{typ: itemRightBracket},
		item{typ: itemKey},
		item{typ: itemEqualSign},
		item{typ: itemValue},
	}},
	{"successcomplex", "[foo]\nfoo=foo\n\n[bar]\nbar=bar\n", []item{
		item{typ: itemLeftBracket},
		item{typ: itemSection},
		item{typ: itemRightBracket},
		item{typ: itemKey},
		item{typ: itemEqualSign},
		item{typ: itemValue},
		item{typ: itemLeftBracket},
		item{typ: itemSection},
		item{typ: itemRightBracket},
		item{typ: itemKey},
		item{typ: itemEqualSign},
		item{typ: itemValue},
	}},
}

func TestLexer(t *testing.T) {
	for _, lt := range lexTests {
		l := lex(lt.input)
		for _, lti := range lt.items {
			li := <-l.items
			if lti.typ != li.typ {
				t.Errorf("test %s did not pass\n", lt.name)
				return
			}
		}
	}
}
