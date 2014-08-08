package data

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Data holds data obtained from the request body and url query parameters.
// Because Data is built from multiple sources, sometimes there will be more
// than one value for a given key. You can use Get, Set, Add, and Del to access
// the first element for a given key or access the map directly to access additional
// elements for a given key. You can also use helper methods to convert the first
// value for a given key to a different type (e.g. bool or int).
type Data url.Values

// Parse parses the request body and url query parameters into
// Data. The content in the body of the request has a higher priority,
// will be added to Data first, and will be the result of any operation
// which gets the first element for a given key (e.g. Get, GetInt, or GetBool).
func Parse(req *http.Request) (Data, error) {
	values := url.Values{}
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(2048); err != nil {
			return nil, err
		}
		for key, val := range req.MultipartForm.Value {
			values.Add(key, value)
		}
	} else if strings.Contains(contentType, "form-urlencoded") {
		if err := req.ParseForm(); err != nil {
			return nil, err
		}
		for key, val := range req.PostForm {
			values.Add(key, value)
		}
	}
	for key, val := range req.URL.Query() {
		values.Add(key, value)
	}
	return Data(values), nil
}

// Add adds the value to key. It appends to any existing values associated with key.
func (d Data) Add(key string, value string) {
	url.Values(d).Add(key, value)
}

// Del deletes the values associated with key.
func (d Data) Del(key string) {
	url.Values(d).Del(key)
}

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux") sorted by key.
func (d Data) Encode() string {
	return url.Values(d).Encode()
}

// Get gets the first value associated with the given key. If there are no values
// associated with the key, Get returns the empty string. To access multiple values,
// use the map directly.
func (d Data) Get(key string) string {
	return url.Values(d).Get(key)
}

// Set sets the key to value. It replaces any existing values.
func (d Data) Set(key string, value string) {
	url.Values(d).Set(key, value)
}

// KeyExists returns true iff data[key] exists. If the value of data[key] is an empty
// string, it is still considered to be in existence.
func (d Data) KeyExists(key string) bool {
	_, found := d[key]
	return found
}

// GetInt returns the first element in data[key] converted to an int.
func (d Data) GetInt(key string) int {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return 0
	}
	str := d[key][0]
	if result, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return result
	}
}

// GetBool returns the first element in data[key] converted to a bool.
func (d Data) GetBool(key string) bool {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return false
	}
	str := d[key][0]
	if result, err := strconv.ParseBool(str); err != nil {
		panic(err)
	} else {
		return result
	}
}

// GetStringsSplit returns the first element in data[key] split into a slice delimited by delim.
func (d Data) GetStringsSplit(key string, delim string) []string {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return nil
	}
	return strings.Split(d[key][0], delim)
}

// Validator returns a Validator which can be used to easily validate data.
func (d Data) Validator() *Validator {
	return &Validator{
		data: d,
	}
}