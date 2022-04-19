package requestidmw

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequestMiddlewareHTTP(t *testing.T) {

	testCases := []struct {
		Header           string
		HasRequestHeader bool
	}{
		{
			uuid.NewString(),
			true,
		},
		{
			"",
			false,
		},
	}

	for i, v := range testCases {
		req, err := http.NewRequest("GET", "/some-endpoint", nil)
		if err != nil {
			t.Fatal(err)
		}

		// set header
		if v.HasRequestHeader {
			req.Header.Set(HeaderKey, v.Header)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Useage middleware
		handler := HttpMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// validate request
			if v.HasRequestHeader {
				if GetRequestId(r) != v.Header {
					t.Error("Request id not match", i)
				}

			} else {
				if GetRequestId(r) == "" {
					t.Error("Request is empty", i)
				}

			}
		}))

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// validate response
		respHeader := rr.Header().Get(HeaderKey)
		if v.HasRequestHeader {
			if respHeader != v.Header {
				t.Error("Response id not match", i)
			}

		} else {
			if respHeader == "" {
				t.Error("Response is empty", i)
			}

		}

	}

}

func TestRequestMiddlewareGin(t *testing.T) {

	testCases := []struct {
		Header           string
		HasRequestHeader bool
	}{
		{
			uuid.NewString(),
			true,
		},
		{
			"",
			false,
		},
	}

	for i, v := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		// Usage global middleware
		// r.Use(GinMiddleware())

		r.POST("/some-endpoint",
			// Use middleware specific endpoint
			GinMiddleware(),
			func(c *gin.Context) {

				// validate request
				if v.HasRequestHeader {
					if GetGinRequestId(c) != v.Header {
						t.Error("Request id not match", i)
					}

				} else {
					if GetGinRequestId(c) == "" {
						t.Error("Request is empty", i)
					}

				}
				c.JSON(http.StatusOK, "")
			})

		c.Request, _ = http.NewRequest(http.MethodPost, "/some-endpoint", bytes.NewBuffer([]byte("{}")))
		// set header
		if v.HasRequestHeader {
			c.Request.Header.Set(HeaderKey, v.Header)
		}
		r.ServeHTTP(w, c.Request)

		// validate response
		respHeader := w.Header().Get(HeaderKey)
		if v.HasRequestHeader {
			if respHeader != v.Header {
				t.Error("Response id not match", i)
			}

		} else {
			if respHeader == "" {
				t.Error("Response is empty", i)
			}

		}
	}
}
