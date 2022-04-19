package authmw

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	libtoken "bitbucket.org/efishery/go-efishery/libs/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareHTTP(t *testing.T) {
	secretKey := "secret-key"
	tokenSvc := libtoken.NewService(secretKey)

	tokenString, err := tokenSvc.Generate(jwt.StandardClaims{Id: "userId"})
	if err != nil {
		t.Error(err)
		return
	}

	testCases := []struct {
		Token      string
		StatusCode int
	}{
		{
			"",
			http.StatusUnauthorized,
		},
		{
			tokenString,
			http.StatusOK,
		},
	}

	for _, v := range testCases {
		req, err := http.NewRequest("GET", "/secret-endpoint?token="+v.Token, nil)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HTTPAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
		}), secretKey)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != v.StatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, v.StatusCode)
		}

	}

}

func TestAuthMiddlewareGin(t *testing.T) {

	secretKey := "secret-key"
	tokenSvc := libtoken.NewService(secretKey)

	tokenString, err := tokenSvc.Generate(jwt.StandardClaims{Id: "userId"})
	if err != nil {
		t.Error(err)
		return
	}

	handlers := gin.HandlersChain{
		GinAuthMiddleware(secretKey),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		},
	}
	testCases := []struct {
		Token      string
		StatusCode int
	}{
		{
			"",
			http.StatusUnauthorized,
		},
		{
			tokenString,
			http.StatusOK,
		},
	}

	for _, v := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.POST("/secret-endpoint", handlers...)

		c.Request, _ = http.NewRequest(http.MethodPost, "/secret-endpoint?token="+v.Token, bytes.NewBuffer([]byte("{}")))

		r.ServeHTTP(w, c.Request)

		if w.Code != v.StatusCode {
			t.Errorf("Expected status %d, got %d", v.StatusCode, w.Code)
		}
	}
}
