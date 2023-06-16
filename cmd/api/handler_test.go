package main

import (
	"authentication/cmd/mock/dbmocks"
	"authentication/cmd/mock/loggermocks"
	"authentication/cmd/model"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func mockGetInfoByEmail(ui *dbmocks.UserInterface, email string, result *model.User, err error) {
	(*ui).On("GetInfoByEmail", email).Return(result, err)
}

func mockMatchPassword(ui *dbmocks.UserInterface, pw string, result bool) {
	(*ui).On("MatchPassword", pw).Return(result)
}

var _ = Describe("Test Chi Route", func() {
	When("Walk chi route", func() {
		It("should return found", func() {
			walkRoutes := func(chiRoutes *chi.Mux, route string) error {
				found := false
				chi.Walk(chiRoutes, func(method string,
					foundRoute string,
					handler http.Handler,
					middlewares ...func(http.Handler) http.Handler) error {
					if foundRoute == route {
						found = true
					}
					return nil
				})
				if !found {
					return errors.New("did not find the route")
				}
				return nil
			}
			testApp := Config{}
			testRoutes := testApp.NewHandler()
			chiRoutes := testRoutes.router
			routes := []string{"/auth", "/"}

			for _, route := range routes {
				Expect(walkRoutes(chiRoutes, route)).To(BeNil())
			}
		})
	})
})

var _ = Describe("Test handler", func() {
	When("Check User", func() {
		var userInterface *dbmocks.UserInterface
		var loggerInterface *loggermocks.LoggerInterface
		var testApp Config
		postBody := map[string]interface{}{
			"email":    "me@here.com",
			"password": "verysecret",
		}
		BeforeEach(func() {
			userInterface = dbmocks.NewUserInterface(GinkgoT())
			loggerInterface = loggermocks.NewLoggerInterface(GinkgoT())
			var db *sql.DB
			testApp = Config{
				DB:              db,
				userInterface:   userInterface,
				loggerInterface: loggerInterface,
			}
		})
		It("should fail if GetInfobyEmail fails", func() {
			mockGetInfoByEmail(userInterface, "me@here.com",
				&model.User{},
				errors.New("GetInfobyEmail failed"))
			body, _ := json.Marshal(postBody)
			req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.CheckUser)
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("should fail if password not match", func() {
			mockGetInfoByEmail(userInterface, "me@here.com",
				&model.User{
					ID:        "1",
					FirstName: "First",
					LastName:  "Last",
					Email:     "me@here.com",
					Password:  "",
					Active:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				nil)
			mockMatchPassword(userInterface, "verysecret", false)
			(*loggerInterface).On("HandleLog", "me@here.com", false).Return()
			body, _ := json.Marshal(postBody)
			req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.CheckUser)
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("should be successful if everything passes", func() {
			(*loggerInterface).On("HandleLog", "me@here.com", true).Return()
			mockGetInfoByEmail(userInterface, "me@here.com",
				&model.User{
					ID:        "1",
					FirstName: "First",
					LastName:  "Last",
					Email:     "me@here.com",
					Password:  "",
					Active:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				nil)
			mockMatchPassword(userInterface, "verysecret", true)
			body, _ := json.Marshal(postBody)
			req, _ := http.NewRequest("POST", "/auth", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.CheckUser)
			handler.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusAccepted))
		})
	})
})
