package parser

import (
	"strings"
)

type curlParser struct {
	source  string
	url     string
	headers map[string]string
	args    map[string]string
}

func NewCurlParser(source string) *curlParser {
	return &curlParser{
		source:  source,
		url:     "",
		headers: map[string]string{},
		args:    map[string]string{},
	}
}

func (p *curlParser) Parse() *parserResult {
	p.preParse()

	return &parserResult{
		URL:       p.URL(),
		Method:    p.Method(),
		Cookie:    p.Cookie(),
		Referer:   p.Referer(),
		UserAgent: p.UserAgent(),
		Headers:   p.Headers(),
		Args:      p.Args(),
	}
}

func (p *curlParser) preParse() {
	slipt := strings.Split(p.source, "\n")
	for i := 0; i < len(slipt); i++ {
		line := strings.TrimSpace(strings.Trim(slipt[i], "\\"))
		parts := strings.SplitN(line, " ", 2)

		part1 := parts[0]
		part2 := ""
		if len(parts) >= 2 && parts[1] != "" {
			part2 = parts[1][1 : len(parts[1])-1]
		}

		if part1 == "curl" {
			p.url = part2
		} else if part1 == "-H" {
			part2Split := strings.SplitN(part2, " ", 2)
			part3 := ""
			if len(part2Split) >= 2 && part2Split[1] != "" {
				part3 = part2Split[1]
			}
			p.headers[strings.ToLower(part2Split[0][0:len(part2Split[0])-1])] = part3
		} else if strings.Index(part1, "-") == 0 {
			p.args[strings.ToLower(part1)] = part2
		}
	}
}

func (p *curlParser) URL() string {
	return p.url
}

func (p *curlParser) Method() string {
	if v, ok := p.args["-x"]; ok {
		//delete(p.args, "-x")
		return strings.ToUpper(v)
	}
	return "GET"
}

func (p *curlParser) Cookie() string {
	if v, ok := p.headers["cookie"]; ok {
		//delete(p.headers, "cookie")
		return v
	}
	return ""
}

func (p *curlParser) Referer() string {
	if v, ok := p.headers["referer"]; ok {
		//delete(p.headers, "referer")
		return v
	}
	return ""
}

func (p *curlParser) UserAgent() string {
	if v, ok := p.headers["user-agent"]; ok {
		//delete(p.headers, "user-agent")
		return v
	}
	return ""
}

func (p *curlParser) Headers() map[string]string {
	return p.headers
}

func (p *curlParser) Args() map[string]string {
	return p.args
}
