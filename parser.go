package lexer

import (
	"fmt"
	"strings"
)

type kp map[string]string
type section map[string]kp

func parse(d string) (section, error) {
	if d == "" {
		return nil, fmt.Errorf("parse error: nothing to parse")
	}

	sections := make(section, 0)
	var sec string
	var key string

	l := lex(d)

	for t := range l.items {
		if (t.typ == itemSection || t.typ == itemKey || t.typ == itemValue) && strings.ContainsAny(t.val, "\n") {
			return nil, fmt.Errorf("parse error: %s can't contain new lines", t.val)
		}

		switch t.typ {
		case itemSection:
			sec = t.val
			_, ok := sections[sec]
			if !ok {
				sections[sec] = make(kp, 0)
			}
		case itemKey:
			key = strings.TrimSpace(t.val)
		case itemValue:
			s := sections[sec]
			s[key] = strings.TrimSpace(t.val)
		case itemError:
			return nil, fmt.Errorf("parse error: %s", t.val)
		}
	}

	return sections, nil
}
