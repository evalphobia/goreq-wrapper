package request

import (
	"github.com/franela/goreq"
)

var isDebug = false

func SetDebug() {
	isDebug = true
}

// DSN has request data
type DSN struct {
	goreq.Request

	Method string
	Uri    string
	Params interface{}
}

func (d DSN) BasicAuth(user, pass string) DSN {
	d.Request.BasicAuthUsername = user
	d.Request.BasicAuthPassword = pass
	return d
}

func (d DSN) Debug() DSN {
	d.Request.ShowDebug = true
	return d
}

// Call sends request to endpoint with parameter(p)
func (d DSN) Call() (*Body, error) {
	switch d.Method {
	case "POST":
		return d.POST()
	default:
		return d.GET()
	}
}

func (d DSN) GET() (*Body, error) {
	d.Request.Method = "GET"
	if d.Params != nil {
		d.Request.QueryString = d.Params
	}
	return d.call()
}

func (d DSN) POST() (*Body, error) {
	d.Request.Method = "POST"
	if d.Params != nil {
		d.Request.Body = d.Params
		if _, ok := d.Params.(string); ok {
			d.Request.ContentType = "application/x-www-form-urlencoded"
		}
	}
	return d.call()
}

func (d DSN) call() (*Body, error) {
	if isDebug {
		d.Request.ShowDebug = true
	}
	d.Request.Uri = d.Uri
	res, err := d.Request.Do()
	return &Body{res.Body}, err
}
