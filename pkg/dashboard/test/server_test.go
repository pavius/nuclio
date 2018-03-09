/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nuclio/nuclio/pkg/dashboard"
	_ "github.com/nuclio/nuclio/pkg/dashboard/resource"
	"github.com/nuclio/nuclio/pkg/platform"
	"github.com/nuclio/nuclio/pkg/platformconfig"
	"github.com/nuclio/nuclio/pkg/restful"
	"github.com/nuclio/nuclio/test/compare"

	"github.com/nuclio/logger"
	"github.com/nuclio/zap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

//
// Platform mock
//

// Platform defines the interface that any underlying function platform must provide for nuclio
// to run over it
type mockPlatform struct {
	mock.Mock
}

// Build will locally build a processor image and return its name (or the error)
func (mp *mockPlatform) CreateFunctionBuild(createFunctionBuildOptions *platform.CreateFunctionBuildOptions) (*platform.CreateFunctionBuildResult, error) {
	args := mp.Called(createFunctionBuildOptions)
	return args.Get(0).(*platform.CreateFunctionBuildResult), args.Error(1)
}

// Deploy will deploy a processor image to the platform (optionally building it, if source is provided)
func (mp *mockPlatform) CreateFunction(createFunctionOptions *platform.CreateFunctionOptions) (*platform.CreateFunctionResult, error) {
	args := mp.Called(createFunctionOptions)
	return args.Get(0).(*platform.CreateFunctionResult), args.Error(1)
}

// UpdateFunctionOptions will update a previously deployed function
func (mp *mockPlatform) UpdateFunction(updateFunctionOptions *platform.UpdateFunctionOptions) error {
	args := mp.Called(updateFunctionOptions)
	return args.Error(0)
}

// DeleteFunction will delete a previously deployed function
func (mp *mockPlatform) DeleteFunction(deleteFunctionOptions *platform.DeleteFunctionOptions) error {
	args := mp.Called(deleteFunctionOptions)
	return args.Error(0)
}

// CreateFunctionInvocation will invoke a previously deployed function
func (mp *mockPlatform) CreateFunctionInvocation(createFunctionInvocationOptions *platform.CreateFunctionInvocationOptions) (*platform.CreateFunctionInvocationResult, error) {
	args := mp.Called(createFunctionInvocationOptions)
	return args.Get(0).(*platform.CreateFunctionInvocationResult), args.Error(1)
}

// CreateFunctionInvocation will invoke a previously deployed function
func (mp *mockPlatform) GetFunctions(getFunctionsOptions *platform.GetFunctionsOptions) ([]platform.Function, error) {
	args := mp.Called(getFunctionsOptions)
	return args.Get(0).([]platform.Function), args.Error(1)
}


// Deploy will deploy a processor image to the platform (optionally building it, if source is provided)
func (mp *mockPlatform) CreateProject(createProjectOptions *platform.CreateProjectOptions) error {
	args := mp.Called(createProjectOptions)
	return args.Error(0)
}

// UpdateProjectOptions will update a previously deployed function
func (mp *mockPlatform) UpdateProject(updateProjectOptions *platform.UpdateProjectOptions) error {
	args := mp.Called(updateProjectOptions)
	return args.Error(0)
}

// DeleteProject will delete a previously deployed function
func (mp *mockPlatform) DeleteProject(deleteProjectOptions *platform.DeleteProjectOptions) error {
	args := mp.Called(deleteProjectOptions)
	return args.Error(0)
}

// CreateProjectInvocation will invoke a previously deployed function
func (mp *mockPlatform) GetProjects(getProjectsOptions *platform.GetProjectsOptions) ([]platform.Project, error) {
	args := mp.Called(getProjectsOptions)
	return args.Get(0).([]platform.Project), args.Error(1)
}

// GetDeployRequiresRegistry returns true if a registry is required for deploy, false otherwise
func (mp *mockPlatform) GetDeployRequiresRegistry() bool {
	args := mp.Called()
	return args.Bool(0)
}

// GetName returns the platform name
func (mp *mockPlatform) GetName() string {
	args := mp.Called()
	return args.String(0)
}

// GetNodes returns a slice of nodes currently in the cluster
func (mp *mockPlatform) GetNodes() ([]platform.Node, error) {
	args := mp.Called()
	return args.Get(0).([]platform.Node), args.Error(1)
}

//
// Test suite
//

type dashboardTestSuite struct {
	suite.Suite
	logger          logger.Logger
	dashboardServer *dashboard.Server
	httpServer      *httptest.Server
	mockPlatform    *mockPlatform
}

func (suite *dashboardTestSuite) SetupTest() {
	var err error
	trueValue := true

	suite.logger, _ = nucliozap.NewNuclioZapTest("test")
	suite.mockPlatform = &mockPlatform{}

	// create a mock platform
	suite.dashboardServer, err = dashboard.NewServer(suite.logger,
		"",
		"",
		"",
		"",
		suite.mockPlatform,
		true,
		&platformconfig.WebServer{Enabled: &trueValue},
		nil)

	if err != nil {
		panic("Failed to create server")
	}

	// create an http server from the dashboard server
	suite.httpServer = httptest.NewServer(suite.dashboardServer.Router)
}

func (suite *dashboardTestSuite) TeardownTest() {
	suite.httpServer.Close()
}

func (suite *dashboardTestSuite) sendRequest(method string,
	path string,
	requestHeaders map[string]string,
	requestBody io.Reader,
	expectedStatusCode *int,
	encodedExpectedResponse interface{}) (*http.Response, map[string]interface{}) {

	request, err := http.NewRequest(method, suite.httpServer.URL+path, requestBody)
	suite.Require().NoError(err)

	for headerKey, headerValue := range requestHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	response, err := http.DefaultClient.Do(request)
	suite.Require().NoError(err)

	encodedResponseBody, err := ioutil.ReadAll(response.Body)
	suite.Require().NoError(err)

	defer response.Body.Close()

	suite.logger.DebugWith("Got response",
		"status", response.StatusCode,
		"response", string(encodedResponseBody))

	// check if status code was passed
	if expectedStatusCode != nil {
		suite.Require().Equal(*expectedStatusCode, response.StatusCode)
	}

	// if there's an expected status code, verify it
	decodedResponseBody := map[string]interface{}{}

	if encodedExpectedResponse != nil {

		err = json.Unmarshal(encodedResponseBody, &decodedResponseBody)
		suite.Require().NoError(err)

		suite.logger.DebugWith("Comparing expected", "expected", encodedExpectedResponse)

		switch typedEncodedExpectedResponse := encodedExpectedResponse.(type) {
		case string:
			decodedExpectedResponseBody := map[string]interface{}{}

			err = json.Unmarshal([]byte(typedEncodedExpectedResponse), &decodedExpectedResponseBody)
			suite.Require().NoError(err)

			suite.Require().True(compare.CompareNoOrder(decodedExpectedResponseBody, decodedResponseBody))

		case func(response map[string]interface{}) bool:
			suite.Require().True(typedEncodedExpectedResponse(decodedResponseBody))

		default:
			panic("Unsupported expected response verifier")
		}
	}

	return response, decodedResponseBody
}

//
// Function
//

type functionTestSuite struct {
	dashboardTestSuite
}

func (suite *functionTestSuite) TestGetDetailSuccessful() {
	returnedFunction := platform.AbstractFunction{}
	returnedFunction.Config.Meta.Name = "f1"
	returnedFunction.Config.Meta.Namespace = "f1Namespace"
	returnedFunction.Config.Spec.Replicas = 10

	// verify
	verifyGetFunctions := func(getFunctionsOptions *platform.GetFunctionsOptions) bool {
		suite.Require().Equal("f1", getFunctionsOptions.Name)
		suite.Require().Equal("f1Namespace", getFunctionsOptions.Namespace)

		return true
	}

	suite.mockPlatform.
		On("GetFunctions", mock.MatchedBy(verifyGetFunctions)).
		Return([]platform.Function{&returnedFunction}, nil).
		Once()

	headers := map[string]string{
		"x-nuclio-function-namespace": "f1Namespace",
	}

	expectedStatusCode := http.StatusOK
	expectedResponseBody := `{
	"metadata": {
		"name": "f1",
		"namespace": "f1Namespace"
	},
	"spec": {
		"resources": {},
		"build": {},
		"replicas": 10
	}
}`

	suite.sendRequest("GET",
		"/functions/f1",
		headers,
		nil,
		&expectedStatusCode,
		expectedResponseBody)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestGetDetailNoNamespace() {
	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Namespace must exist"})
	suite.sendRequest("GET",
		"/functions/f1",
		nil,
		nil,
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestGetListSuccessful() {
	returnedFunction1 := platform.AbstractFunction{}
	returnedFunction1.Config.Meta.Name = "f1"
	returnedFunction1.Config.Meta.Namespace = "fNamespace"
	returnedFunction1.Config.Spec.Runtime = "r1"

	returnedFunction2 := platform.AbstractFunction{}
	returnedFunction2.Config.Meta.Name = "f2"
	returnedFunction2.Config.Meta.Namespace = "fNamespace"
	returnedFunction2.Config.Spec.Runtime = "r2"

	// verify
	verifyGetFunctions := func(getFunctionsOptions *platform.GetFunctionsOptions) bool {
		suite.Require().Equal("", getFunctionsOptions.Name)
		suite.Require().Equal("fNamespace", getFunctionsOptions.Namespace)

		return true
	}

	suite.mockPlatform.
		On("GetFunctions", mock.MatchedBy(verifyGetFunctions)).
		Return([]platform.Function{&returnedFunction1, &returnedFunction2}, nil).
		Once()

	headers := map[string]string{
		"x-nuclio-function-namespace": "fNamespace",
	}

	expectedStatusCode := http.StatusOK
	expectedResponseBody := `{
	"f1": {
		"metadata": {
			"name": "f1",
			"namespace": "fNamespace"
		},
		"spec": {
			"resources": {},
			"build": {},
			"runtime": "r1"
		}
	},
	"f2": {
		"metadata": {
			"name": "f2",
			"namespace": "fNamespace"
		},
		"spec": {
			"resources": {},
			"build": {},
			"runtime": "r2"
		}
	}
}`

	suite.sendRequest("GET",
		"/functions",
		headers,
		nil,
		&expectedStatusCode,
		expectedResponseBody)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestGetListNoNamespace() {
	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Namespace must exist"})
	suite.sendRequest("GET",
		"/functions",
		nil,
		nil,
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestCreateSuccessful() {

	// verify
	verifyCreateFunction := func(createFunctionOptions *platform.CreateFunctionOptions) bool {
		suite.Require().Equal("f1", createFunctionOptions.FunctionConfig.Meta.Name)
		suite.Require().Equal("f1Namespace", createFunctionOptions.FunctionConfig.Meta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("CreateFunction", mock.MatchedBy(verifyCreateFunction)).
		Return(&platform.CreateFunctionResult{}, nil).
		Once()

	headers := map[string]string{
		"x-nuclio-wait-function-action": "true",
	}

	expectedStatusCode := http.StatusAccepted
	requestBody := `{
	"metadata": {
		"name": "f1",
		"namespace": "f1Namespace"
	},
	"spec": {
		"resources": {},
		"build": {},
		"runtime": "r1"
	}
}`

	suite.sendRequest("POST",
		"/functions",
		headers,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		nil)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestCreateNoMetadata() {
	suite.sendRequestNoMetadata("POST")
}

func (suite *functionTestSuite) TestCreateNoName() {
	suite.sendRequestNoName("POST")
}

func (suite *functionTestSuite) TestCreateNoNamespace() {
	suite.sendRequestNoNamespace("POST")
}

func (suite *functionTestSuite) TestUpdateSuccessful() {

	// verify
	verifyUpdateFunction := func(updateFunctionOptions *platform.UpdateFunctionOptions) bool {
		suite.Require().Equal("f1", updateFunctionOptions.FunctionMeta.Name)
		suite.Require().Equal("f1Namespace", updateFunctionOptions.FunctionMeta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("UpdateFunction", mock.MatchedBy(verifyUpdateFunction)).
		Return(nil).
		Once()

	headers := map[string]string{
		"x-nuclio-wait-function-action": "true",
	}

	expectedStatusCode := http.StatusAccepted
	requestBody := `{
	"metadata": {
		"name": "f1",
		"namespace": "f1Namespace"
	},
	"spec": {
		"resources": {},
		"build": {},
		"runtime": "r1"
	}
}`

	suite.sendRequest("PUT",
		"/functions",
		headers,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		nil)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestUpdateNoMetadata() {
	suite.sendRequestNoMetadata("PUT")
}

func (suite *functionTestSuite) TestUpdateNoName() {
	suite.sendRequestNoName("PUT")
}

func (suite *functionTestSuite) TestUpdateNoNamespace() {
	suite.sendRequestNoNamespace("PUT")
}

func (suite *functionTestSuite) TestDeleteSuccessful() {

	// verify
	verifyDeleteFunction := func(deleteFunctionOptions *platform.DeleteFunctionOptions) bool {
		suite.Require().Equal("f1", deleteFunctionOptions.FunctionConfig.Meta.Name)
		suite.Require().Equal("f1Namespace", deleteFunctionOptions.FunctionConfig.Meta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("DeleteFunction", mock.MatchedBy(verifyDeleteFunction)).
		Return(nil).
		Once()

	headers := map[string]string{
		"x-nuclio-wait-function-action": "true",
	}

	expectedStatusCode := http.StatusNoContent
	requestBody := `{
	"metadata": {
		"name": "f1",
		"namespace": "f1Namespace"
	}
}`

	suite.sendRequest("DELETE",
		"/functions",
		headers,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		nil)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestDeleteNoMetadata() {
	suite.sendRequestNoMetadata("DELETE")
}

func (suite *functionTestSuite) TestDeleteNoName() {
	suite.sendRequestNoName("DELETE")
}

func (suite *functionTestSuite) TestDeleteNoNamespace() {
	suite.sendRequestNoNamespace("DELETE")
}

func (suite *functionTestSuite) TestInvokeSuccessful() {
	functionName := "f1"
	functionNamespace := "f1Namespace"

	requestMethod := "PUT"
	requestPath := "/some/path"
	requestBody := []byte("request body")
	responseBody := []byte(`{"response": "body"}`)

	// headers we want to pass to the actual function
	functionRequestHeaders := map[string]string{
		"request_h1": "request_h1v",
		"request_h2": "request_h2v",
	}

	// headers we need to pass to dashboard for invocation
	requestHeaders := map[string]string{
		"x-nuclio-path":               requestPath,
		"x-nuclio-function-name":      functionName,
		"x-nuclio-function-namespace": functionNamespace,
		"x-nuclio-invoke-via":         "external-ip",
	}

	// add functionRequestHeaders to requestHeaders so that dashboard will invoke the functions with them
	for headerKey, headerValue := range functionRequestHeaders {
		requestHeaders[headerKey] = headerValue
	}

	// CreateFunctionInvocationResult holds the result of a single invocation
	expectedInvokeResult := platform.CreateFunctionInvocationResult{
		Headers: map[string][]string{
			"response_h1": {"response_h1v"},
			"response_h2": {"response_h2v"},
		},
		Body:       responseBody,
		StatusCode: http.StatusCreated,
	}

	// verify call to invoke function
	verifyCreateFunctionInvocation := func(createFunctionInvocationOptions *platform.CreateFunctionInvocationOptions) bool {
		suite.Require().Equal(functionName, createFunctionInvocationOptions.Name)
		suite.Require().Equal(functionNamespace, createFunctionInvocationOptions.Namespace)
		suite.Require().Equal(requestBody, createFunctionInvocationOptions.Body)
		suite.Require().Equal(requestMethod, createFunctionInvocationOptions.Method)
		suite.Require().Equal(platform.InvokeViaExternalIP, createFunctionInvocationOptions.Via)

		// dashboard will trim the first "/"
		suite.Require().Equal(requestPath[1:], createFunctionInvocationOptions.Path)

		// expect only to receive the function headers (those that don't start with x-nuclio
		for headerKey, _ := range createFunctionInvocationOptions.Headers {
			suite.Require().False(strings.HasPrefix(headerKey, "x-nuclio"))
		}

		// expect all the function headers to be there
		for headerKey, headerValue := range functionRequestHeaders {
			suite.Require().Equal(headerValue, createFunctionInvocationOptions.Headers.Get(headerKey))
		}

		return true
	}

	suite.mockPlatform.
		On("CreateFunctionInvocation", mock.MatchedBy(verifyCreateFunctionInvocation)).
		Return(&expectedInvokeResult, nil).
		Once()

	expectedStatusCode := expectedInvokeResult.StatusCode

	suite.sendRequest(requestMethod,
		"/function_invocations",
		requestHeaders,
		bytes.NewBuffer(requestBody),
		&expectedStatusCode,
		string(responseBody))

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestInvokeNoName() {

	// headers we need to pass to dashboard for invocation
	requestHeaders := map[string]string{
		"x-nuclio-path":               "p",
		"x-nuclio-function-namespace": "ns",
		"x-nuclio-invoke-via":         "external-ip",
	}

	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Function name and namespace must be provided"})

	expectedStatusCode := http.StatusBadRequest
	suite.sendRequest("POST",
		"/function_invocations",
		requestHeaders,
		bytes.NewBufferString("request body"),
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) TestInvokeNoNamespace() {

	// headers we need to pass to dashboard for invocation
	requestHeaders := map[string]string{
		"x-nuclio-path":          "p",
		"x-nuclio-function-name": "n",
		"x-nuclio-invoke-via":    "external-ip",
	}

	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Function name and namespace must be provided"})

	expectedStatusCode := http.StatusBadRequest
	suite.sendRequest("POST",
		"/function_invocations",
		requestHeaders,
		bytes.NewBufferString("request body"),
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *functionTestSuite) sendRequestNoMetadata(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"spec": {
		"resources": {},
		"build": {},
		"runtime": "r1"
	}
}`)
}

func (suite *functionTestSuite) sendRequestNoNamespace(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"metadata": {
		"name": "f1Name"
	},
	"spec": {
		"resources": {},
		"build": {},
		"runtime": "r1"
	}
}`)
}

func (suite *functionTestSuite) sendRequestNoName(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"metadata": {
		"namespace": "f1Namespace"
	},
	"spec": {
		"resources": {},
		"build": {},
		"runtime": "r1"
	}
}`)
}

func (suite *functionTestSuite) sendRequestWithInvalidBody(method string, body string) {
	headers := map[string]string{
		"x-nuclio-wait-function-action": "true",
	}

	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Function name and namespace must be provided in metadata"})
	requestBody := body

	suite.sendRequest(method,
		"/functions",
		headers,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

//
// Project
//

type projectTestSuite struct {
	dashboardTestSuite
}

func (suite *projectTestSuite) TestGetDetailSuccessful() {
	returnedProject := platform.AbstractProject{}
	returnedProject.ProjectConfig.Meta.Name = "p1"
	returnedProject.ProjectConfig.Meta.Namespace = "p1Namespace"
	returnedProject.ProjectConfig.Spec.DisplayName = "p1DisplayName"
	returnedProject.ProjectConfig.Spec.Description = "p1Desc"

	// verify
	verifyGetProjects := func(getProjectsOptions *platform.GetProjectsOptions) bool {
		suite.Require().Equal("p1", getProjectsOptions.Meta.Name)
		suite.Require().Equal("p1Namespace", getProjectsOptions.Meta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("GetProjects", mock.MatchedBy(verifyGetProjects)).
		Return([]platform.Project{&returnedProject}, nil).
		Once()

	headers := map[string]string{
		"x-nuclio-project-namespace": "p1Namespace",
	}

	expectedStatusCode := http.StatusOK
	expectedResponseBody := `{
	"metadata": {
		"name": "p1",
		"namespace": "p1Namespace"
	},
	"spec": {
		"displayName": "p1DisplayName",
		"description": "p1Desc"
	}
}`

	suite.sendRequest("GET",
		"/projects/p1",
		headers,
		nil,
		&expectedStatusCode,
		expectedResponseBody)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *projectTestSuite) TestGetDetailNoNamespace() {
	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Namespace must exist"})
	suite.sendRequest("GET",
		"/projects/p1",
		nil,
		nil,
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *projectTestSuite) TestGetListSuccessful() {
	returnedProject1 := platform.AbstractProject{}
	returnedProject1.ProjectConfig.Meta.Name = "p1"
	returnedProject1.ProjectConfig.Meta.Namespace = "pNamespace"
	returnedProject1.ProjectConfig.Spec.DisplayName = "p1DisplayName"
	returnedProject1.ProjectConfig.Spec.Description = "p1Desc"

	returnedProject2 := platform.AbstractProject{}
	returnedProject2.ProjectConfig.Meta.Name = "p2"
	returnedProject2.ProjectConfig.Meta.Namespace = "pNamespace"
	returnedProject2.ProjectConfig.Spec.DisplayName = "p2DisplayName"
	returnedProject2.ProjectConfig.Spec.Description = "p2Desc"

	// verify
	verifyGetProjects := func(getProjectsOptions *platform.GetProjectsOptions) bool {
		suite.Require().Equal("", getProjectsOptions.Meta.Name)
		suite.Require().Equal("pNamespace", getProjectsOptions.Meta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("GetProjects", mock.MatchedBy(verifyGetProjects)).
		Return([]platform.Project{&returnedProject1, &returnedProject2}, nil).
		Once()

	headers := map[string]string{
		"x-nuclio-project-namespace": "pNamespace",
	}

	expectedStatusCode := http.StatusOK
	expectedResponseBody := `{
	"p1": {
		"metadata": {
			"name": "p1",
			"namespace": "pNamespace"
		},
		"spec": {
			"displayName": "p1DisplayName",
			"description": "p1Desc"
		}
	},
	"p2": {
		"metadata": {
			"name": "p2",
			"namespace": "pNamespace"
		},
		"spec": {
			"displayName": "p2DisplayName",
			"description": "p2Desc"
		}
	}
}`

	suite.sendRequest("GET",
		"/projects",
		headers,
		nil,
		&expectedStatusCode,
		expectedResponseBody)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *projectTestSuite) TestGetListNoNamespace() {
	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Namespace must exist"})
	suite.sendRequest("GET",
		"/projects",
		nil,
		nil,
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}


func (suite *projectTestSuite) TestCreateSuccessful() {

	// verify
	verifyCreateProject := func(createProjectOptions *platform.CreateProjectOptions) bool {
		suite.Require().Equal("p1", createProjectOptions.ProjectConfig.Meta.Name)
		suite.Require().Equal("p1Namespace", createProjectOptions.ProjectConfig.Meta.Namespace)
		suite.Require().Equal("p1DisplayName", createProjectOptions.ProjectConfig.Spec.DisplayName)
		suite.Require().Equal("p1Description", createProjectOptions.ProjectConfig.Spec.Description)

		return true
	}

	suite.mockPlatform.
		On("CreateProject", mock.MatchedBy(verifyCreateProject)).
		Return(nil).
		Once()

	expectedStatusCode := http.StatusNoContent
	requestBody := `{
	"metadata": {
		"name": "p1",
		"namespace": "p1Namespace"
	},
	"spec": {
		"displayName": "p1DisplayName",
		"description": "p1Description"
	}
}`

	suite.sendRequest("POST",
		"/projects",
		nil,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		nil)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *projectTestSuite) TestCreateNoMetadata() {
	suite.sendRequestNoMetadata("POST")
}

func (suite *projectTestSuite) TestCreateNoName() {
	suite.sendRequestNoName("POST")
}

func (suite *projectTestSuite) TestCreateNoNamespace() {
	suite.sendRequestNoNamespace("POST")
}
//
//func (suite *projectTestSuite) TestUpdateSuccessful() {
//
//	// verify
//	verifyUpdateProject := func(updateProjectOptions *platform.UpdateProjectOptions) bool {
//		suite.Require().Equal("p1", updateProjectOptions.ProjectConfig.Meta.Name)
//		suite.Require().Equal("p1Namespace", updateProjectOptions.ProjectConfig.Meta.Namespace)
//		suite.Require().Equal("p1DisplayName", updateProjectOptions.ProjectConfig.Spec.DisplayName)
//		suite.Require().Equal("p1Description", updateProjectOptions.ProjectConfig.Spec.Description)
//
//		return true
//	}
//
//	suite.mockPlatform.
//		On("UpdateProject", mock.MatchedBy(verifyUpdateProject)).
//		Return(nil).
//		Once()
//
//	expectedStatusCode := http.StatusAccepted
//	requestBody := `{
//	"metadata": {
//		"name": "p1",
//		"namespace": "p1Namespace"
//	},
//	"spec": {
//		"displayName": "p1DisplayName",
//		"description": "p1Description"
//	}
//}`
//
//	suite.sendRequest("PUT",
//		"/projects",
//		nil,
//		bytes.NewBufferString(requestBody),
//		&expectedStatusCode,
//		nil)
//
//	suite.mockPlatform.AssertExpectations(suite.T())
//}
//
//func (suite *projectTestSuite) TestUpdateNoMetadata() {
//	suite.sendRequestNoMetadata("PUT")
//}
//
//func (suite *projectTestSuite) TestUpdateNoName() {
//	suite.sendRequestNoName("PUT")
//}
//
//func (suite *projectTestSuite) TestUpdateNoNamespace() {
//	suite.sendRequestNoNamespace("PUT")
//}

func (suite *projectTestSuite) TestDeleteSuccessful() {

	// verify
	verifyDeleteProject := func(deleteProjectOptions *platform.DeleteProjectOptions) bool {
		suite.Require().Equal("p1", deleteProjectOptions.Meta.Name)
		suite.Require().Equal("p1Namespace", deleteProjectOptions.Meta.Namespace)

		return true
	}

	suite.mockPlatform.
		On("DeleteProject", mock.MatchedBy(verifyDeleteProject)).
		Return(nil).
		Once()

	expectedStatusCode := http.StatusNoContent
	requestBody := `{
	"metadata": {
		"name": "p1",
		"namespace": "p1Namespace"
	}
}`

	suite.sendRequest("DELETE",
		"/projects",
		nil,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		nil)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func (suite *projectTestSuite) TestDeleteNoMetadata() {
	suite.sendRequestNoMetadata("DELETE")
}

func (suite *projectTestSuite) TestDeleteNoName() {
	suite.sendRequestNoName("DELETE")
}

func (suite *projectTestSuite) TestDeleteNoNamespace() {
	suite.sendRequestNoNamespace("DELETE")
}

func (suite *projectTestSuite) sendRequestNoMetadata(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"spec": {
		"displayName": "dn",
		"description": "d"
	}
}`)
}

func (suite *projectTestSuite) sendRequestNoNamespace(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"metadata": {
		"name": "name"
	},
	"spec": {
		"displayName": "dn",
		"description": "d"
	}
}`)
}

func (suite *projectTestSuite) sendRequestNoName(method string) {
	suite.sendRequestWithInvalidBody(method, `{
	"metadata": {
		"namespace": "namespace"
	},
	"spec": {
		"displayName": "dn",
		"description": "d"
	}
}`)
}

func (suite *projectTestSuite) sendRequestWithInvalidBody(method string, body string) {
	expectedStatusCode := http.StatusBadRequest
	ecv := restful.NewErrorContainsVerifier(suite.logger, []string{"Project name and namespace must be provided in metadata"})
	requestBody := body

	suite.sendRequest(method,
		"/projects",
		nil,
		bytes.NewBufferString(requestBody),
		&expectedStatusCode,
		ecv.Verify)

	suite.mockPlatform.AssertExpectations(suite.T())
}

func TestDashboardTestSuite(t *testing.T) {
	suite.Run(t, new(functionTestSuite))
	suite.Run(t, new(projectTestSuite))
}
