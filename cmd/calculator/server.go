package main

import (
	"calculator/pkg/calculator"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.New("form").Parse(`
<html>
<head><title>Go Calculator</title></head>
<body>
	<h1>Calculator</h1>
	<form method="POST">
		<input name="expression" value="{{.Input}}" />
		<input type="submit" value="Calculate" />
	</form>
	<p>Result: {{.Result}}</p>
	<p style="color:red">{{.Error}}</p>
</body>
</html>
`))

type PageData struct {
	Input  string
	Result string
	Error  string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}
		if r.Method == http.MethodPost {
			r.ParseForm()
			expr := r.FormValue("expression")
			data.Input = expr
			result, err := calculator.Calc(expr)
			if err != nil {
				data.Error = err.Error()
			} else {
				data.Result = fmt.Sprintf("%v", result)
			}
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("Serving on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
