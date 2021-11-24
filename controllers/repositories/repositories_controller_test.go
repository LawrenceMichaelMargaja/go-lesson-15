 package repositories

 import (
	 "encoding/json"
	 "fmt"
	 "github.com/dembygenesis/go-rest-industry-standard/src/api/clients/restclient"
	 "github.com/dembygenesis/go-rest-industry-standard/src/api/domain/repositories"
	 "github.com/dembygenesis/go-rest-industry-standard/src/api/utils/errors"
	 "github.com/dembygenesis/go-rest-industry-standard/src/api/utils/test_utils"
	 "github.com/stretchr/testify/assert"
	 "io/ioutil"
	 "net/http"
	 "net/http/httptest"
	 "os"
	 "strings"
	 "testing"
 )

 func TestMain(m *testing.M) {
	 restclient.StartMockUps()
	 os.Exit(m.Run())
 }

 func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	 request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(``))
	 response := httptest.NewRecorder()
	 c := test_utils.GetMockedContext(request, response)
	 //c, _ := gin.CreateTestContext(response)
	 //c.Request = request

	 CreateRepo(c)

	 assert.EqualValues(t, http.StatusBadRequest, response.Code)
	 fmt.Println(response.Body.String())

	 apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	 assert.Nil(t, err)
	 assert.NotNil(t, apiErr)
	 assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	 assert.EqualValues(t, "invalid json body", apiErr.Message())
 }

 func TestCreateRepoErrorFromGithub(t *testing.T) {
	 restclient.FlushMockUps()
	 restclient.AddMockUp(restclient.Mock{
		 Url:        "https://api.github.com/user/repos",
		 HttpMethod: http.MethodPost,
		 Response:   &http.Response{
			 StatusCode: http.StatusUnauthorized,
			 Body: ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"}`)),
		 },
	 })

	 request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "testing"}`))
	 response := httptest.NewRecorder()
	 c := test_utils.GetMockedContext(request, response)



	 CreateRepo(c)

	 assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	 fmt.Println(response.Body.String())

	 apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	 assert.Nil(t, err)
	 assert.NotNil(t, apiErr)
	 assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	 assert.EqualValues(t, "Requires authentication", apiErr.Message())
 }


 func TestCreateRepoNoError(t *testing.T) {
	 restclient.FlushMockUps()
	 restclient.AddMockUp(restclient.Mock{
		 Url:        "https://api.github.com/user/repos",
		 HttpMethod: http.MethodPost,
		 Response:   &http.Response{
			 StatusCode: http.StatusCreated,
			 Body: ioutil.NopCloser(strings.NewReader(`{"id": 123}`)),
		 },
	 })

	 request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "testing"}`))
	 response := httptest.NewRecorder()
	 c := test_utils.GetMockedContext(request, response)


	 CreateRepo(c)

	 assert.EqualValues(t, http.StatusCreated, response.Code)
	 fmt.Println(response.Body.String())

	 var result repositories.CreateRepoResponse
	 err := json.Unmarshal(response.Body.Bytes(), &result)
	 assert.Nil(t, err)
	 assert.EqualValues(t, 123, result.Id)
	 assert.EqualValues(t, "", result.Name)
	 assert.EqualValues(t, "", result.Owner)
 }