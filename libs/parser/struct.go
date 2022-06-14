package parser

type parserResult struct {
	URL       string
	Method    string
	Cookie    string
	Referer   string
	UserAgent string
	Headers   map[string]string
	Args      map[string]string
}
