package request

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	"github.com/franela/goreq"
)

// Body is wrapper struct for goreq.Body
type Body struct {
	*goreq.Body
}

// ToByte converts request body to []byte
func (b *Body) ToByte() ([]byte, error) {
	return ioutil.ReadAll(b.Body)
}

// ToMap converts request body to map
func (b *Body) ToMap() (map[string]interface{}, error) {
	var r map[string]interface{}
	byt, err := b.ToByte()
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(byt, &r)
	return r, err
}

// ParseXML assigns request body data(XML) to entity
func (b *Body) ParseXML(s interface{}) error {
	body, err := b.ToByte()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, s)
}
