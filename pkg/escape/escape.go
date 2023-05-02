package escape

// String 进行转义处理。
func String(in string) string {

	out := ""
	for _, r := range in {
		switch r {
		case '\n':
			out += "\\n"
		case '\t':
			out += "\\t"
		case '\\':
			out += "\\\\"
		case '"':
			out += "\\\""
		default:
			out += string(r)
		}
	}
	return out
}
