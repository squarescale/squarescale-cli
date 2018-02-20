package squarescale_test

import (
	"github.com/onsi/gomega/ghttp"
	. "github.com/squarescale/squarescale-cli/squarescale"
	"net/http"
)

var _ = Describe("NewEnvironment", func() {
	var (
		server     *ghttp.Server
		client     *Client
		project    string
		statusCode int
		response   interface{}
		env        *Environment
		err        error
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL(), "token")
		project = "whatever"

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/projects/"+project+"/environment"),
			ghttp.VerifyContentType("application/json"),
			ghttp.RespondWithJSONEncodedPtr(&statusCode, &response),
		))
	})

	AfterEach(func() {
		server.Close()
	})

	It("should make a request to fetch environment variables", func() {
		NewEnvironment(client, project)
		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	Describe("when API response object matches the Environment struct", func() {
		BeforeEach(func() {
			statusCode = http.StatusOK
			response = map[string]interface{}{
				"default": map[string]string{"DB_NAME": "dbstaging"},
				"global":  map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				"per_service": map[string]map[string]string{
					"wordpress": map[string]string{"WORDPRESS_DB": "wordpress"},
				},
			}
		})

		It("can unmarshal the response into an Environment", func() {
			expectedEnv := &Environment{
				Preset: map[string]string{"DB_NAME": "dbstaging"},
				Global: map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				PerService: map[string]map[string]string{
					"wordpress": map[string]string{"WORDPRESS_DB": "wordpress"},
				},
			}
			env, err = NewEnvironment(client, "whatever")

			Expect(env).To(Equal(expectedEnv))
		})
	})

	Describe("when API response contains keys not matched in the Environment struct", func() {
		BeforeEach(func() {
			statusCode = http.StatusOK
			response = map[string]interface{}{
				"default": map[string]string{"DB_NAME": "dbstaging"},
				"global":  map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				"per_service": map[string]map[string]interface{}{
					"wordpress": map[string]interface{}{
						"WORDPRESS_DB": "wordpress",
						"default":      map[string]string{},
						"custom":       map[string]string{"WORDPRESS_DB": "wordpress"},
					},
				},
				"project": map[string]map[string]string{
					"default": map[string]string{"DB_NAME": "dbstaging"},
					"custom":  map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				},
			}
		})

		It("can unmarshal the response into an Environment", func() {
			expectedEnv := &Environment{
				Preset: map[string]string{"DB_NAME": "dbstaging"},
				Global: map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				PerService: map[string]map[string]string{
					"wordpress": map[string]string{"WORDPRESS_DB": "wordpress"},
				},
			}
			env, err = NewEnvironment(client, "whatever")

			Expect(env).To(Equal(expectedEnv))
		})
	})
})
