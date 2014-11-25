package pintu

import (
	"bytes"
	"html/template"
	"net/http"
)

type LoginPartial struct {
	Action   string
	Redirect string
	Name     string
}

func (p *LoginPartial) GetForm(r *http.Request) string {
	t := LoginFormTemplate()
	return p.RenderPartial(r, t)
}

func (p *LoginPartial) GetLink(r *http.Request) string {
	t := LoginFormTemplate()
	return p.RenderPartial(r, t)
}

func (p *LoginPartial) RenderPartial(r *http.Request, t *template.Template) string {
	buffer := new(bytes.Buffer)
	t.ExecuteTemplate(buffer, "partial.html", p)
	return buffer.String()
}

func GetTemplates() *template.Template {
	t := template.Must(template.New("main").Parse(`{{define "login.html"}}
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Login</title>
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
    <style>
      body {padding-top: 40px; padding-bottom: 40px; background-color: #fff}
      .form-control:focus {z-index: 2}
      .form-control {position: relative; height: auto; -webkit-box-sizing: border-box; -moz-box-sizing: border-box; box-sizing: border-box; padding: 10px; font-size: 16px;}
      button[type=submit] {margin: 10px 0}
    </style>
  </head>
  <body onload="document.getElementsByName('rd')[0].value=window.location.href">
    <div class="container-fluid">
      <h2 class="text-center">Please sign in</h2>
      {{.}}
    </div>
  </body>
</html>
{{end}}`))

	t = template.Must(t.Parse(`{{define "error.html"}}
<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Oops!</title>
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
    <style>
      .error-template {padding: 40px 15px;text-align: center;}
      .error-actions {margin-top:15px;margin-bottom:15px;}
      .error-actions .btn { margin-right:10px; }
    </style>
  </head>
  <body>
    <div class="container-fluid">
      <div class="row-fluid">
        <div class="col-md-12">
          <div class="error-template">
            <h1>Oops!</h1>
            <h2>{{.Title}}</h2>
            <div class="error-details">
	      {{.Message}}
            </div>
	    <div class="error-actions">
	      <a href="{{.LoginPath}}" class="btn btn-primary btn-lg">Sign In</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </body>
</html>
{{end}}`))

	return t
}

func LoginLinkTemplate() *template.Template {
	t := template.Must(template.New("LoginLink").Parse(`{{define "partial.html"}}
<div class="row-fluid">
  <form class="col-md-offset-4 col-md-4" method="GET" action="{{.Action}}" role="form">
    <input type="hidden" name="rd" value="{{.Redirect}}">
    <button class="btn btn-lg btn-primary btn-block" type="submit">Sign In with {{.Name}}</button>
  </form>
</div>
{{end}}`))
	return t
}

func LoginFormTemplate() *template.Template {
	t := template.Must(template.New("LoginForm").Parse(`{{define "partial.html"}}
<div class="row-fluid">
  <form class="col-md-offset-4 col-md-4" method="POST" action="{{.Action}}" role="form">
    <input type="hidden" name="rd" value="{{.Redirect}}">
    <input type="login" name="username" class="form-control" placeholder="Username" required autofocus>
    <input type="password" name="password" class="form-control" placeholder="Password" required>
    <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
  </form>
</div>
{{end}}`))
	return t
}
