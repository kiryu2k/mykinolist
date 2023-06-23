package controller

import (
	"bytes"
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

func TestController_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, userDTO *model.SignUpUserDTO)
	type testCase struct {
		name                 string
		inputBody            string
		inputUser            model.SignUpUserDTO
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}
	testCases := []testCase{
		{
			name:      "OK",
			inputBody: `{"username":"testUser2023","email":"test-user@gmail.com","password":"qweRty2023"}`,
			inputUser: model.SignUpUserDTO{
				Username: "testUser2023",
				Email:    "test-user@gmail.com",
				Password: "qweRty2023",
			},
			mockBehavior: func(s *mock_service.MockAuthService, userDTO *model.SignUpUserDTO) {
				s.EXPECT().SignUp(userDTO).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1}\n",
		},
		{
			name:      "Empty fields",
			inputBody: `{"username":"","email":"","password":""}`,
			inputUser: model.SignUpUserDTO{},
			mockBehavior: func(s *mock_service.MockAuthService, userDTO *model.SignUpUserDTO) {
				s.EXPECT().SignUp(userDTO).Return(
					int64(0),
					fmt.Errorf("username must consist of letters or numbers, also it must contain from 6 to 50 characters"),
				)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"username must consist of letters or numbers, also it must contain from 6 to 50 characters\"}\n",
		},
		{
			name:      "Without email",
			inputBody: `{"username":"testUser2023","password":"qweRty2023"}`,
			inputUser: model.SignUpUserDTO{
				Username: "testUser2023",
				Password: "qweRty2023",
			},
			mockBehavior: func(s *mock_service.MockAuthService, userDTO *model.SignUpUserDTO) {
				s.EXPECT().SignUp(userDTO).Return(
					int64(0),
					fmt.Errorf("email must consist of letters and numbers, also it mustn't exceed 100 characters"),
				)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"email must consist of letters and numbers, also it mustn't exceed 100 characters\"}\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuthService(c)
			tc.mockBehavior(auth, &tc.inputUser)
			var (
				services = &service.Service{AuthService: auth}
				handler  = &authHandler{service: services}
				router   = mux.NewRouter()
			)
			router.HandleFunc("/signup", handler.signUp).Methods(http.MethodPost)
			var (
				w   = httptest.NewRecorder()
				req = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tc.inputBody))
			)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestController_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthService, userDTO *model.SignInUserDTO)
	type testCase struct {
		name                   string
		inputBody              string
		inputUser              model.SignInUserDTO
		mockBehavior           mockBehavior
		expectedStatusCode     int
		expectedResponseBody   string
		expectedResponseHeader http.Header
	}
	testCases := []testCase{
		{
			name:      "OK",
			inputBody: `{"email":"testUserMail@yahoo.com","password":"PA55WorD"}`,
			inputUser: model.SignInUserDTO{
				Email:    "testUserMail@yahoo.com",
				Password: "PA55WorD",
			},
			mockBehavior: func(s *mock_service.MockAuthService, userDTO *model.SignInUserDTO) {
				s.EXPECT().SignIn(userDTO).Return(&model.Tokens{
					AccessToken:  "IM_ACCESS_TOKEN_TRUST_ME",
					RefreshToken: "WELL_I_GUESS_IM_REFRESH_TOKEN",
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"access_token\":\"IM_ACCESS_TOKEN_TRUST_ME\",\"refresh_token\":\"WELL_I_GUESS_IM_REFRESH_TOKEN\"}\n",
			expectedResponseHeader: http.Header{
				"Authorization": {"Bearer IM_ACCESS_TOKEN_TRUST_ME"},
				"Content-Type":  {"application/json"},
				"Set-Cookie":    {"refreshToken=WELL_I_GUESS_IM_REFRESH_TOKEN; Path=/auth; Max-Age=2592000; HttpOnly"},
			},
		},
		{
			name:      "Without password",
			inputBody: `{"email":"testUserMail@yahoo.com","password":""}`,
			inputUser: model.SignInUserDTO{
				Email: "testUserMail@yahoo.com",
			},
			mockBehavior: func(s *mock_service.MockAuthService, userDTO *model.SignInUserDTO) {
				s.EXPECT().SignIn(userDTO).Return(
					nil,
					fmt.Errorf("password must contain from 8 to 30 characters, be at least one uppercase letter, one lowercase letter and one number"))
			},
			expectedStatusCode:     http.StatusBadRequest,
			expectedResponseBody:   "{\"error\":\"password must contain from 8 to 30 characters, be at least one uppercase letter, one lowercase letter and one number\"}\n",
			expectedResponseHeader: http.Header{"Content-Type": {"application/json"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuthService(c)
			tc.mockBehavior(auth, &tc.inputUser)
			var (
				services = &service.Service{AuthService: auth}
				handler  = &authHandler{service: services}
				router   = mux.NewRouter()
			)
			router.HandleFunc("/signin", handler.signIn).Methods(http.MethodPost)
			var (
				w   = httptest.NewRecorder()
				req = httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(tc.inputBody))
			)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
			assert.Equal(t, tc.expectedResponseHeader, w.Header())
		})
	}
}
