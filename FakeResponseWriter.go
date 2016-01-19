package rizo

import "net/http"

//FakeResponseWriter ...
type FakeResponseWriter struct {
	Data         []byte
	ResponseCode int
	header       http.Header
}

//Header ...
func (instance *FakeResponseWriter) Header() http.Header {
	if instance.header == nil {
		instance.header = http.Header{}
	}
	return instance.header
}

//Write ...
func (instance *FakeResponseWriter) Write(data []byte) (int, error) {
	if instance.Data == nil {
		instance.Data = make([]byte, 0)
	}
	instance.Data = append(instance.Data, data...)
	return len(data), nil
}

//WriteHeader ...
func (instance *FakeResponseWriter) WriteHeader(header int) {
	instance.ResponseCode = header
}

//Reset ...
func (instance *FakeResponseWriter) Reset() {
	instance.Data = make([]byte, 0)
	instance.header = http.Header{}
}
