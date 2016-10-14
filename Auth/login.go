package Auth

import (
	"net/http"
)

const htmlIndex = `<html><body>
Log in with <a href="/login/facebook">facebook</a>
</body></html>
`

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}
