package main

import (
	"go/token"
	"strings"
)

func checkGoDocs(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
	tch := Filter(lc, DeclCommentFilter)
	for {
		ll := []*Lexeme{}
		for {
			l, ok := <-tch
			if !ok {
				return
			}
			if l.tok == token.ILLEGAL {
				break
			}
			ll = append(ll, l)
		}

		comm := beginGoDoc(ll)

		// does the comment line up with the next line?
		after := afterGoDoc(ll)
		if after.pos.Column != comm.pos.Column {
			continue
		}

		// does the comment have a token for documentation?
		fields := strings.Fields(comm.lit)
		if len(fields) < 2 {
			continue
		}

		// what token should the documentation match?
		cmplex := ll[len(ll)-1]
		if len(ll) >= 2 && ll[len(ll)-2].tok == token.IDENT {
			cmplex = ll[len(ll)-2]
		}
		if fields[1] == cmplex.lit {
			continue
		}

		// bad godoc
		cw := []CheckedWord{{fields[1], cmplex.lit}}
		cl := &CheckedLexeme{comm, "godoc", cw}
		outc <- cl
	}
}

// beginGoDoc gets the last comment block from a string of comments
func beginGoDoc(ll []*Lexeme) (comm *Lexeme) {
	wantLine := 0
	for _, l := range ll {
		if l.tok != token.COMMENT {
			break
		}
		if l.pos.Line != wantLine {
			comm = l
		}
		wantLine = l.pos.Line + 1
	}
	return comm
}

// afterGoDoc gets the first token following the comments
func afterGoDoc(ll []*Lexeme) *Lexeme {
	for _, l := range ll {
		if l.tok != token.COMMENT {
			return l
		}
	}
	return nil
}
