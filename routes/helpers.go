package routes

import (
	fk "github.com/andrei-nita/goframework/framework"
	"github.com/gorilla/csrf"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"runtime"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

var (
	t        *template.Template
	upgrader = websocket.Upgrader{}
)

func parseTemplates() error {
	var err error
	files := fk.GetFilesAndDirRecursively("static/pages")
	t, err = template.New("").Funcs(template.FuncMap{"IsDev": fk.IsDev}).ParseFiles(files...)
	if err != nil {
		return err
	}
	return err
}

func allowMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		w.Header().Set("Allow", r.Method)
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func tmplExecute(w http.ResponseWriter, r *http.Request, tmplName string, data H) {
	var err error
	if !fk.Server.CacheTempls {
		err = parseTemplates()
	}
	err = t.Lookup(tmplName).Execute(w, data)
	if err != nil {
		http.Error(w, "please retry later", http.StatusInternalServerError)
		pc, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
	}
}

func createTemplateDataWithCSRF(r *http.Request) H {
	return H{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
}

func userIsAuth(r *http.Request, data H) bool {
	var auth bool

	isAuth, ok := r.Context().Value("cookie").(bool)
	if ok && isAuth {
		if data != nil {
			data["isAuth"] = true
		}
		auth = true
	}

	return auth
}

func methodCSRFAuth(w http.ResponseWriter, r *http.Request, method string) (data H, isAuth, next bool) {
	if ok := allowMethod(w, r, method); !ok {
		return
	}
	data = createTemplateDataWithCSRF(r)
	isAuth = userIsAuth(r, data)
	next = true
	return
}

func csrfAuth(r *http.Request) (data H, isAuth bool) {
	data = createTemplateDataWithCSRF(r)
	isAuth = userIsAuth(r, data)
	return
}

func methodAuth(w http.ResponseWriter, r *http.Request, method string) (data H, isAuth, next bool) {
	if ok := allowMethod(w, r, method); !ok {
		return
	}
	data = H{}
	isAuth = userIsAuth(r, data)
	next = true
	return
}

func methodCsrfFlashAuth(w http.ResponseWriter, r *http.Request, method, cookieName string) (data H, isAuth, next bool) {
	if ok := allowMethod(w, r, method); !ok {
		return
	}
	data = createTemplateDataWithCSRF(r)
	if msg := fk.CookieGetFlash(cookieName, w, r); msg != "" {
		data[cookieName] = msg
	}
	isAuth = userIsAuth(r, data)
	next = true

	return
}
