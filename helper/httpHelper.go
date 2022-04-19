package helper

import (
	"encoding/json"
	"net/http"
)

// HTTPResponse is eFishery  standart response object
type HTTPResponse struct {
	Status  string
	Success bool
	Data    interface{}
}

// HTTPWriteResponse is http response json wrapper
// Default status code success
func HTTPWriteResponse(rw http.ResponseWriter, data interface{}, statusCode ...int) error {
	_statusCode := http.StatusOK
	_isSuccess := true

	var _data json.RawMessage
	var _res interface{}
	var err error

	if len(statusCode) > 0 {
		_statusCode = statusCode[0]
	}

	if _statusCode >= 400 {
		_isSuccess = false
	}

	switch v := data.(type) {
	case HTTPResponse, *HTTPResponse:
		_res = v

	case error:
		if len(statusCode) < 1 {
			_statusCode = http.StatusBadRequest
			_isSuccess = false
		}

		_res = HTTPResponse{
			Status:  http.StatusText(_statusCode),
			Success: false,
			Data:    v.Error(),
		}
	case []byte:
		_data = v
	default:

		_res = HTTPResponse{
			Status:  http.StatusText(_statusCode),
			Success: _isSuccess,
			Data:    v,
		}
	}

	if _data == nil {
		_data, err = json.Marshal(_res)
		if err != nil {
			rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=utf-8")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("{}"))
			return err
		}
	}

	rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=utf-8")
	rw.WriteHeader(_statusCode)
	_, err = rw.Write(_data)

	return err
}

func errResponse() {

}
