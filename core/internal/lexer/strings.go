package lexer

// scanString scans a single- or double-quoted string literal. The
// Value includes the surrounding quotes. Unterminated strings are
// reported and the scanner stops at the end of the line rather than
// consuming the rest of the file, which keeps error recovery useful.
func (l *Lexer) scanString(line, col, start int, nl bool) Token {
	quote := l.cur.advance()
	for !l.cur.eof() {
		c := l.cur.peek()
		if c == '\\' {
			// Consume the backslash and whatever follows it. This is
			// enough to keep the scanner from seeing \" as the end
			// of the string; the printer is responsible for
			// re-escaping when it emits the token.
			l.cur.advance()
			if !l.cur.eof() {
				l.cur.advance()
			}
			continue
			
		}
		if c == quote {
			l.cur.advance()
			return l.emit(String, l.cur.src[start:l.cur.pos], line, col, start, nl)
		}
		if c == '\n' || c == '\r' {
			break
		}
		l.cur.advance()
	}
	l.errors = append(l.errors,
		errorf(line, col, "unterminated string literal"))
	return l.emit(Illegal, l.cur.src[start:l.cur.pos], line, col, start, nl)
}