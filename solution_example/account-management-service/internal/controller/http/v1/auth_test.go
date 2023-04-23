package v1

import (
	"account-management-service/internal/mocks/servicemocks"
	"account-management-service/internal/service"
	"account-management-service/pkg/validator"
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRoutes_SignUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.AuthCreateUserInput
	}

	type MockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name            string
		args            args
		inputBody       string
		mockBehaviour   MockBehaviour
		wantStatusCode  int
		wantRequestBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.AuthCreateUserInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(1, nil)
			},
			wantStatusCode:  201,
			wantRequestBody: `{"id":1}` + "\n",
		},
		{
			name:            "Invalid password: not provided",
			args:            args{},
			inputBody:       `{"username":"test"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password is required"}` + "\n",
		},
		{
			name:            "Invalid password: too short",
			args:            args{},
			inputBody:       `{"username":"test","password":"Qw!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must be between 8 and 32 characters"}` + "\n",
		},
		{
			name:            "Invalid password: too long",
			args:            args{},
			inputBody:       `{"username":"test","password":"Qwerty!123456789012345678901234567890"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must be between 8 and 32 characters"}` + "\n",
		},
		{
			name:            "Invalid password: no uppercase",
			args:            args{},
			inputBody:       `{"username":"test","password":"qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 uppercase letter(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no lowercase",
			args:            args{},
			inputBody:       `{"username":"test","password":"QWERTY!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 lowercase letter(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no digits",
			args:            args{},
			inputBody:       `{"username":"test","password":"Qwerty!!"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 digit(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no special characters",
			args:            args{},
			inputBody:       `{"username":"test","password":"Qwerty11"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 special character(s)"}` + "\n",
		},
		{
			name:            "Invalid username: not provided",
			args:            args{},
			inputBody:       `{"password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username is required"}` + "\n",
		},
		{
			name:            "Invalid username: too short",
			args:            args{},
			inputBody:       `{"username":"t","password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username must be at least 4 characters"}` + "\n",
		},
		{
			name:            "Invalid username: too long",
			args:            args{},
			inputBody:       `{"username":"testtesttesttesttesttesttesttesttest","password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username must be at most 32 characters"}` + "\n",
		},
		{
			name:            "Invalid request body",
			args:            args{},
			inputBody:       `{"username" test","password":"Qwerty!1"`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"invalid request body"}` + "\n",
		},
		{
			name: "Auth service error",
			args: args{
				ctx: context.Background(),
				input: service.AuthCreateUserInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(0, service.ErrUserAlreadyExists)
			},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"user already exists"}` + "\n",
		},
		{
			name: "Internal server error",
			args: args{
				ctx: context.Background(),
				input: service.AuthCreateUserInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(0, errors.New("some error"))
			},
			wantStatusCode:  500,
			wantRequestBody: `{"message":"internal server error"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			auth := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(auth, tc.args)
			services := &service.Services{Auth: auth}

			// create test server
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			g := e.Group("/auth")
			newAuthRoutes(g, services.Auth)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}

func TestAuthRoutes_SignIn(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.AuthGenerateTokenInput
	}

	type mockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name            string
		args            args
		inputBody       string
		mockBehaviour   mockBehaviour
		wantStatusCode  int
		wantRequestBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.AuthGenerateTokenInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return("token", nil)
			},
			wantStatusCode:  200,
			wantRequestBody: `{"token":"token"}` + "\n",
		},
		{
			name:            "Invalid username: not provided",
			args:            args{},
			inputBody:       `{"password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username is required"}` + "\n",
		},
		{
			name:            "Invalid password: not provided",
			args:            args{},
			inputBody:       `{"username":"test"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password is required"}` + "\n",
		},
		{
			name: "Wrong username or password",
			args: args{
				ctx: context.Background(),
				input: service.AuthGenerateTokenInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return("", service.ErrUserNotFound)
			},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"invalid username or password"}` + "\n",
		},
		{
			name: "Internal server error",
			args: args{
				ctx: context.Background(),
				input: service.AuthGenerateTokenInput{
					Username: "test",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"username":"test","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return("", errors.New("some error"))
			},
			wantStatusCode:  500,
			wantRequestBody: `{"message":"internal server error"}` + "\n",
		},
		{
			name:            "Invalid request body",
			args:            args{},
			inputBody:       `{qw"qwdf)00)))`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"invalid request body"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create service mock
			auth := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(auth, tc.args)
			services := &service.Services{Auth: auth}

			// create test server
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			g := e.Group("/auth")
			newAuthRoutes(g, services.Auth)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}
