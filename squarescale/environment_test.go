package squarescale_test

import (
	. "github.com/squarescale/squarescale-cli/squarescale"
)

type mockedAPI struct {
	statusCode int
	response   []byte
	err        error
}

func (api *mockedAPI) Get(path string) (int, []byte, error) {
	return api.statusCode, api.response, api.err
}

var _ = Describe("NewEnvironment", func() {
	var (
		client        mockedAPI
		apiStatusCode int
		apiResponse   []byte
		apiError      error
		env           *Environment
		err           error
	)

	JustBeforeEach(func() {
		client = mockedAPI{
			statusCode: apiStatusCode,
			response:   apiResponse,
			err:        apiError,
		}
	})

	Describe("when API response object matches the Environment struct", func() {
		BeforeEach(func() {
			apiStatusCode = 200
			apiResponse = []byte(`{
			  "default": {"DB_NAME":"dbstaging"},
			  "global": {"MY_CUSTOM_GLOBAL":"one"},
			  "per_service": {
				"wordpress": {"WORDPRESS_DB":"wordpress"}
			  }
			}`)
			apiError = nil
		})

		It("can unmarshal the response into an Environment", func() {
			expectedEnv := &Environment{
				Preset: map[string]string{"DB_NAME": "dbstaging"},
				Global: map[string]string{"MY_CUSTOM_GLOBAL": "one"},
				PerService: map[string]map[string]string{
					"wordpress": map[string]string{"WORDPRESS_DB": "wordpress"},
				},
			}
			env, err = NewEnvironment(&client, "whatever")

			Expect(env).To(Equal(expectedEnv))
		})
	})
})
