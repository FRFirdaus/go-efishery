package parser_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"bitbucket.org/efishery/go-efishery/libs/parser"
	"gotest.tools/assert"
)

type Dependency struct {
	JSONByte []byte
	JSONP    parser.JSONParser
}

type photo struct {
	AlbumID      string `json:"album_id"`
	Title        string `json:"title"`
	ID           int64  `json:"id"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

var d Dependency

func TestMain(m *testing.M) {

	parse := parser.Init(
		parser.Options{
			parser.JSONOptions{
				Config: parser.JSONConfigDefault,
			},
		},
	)

	d.JSONP = parse.JSONParser()

	// get json data from API 5000 photos
	resp, err := http.Get("https://jsonplaceholder.typicode.com/photos")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	d.JSONByte = body

	// Run tests!
	exitVal := m.Run()

	os.Exit(exitVal)
}

func BenchmarkJSONParser(b *testing.B) {
	var result []photo
	for i := 0; i < b.N; i++ {
		err := d.JSONP.Unmarshal(d.JSONByte, &result)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkJSONStd(b *testing.B) {
	var result []photo
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(d.JSONByte, &result)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func TestUnmarshalJSONParser(t *testing.T) {
	testCases := []struct {
		name string
		in   []byte
		out  struct {
			expected map[string]interface{}
			err      error
		}
	}{
		{
			name: "Test Unmarshal String",
			in: []byte(`{
		       "name": "test"
		     }`,
			),
			out: struct {
				expected map[string]interface{}
				err      error
			}{
				expected: map[string]interface{}{
					"name": "test",
				},
				err: nil,
			},
		},
		{
			name: "Test Unmarshal Array",
			in: []byte(`{
		       "name": "test array",
					 "address": [
						 "alamat1",
						 "alamat2"
					 ]
	       }`,
			),
			out: struct {
				expected map[string]interface{}
				err      error
			}{
				expected: map[string]interface{}{
					"name":    "test array",
					"address": []interface{}{string("alamat1"), string("alamat2")},
				},
				err: nil,
			},
		},
		{
			name: "Test Unmarshal Array Object",
			in: []byte(`{
		       "name": "test array object",
					 "promo": [
						 {
							 "name": "promo 1",
							 "value": 800
						 },
						 {
							"name": "promo 2",
							"value": 2000
						}
					 ]
	       }`,
			),
			out: struct {
				expected map[string]interface{}
				err      error
			}{
				expected: map[string]interface{}{
					"name": "test array object",
					"promo": []interface{}{
						map[string]interface{}{"name": string("promo 1"), "value": float64(800)},
						map[string]interface{}{"name": string("promo 2"), "value": float64(2000)},
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		var result map[string]interface{}
		t.Run(tc.name, func(t *testing.T) {
			err := d.JSONP.Unmarshal(tc.in, &result)

			assert.NilError(t, err)
			assert.DeepEqual(t, result, tc.out.expected)
		})
	}
}
