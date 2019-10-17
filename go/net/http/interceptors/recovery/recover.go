package recovery

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/searKing/golang/go/error/builtin"
)

func Recover(logger *log.Logger, req *http.Request, recoverHandler func(err interface{})) {
	builtin.Recover(logger, recoverHandler, func() string {
		httpRequest, _ := httputil.DumpRequest(req, false)
		headers := strings.Split(string(httpRequest), "\r\n")
		for idx, header := range headers {
			current := strings.Split(header, ":")
			if current[0] == "Authorization" {
				headers[idx] = current[0] + ": *"
			}
		}
		return string(httpRequest)
	})
}