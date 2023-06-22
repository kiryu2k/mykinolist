package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
	mock_service "github.com/kiryu-dev/mykinolist/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestController_IdentifyUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, tokens *model.Tokens)
	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		cookieName          string
		cookieValue         string
		tokens              model.Tokens
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer sOmEt0kee3N",
			cookieName:  "refreshToken",
			tokens: model.Tokens{
				AccessToken: "sOmEt0kee3N",
			},
			mockBehavior: func(s *mock_service.MockAuthService, tokens *model.Tokens) {
				s.EXPECT().ParseAccessToken(tokens.AccessToken).Return(int64(666), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "666",
		},
		{
			name:        "Access token out of date",
			headerName:  "Authorization",
			headerValue: "Bearer tokenOutOfDate",
			cookieName:  "refreshToken",
			cookieValue: "ImValidToken",
			tokens: model.Tokens{
				AccessToken:  "tokenOutOfDate",
				RefreshToken: "ImValidToken",
			},
			mockBehavior: func(s *mock_service.MockAuthService, tokens *model.Tokens) {
				var id int64 = 88
				s.EXPECT().ParseAccessToken(tokens.AccessToken).Return(
					id, fmt.Errorf("token expiration date has passed"),
				)
				s.EXPECT().ParseRefreshToken(tokens.RefreshToken).Return(id, nil)
				s.EXPECT().UpdateTokens(id).Return(&model.Tokens{}, nil) // user gets new tokens
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "88",
		},
		{
			name:        "Invalid tokens",
			headerName:  "Authorization",
			headerValue: "Bearer invalidToken",
			cookieName:  "refreshToken",
			tokens: model.Tokens{
				AccessToken: "invalidToken",
			},
			mockBehavior: func(s *mock_service.MockAuthService, tokens *model.Tokens) {
				s.EXPECT().ParseAccessToken(tokens.AccessToken).Return(
					int64(0),
					fmt.Errorf("token is malformed: could not base64 decode header: illegal base64 data at input byte 36"),
				)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: "{\"error\":\"token is malformed: could not base64 decode header: illegal base64 data at input byte 36\"}\n",
		},
		{
			name:        "Access token out of date but there is no refresh token",
			headerName:  "Authorization",
			headerValue: "Bearer tokenOutOfDate",
			cookieName:  "refreshToken",
			tokens: model.Tokens{
				AccessToken: "tokenOutOfDate",
			},
			mockBehavior: func(s *mock_service.MockAuthService, tokens *model.Tokens) {
				var id int64 = 88
				s.EXPECT().ParseAccessToken(tokens.AccessToken).Return(
					id, fmt.Errorf("token expiration date has passed"),
				)
				s.EXPECT().ParseRefreshToken(tokens.RefreshToken).Return(int64(0), fmt.Errorf("token is malformed: token contains an invalid number of segments"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "{\"error\":\"token expiration date has passed\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuthService(c)
			testCase.mockBehavior(auth, &testCase.tokens)
			services := &service.Service{AuthService: auth}
			handler := &authHandler{service: services}
			router := mux.NewRouter()
			router.Use(handler.identifyUser)
			router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
				id := r.Context().Value(userIDKey{}).(int64)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", id)))
			}).Methods(http.MethodGet)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/auth", nil)
			http.SetCookie(w, &http.Cookie{
				Name:     testCase.cookieName,
				Value:    testCase.cookieValue,
				Path:     "/auth",
				MaxAge:   cookieMaxAge,
				HttpOnly: true,
			})
			req.Header.Add(testCase.headerName, testCase.headerValue)
			cookieHeader := fmt.Sprintf("%s=%s", w.Result().Cookies()[0].Name, w.Result().Cookies()[0].Value)
			req.Header.Add("Cookie", cookieHeader)
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
