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
	testTable := []struct {
		name                string
		inputBody           string
		inputUser           model.SignUpUserDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
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
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "{\"id\":1}\n",
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
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "{\"error\":\"username must consist of letters or numbers, also it must contain from 6 to 50 characters\"}\n",
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
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "{\"error\":\"email must consist of letters and numbers, also it mustn't exceed 100 characters\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuthService(c)
			testCase.mockBehavior(auth, &testCase.inputUser)
			services := &service.Service{AuthService: auth}
			handler := &authHandler{service: services}
			router := mux.NewRouter()
			router.HandleFunc("/signup", handler.signUp).Methods(http.MethodPost)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(testCase.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
