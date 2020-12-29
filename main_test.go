package main

import (
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/adams-sarah/test2doc/test"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
)

type DocsSuite struct {
	suite.Suite
	router     *chi.Mux
	testServer *test.Server
}

func (suite *DocsSuite) SetupSuite() {
	suite.router = chi.NewRouter()
	setupRoutes(suite.router)

	var err error
	suite.testServer, err = test.NewServer(suite.router)
	if err != nil {
		panic(err.Error())
	}
}

func (suite *DocsSuite) TearDownSuite() {
	suite.testServer.Finish()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDocs(t *testing.T) {
	suite.Run(t, new(DocsSuite))
}

func (suite *DocsSuite) TestPotresi() {

	resp, err := suite.Get("potresi.json")
	suite.Require().Nil(err)
	suite.assertJSONOK(resp)
}

func (suite *DocsSuite) TestVreme() {

	resp, err := suite.Get("postaje.json")
	suite.Require().Nil(err)
	suite.assertJSONOK(resp)
}

func (suite *DocsSuite) Get(uri string) (*http.Response, error) {
	u, _ := url.Parse(suite.testServer.URL)
	u.Path = path.Join(u.Path, uri)
	return http.Get(u.String())
}

func (suite *DocsSuite) assertJSONOK(resp *http.Response) {
	suite.Equal(http.StatusOK, resp.StatusCode, "HTTP status")
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "HTTP Content-Type")
}
