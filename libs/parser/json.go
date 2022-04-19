package parser

import (
	jsoniter "github.com/json-iterator/go"
)

type JSONParser interface {
	// Marshal go structs into bytes
	Marshal(orig interface{}) ([]byte, error)
	// Unmarshal bytes into go structs
	Unmarshal(blob []byte, dest interface{}) error
}

type jsonConfig string

/**
 * JSON Parser Configuration.
 * There are 4 options :
 *  JSONConfigDefault
 *  JSONConfigCompatibleWithStdLibrary
 *  JSONConfigFastest
 *  TODO JSONConfigCustom - Use this if you want to apply JSON parser config otherwise any config set will not be applied
 */
const (
	// JSONConfigDefault set :
	//	EscapeHTML :	true
	JSONConfigDefault jsonConfig = `default`

	// JSONConfigCompatibleWithStdLibrary
	//  EscapeHTML:                 true
	//  SortMapKeys:					      true
	//  ValidateJsonRawMessage:			true
	JSONConfigCompatibleWithStdLibrary jsonConfig = `standard`

	// JSONConfigFastest
	//  EscapeHTML:                    	false
	//  MarshalFloatWith6Digits:       	true
	//  ObjectFieldMustBeSimpleString: 	true
	JSONConfigFastest jsonConfig = `fastest`

	// JSONConfigCustom
	//	Custom Configuration which is set in JSONOptions
	JSONConfigCustom jsonConfig = `custom`
)

// JSONOptions contains configuration options used by JSON Parser during parsing.
// see http://jsoniter.com/migrate-from-go-std.html
type JSONOptions struct {
	// Config based on jsoniterator-go default config
	Config jsonConfig
}

type jsonparser struct {
	API jsoniter.API
	opt JSONOptions
}

// initJSONP initialize JSON parser with declared options
func initJSONP(opt JSONOptions) JSONParser {
	var jsonAPI jsoniter.API
	switch opt.Config {

	case JSONConfigDefault:
		jsonAPI = jsoniter.ConfigDefault

	case JSONConfigFastest:
		jsonAPI = jsoniter.ConfigFastest

	case JSONConfigCompatibleWithStdLibrary:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}

	p := &jsonparser{
		API: jsonAPI,
		opt: opt,
	}

	return p
}

func (p *jsonparser) Marshal(orig interface{}) (result []byte, err error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)

	stream.WriteVal(orig)

	result = make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, stream.Error
	}

	copy(result, stream.Buffer())

	return
}

func (p *jsonparser) Unmarshal(blob []byte, dest interface{}) (err error) {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)

	iter.ReadVal(dest)
	if iter.Error != nil {
		return iter.Error
	}

	return
}
