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
	Type     string
	Btn      string
}

func (p *LoginPartial) GetForm(r *http.Request) string {
	t := LoginFormTemplate()
	return p.RenderPartial(r, t)
}

func (p *LoginPartial) GetLink(r *http.Request) string {
	t := LoginLinkTemplate()
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
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.1/css/bootstrap.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.1/css/bootstrap.css.map">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/material.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/ripples.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/material-wfont.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.2.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-social/4.2.1/bootstrap-social.min.css">
  </head>
  <body>
    <div class="container-fluid">
      <h2 class="text-center">Please login to proceed</h2>
      {{.}}
    </div>
    <script src="//code.jquery.com/jquery-1.10.2.min.js"></script>
    <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/js/ripples.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/js/material.min.js"></script>
    <script>$(function () { $.material.init(); $("input[name=rd]").val(window.location.href); });</script>
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
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.1/css/bootstrap.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.1/css/bootstrap.css.map">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/material.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/ripples.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/css/material-wfont.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.2.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-social/4.2.1/bootstrap-social.min.css">
  </head>
  <body>
    <hr />
    <div class="container-fluid">
      <div class="row-fluid">
        <div class="col-md-offset-3 col-md-6">
          <div class="error-template well">
            <h2>Oops! {{.Message}}</h2>
            <div class="error-actions">
	      <a href="{{.LoginPath}}" class="btn btn-primary btn-lg">
                <i class="fa fa-arrow-right"></i> back to login page
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
    <hr />
    <script src="//code.jquery.com/jquery-1.10.2.min.js"></script>
    <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/js/ripples.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/0.1.6/js/material.min.js"></script>
    <script>$(function () { $.material.init(); });</script>
  </body>
</html>
{{end}}`))

	return t
}

func LoginLinkTemplate() *template.Template {
	t := template.Must(template.New("LoginLink").Parse(`{{define "partial.html"}}
<div class="row-fluid">
  <div class="col-md-offset-4 col-md-4">
    <form method="GET" action="{{.Action}}" role="form">
      <input type="hidden" name="rd" value="{{.Redirect}}">
      <button class="btn btn-link btn-lg btn-block btn-social {{.Type}}" type="submit">
        <i class="fa {{.Btn}}"></i>
        Sign In with {{.Name}}
      </button>
    </form>
  </div>
</div>
{{end}}`))
	return t
}

func LoginFormTemplate() *template.Template {
	t := template.Must(template.New("LoginForm").Parse(`{{define "partial.html"}}
<div class="row-fluid">
  <div class="col-md-offset-4 col-md-4">
    <div class="well">
      <form method="POST" action="{{.Action}}" role="form">
        <fieldset>
          <input type="hidden" name="rd" value="{{.Redirect}}">
          <div class="form-group">
            <input type="login" name="username" class="form-control" placeholder="Username" required autofocus>
            <input type="password" name="password" class="form-control" placeholder="Password" required>
          </div>
          <div class="form-group">
            <button class="btn btn-lg btn-primary btn-block" type="submit">
              <i class="fa fa-lg fa-sign-in"></i>
              Sign in
            </button>
          </div>
        </fieldset>
      </form>
    </div>
  </div>
</div>
{{end}}`))
	return t
}
