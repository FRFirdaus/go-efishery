package parser

// Parser defines parser which owns 1 method, JSONParser
// Extended Parser will be defined in this package.
// TODO : JSONParser with schema validation, XMLParser, AvroParser
type Parser interface {
	// JSONParser return JSONParser Object
	JSONParser() JSONParser
}

type parser struct {
	json JSONParser
	opt  Options
}

// Options hold configurations of respective parser type.
type Options struct {
	JSON JSONOptions
}

func Init(opt Options) Parser {
	return &parser{
		json: initJSONP(opt.JSON),
	}
}

// JSONParser return JSONParser
// see :
//  - https://github.com/json-iterator/go
func (p *parser) JSONParser() JSONParser {
	return p.json
}
