package cmux

import (
	"golang.org/x/net/http2/hpack"
	"io"
	"io/ioutil"
	"strings"

	http2_ "github.com/searKing/golang/go/net/cmux/internal/http2"
)

// HTTP2 parses the frame header of the first frame to detect whether the
// connection is an HTTP2 connection.
func HTTP2() MatcherFunc {
	return http2_.HasClientPreface
}

// HTTP2HeaderField returns a matcher matching the header fields of the first
// headers frame.
// writes the settings to the server if sendSetting is true.
// Prefer HTTP2HeaderField over this one, if the client does not block on receiving a SETTING frame.
func HTTP2HeaderField(sendSetting bool,
	match func(actual, expect map[string]hpack.HeaderField) bool,
	expects ...hpack.HeaderField) MatcherFunc {
	return func(w io.Writer, r io.Reader) bool {
		if !sendSetting {
			w = ioutil.Discard
		}
		return http2_.MatchHTTP2Header(w, r, nil, func(parsedHeader map[string]hpack.HeaderField) bool {
			var expectMap = map[string]hpack.HeaderField{}
			for _, expect := range expects {
				expectMap[expect.Name] = expect
			}
			return match(parsedHeader, expectMap)
		})
	}
}

// helper functions
func HTTP2HeaderFieldValue(sendSetting bool, match func(actual, expect string) bool, expects ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderField(sendSetting, func(actual, expect map[string]hpack.HeaderField) bool {
		for name, _ := range expect {
			if match(actual[name].Name, expect[name].Name) {
				return false
			}
		}
		return true
	}, expects...)
}

// HTTP2HeaderFieldEqual returns a matcher matching the header fields.
func HTTP2HeaderFieldEqual(sendSetting bool, headers ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderFieldValue(sendSetting, func(actual string, expect string) bool {
		return actual == expect
	}, headers...)
}

// HTTP2HeaderFieldPrefix returns a matcher matching the header fields.
// If the header with key name has a
// value prefixed with valuePrefix, this will match.
func HTTP2HeaderFieldPrefix(sendSetting bool, headers ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderFieldValue(sendSetting, strings.HasPrefix, headers...)
}
