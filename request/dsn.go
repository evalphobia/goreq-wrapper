package request

import (
	"time"

	"github.com/franela/goreq"
)

var isDebug = false

var timeout = 10 * time.Second

func SetDebug() {
	isDebug = true
}

func SetTimeout(d time.Duration) {
	timeout = d
}

// DSN has request data
type DSN struct {
	goreq.Request

	Method string
	Uri    string
	Params interface{}
}

func (d *DSN) BasicAuth(user, pass string) *DSN {
	d.Request.BasicAuthUsername = user
	d.Request.BasicAuthPassword = pass
	return d
}

func (d *DSN) Param(p interface{}) *DSN {
	d.Params = p
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
	case "PUT":
		return d.PUT()
	case "DELETE":
		return d.DELETE()
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
	prepareParams(&d)
	return d.call()
}

func (d DSN) PUT() (*Body, error) {
	d.Request.Method = "PUT"
	prepareParams(&d)
	return d.call()
}

func (d DSN) DELETE() (*Body, error) {
	d.Request.Method = "DELETE"
	prepareParams(&d)
	return d.call()
}

func prepareParams(d *DSN) {
	if d.Params != nil {
		d.Request.Body = d.Params
		if _, ok := d.Params.(string); ok {
			d.Request.ContentType = "application/x-www-form-urlencoded"
		} else {
			d.Request.ContentType = "application/json"
		}
	}
}

func (d DSN) call() (*Body, error) {
	if isDebug {
		d.Request.ShowDebug = true
	}
	d.Request.Uri = d.Uri
	d.Request.Timeout = timeout
	res, err := d.Request.Do()
	if res == nil {
		return nil, err
	}
	return &Body{res.Body}, err
}
