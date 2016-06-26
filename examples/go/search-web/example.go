package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/ldap.v2"
)

const (
	host  = "ldap.example.com"
	port  = 389
	topDn = "dc=example,dc=com"
)

var (
	fieldsRegExp = regexp.MustCompile(`[^\s,]+`)
	templates    = template.New("").Funcs(template.FuncMap(map[string]interface{}{
		"attribute": func(a *ldap.EntryAttribute) template.HTML {
			h := make([]string, len(a.Values))
			switch a.Name {
			case "labeledURI":
				for i, v := range a.Values {
					pair := strings.SplitN(v, " ", 2)
					h[i] = fmt.Sprintf(
						"<a href='%s'>%s</a>",
						template.HTMLEscapeString(pair[0]),
						template.HTMLEscapeString(pair[1]))
				}
			case "mail":
				for i, v := range a.Values {
					h[i] = fmt.Sprintf(
						"<a href='mailto:%s'>%s</a>",
						template.HTMLEscapeString(v),
						template.HTMLEscapeString(v))
				}
			case "jpegPhoto":
				for i, v := range a.ByteValues {
					h[i] = fmt.Sprintf(
						"<img src='data:image/jpg;base64,%s'>",
						template.HTMLEscapeString(
							base64.StdEncoding.EncodeToString(
								v)))
				}
			case "roleOccupant":
				for i, v := range a.Values {
					h[i] = fmt.Sprintf(
						`<a href='javascript:search("(objectClass=*)", "*", "%s")'>%s</a>`,
						template.HTMLEscapeString(v),
						template.HTMLEscapeString(v))
				}
			default:
				for i, v := range a.Values {
					h[i] = template.HTMLEscapeString(v)
				}
			}
			return template.HTML(strings.Join(h, "<br>"))
		},
	}))
	searchResultsTemplate = template.Must(templates.Parse(`<!DOCTYPE html>
<html>
<head>
<style>
body {
	font-family: DejaVu Sans,Verdana,Geneva,sans-serif;
	font-size: 13px;
	background: #ececec;
}
form {
	display: flex;
	flex-direction: column;
}
.split {
	display: flex;
}
.split > div {
	display: flex;
	flex-direction: column;
	flex: 2;
}
.split > div+div {
	flex: 1;
	margin-left: 10px;
}
label, input {
	display: block;
	margin: 5px 0;
}
input[type=text] {
	font-size: 1.5em;
	padding: 5px;
}
.buttons {
	display: flex;
	align-items: center;
}
.buttons > input {
	width: 5em;
    padding: 5px 0;
}
.buttons > a {
	margin-left: 5px;
	text-decoration: underline;
	cursor: pointer;
}
#result {
	display: flex;
	flex-wrap: wrap;
	margin-left: -5px;
}
.card {
	box-shadow: 0 2px 2px 0 rgba(0,0,0,.14),0 3px 1px -2px rgba(0,0,0,.2),0 1px 5px 0 rgba(0,0,0,.12);
	background: #fff;
    border-radius: 2px;
	margin: 5px;
}
table {
	padding: 5px;
	width: 350px;
}
caption, .no-results {
	border-radius: 2px 2px 0 0;
	color: #fff;
	background: #204a87;
	padding: 10px 5px;
}
td+td {
	word-break: break-all;
}
.no-results {
	padding: 10px;
	border-radius: 2px;
	text-align: center;
	font-style: italic;
}
img {
	vertical-align: top;
}
</style>
</head>
<body>
<form method="GET">
	<label for="filter">Filter (<a href="https://tools.ietf.org/html/rfc4515" title="String Representation of Search Filters">RFC 4515</a>):</label>
	<input type="text" id="filter" name="filter" value="{{.Filter}}"/>
	<div class="split">
		<div>
			<label for="fields">Fields (<a href="https://tools.ietf.org/html/rfc4519" title="Schema for User Applications">RFC 4519</a>, <a href="https://tools.ietf.org/html/rfc2798" title="inetOrgPerson Object Class">RFC 2798</a>):</label>
			<input type="text" id="fields" name="fields" value="{{.Fields}}"/>
		</div>
		<div>
			<label for="base-dn">Base DN:</label>
			<input type="text" id="base-dn" name="base-dn" value="{{.BaseDn}}"/>
		</div>
	</div>
	<div class="buttons">
		<input type="submit" value="Search">
		<a id="search-all">All</a>
		<a id="search-people">People</a>
		<a id="search-roles">Roles</a>
		<a id="search-groups">Groups</a>
	</div>
</form>
<div id="result">
{{- range .Result.Entries}}
<div class="card">
<table>
<caption>{{.DN}}</caption>
<tbody>
	{{- range .Attributes}}
	<tr>
		<td>{{.Name}}</td>
		<td>{{attribute .}}</td>
	</tr>
	{{- end}}
</tbody>
</table>
</div>
{{- else}}
<div class="card no-results">no results found</div>
{{- end}}
</div>
<script>
var baseDnEl = document.getElementById("base-dn");
var filterEl = document.getElementById("filter");
var fieldsEl = document.getElementById("fields");
var searchFormEl = document.querySelector("form");
function search(filter, fields, baseDn) {
	filterEl.value = filter;
	fieldsEl.value = fields || "*";
	baseDnEl.value = baseDn || "{{.TopDn}}";
	searchFormEl.submit();
	return false;
}
function hook(id, filter, fields, baseDn) {
	var el = document.getElementById(id);
	el.title = filter;
	el.onclick = function(e) {
		e.preventDefault();
		search(filter, fields, baseDn);
	};
}
hook("search-all", "(objectClass=*)");
hook("search-people", "(objectClass=person)");
hook("search-roles", "(objectClass=organizationalRole)");
hook("search-groups", "(|(objectClass=groupOfNames)(objectClass=groupOfUniqueNames)(objectClass=group))");
</script>
</body>
</html>
`))
)

func search(filter string, fields string, baseDn string) (html []byte, err error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return
	}
	defer l.Close()

	search := ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		fieldsRegExp.FindAllString(fields, -1),
		nil)

	result, err := l.Search(search)
	if err != nil {
		return
	}

	data := struct {
		Result *ldap.SearchResult
		Filter string
		Fields string
		BaseDn string
		TopDn  string
	}{
		result,
		filter,
		fields,
		baseDn,
		topDn,
	}

	var b bytes.Buffer

	err = searchResultsTemplate.ExecuteTemplate(&b, "", data)
	if err != nil {
		return
	}

	html = b.Bytes()

	return
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		qs := r.URL.Query()
		filter := qs.Get("filter")
		if filter == "" {
			filter = "(objectClass=*)"
		}
		fields := qs.Get("fields")
		if fields == "" {
			fields = "*"
		}
		baseDn := qs.Get("base-dn")
		if baseDn == "" {
			baseDn = topDn
		}

		h, err := search(filter, fields, baseDn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(h)
	})

	fmt.Printf("Listening at http://%s:12345\n", host)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatalf("Failed to ListenAndServe: %v", err)
	}
}
