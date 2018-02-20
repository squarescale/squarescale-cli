package squarescale_test

import (
	"encoding/json"
	"github.com/onsi/gomega/ghttp"
	. "github.com/squarescale/squarescale-cli/squarescale"
	"net/http"
)

var _ = Describe("Environment", func() {
	var environment Environment

	BeforeEach(func() {
		environment = Environment{
			Project: &VariableGroup{
				Name: "Project",
				Variables: []*Variable{
					{Key: "DB_NAME", Value: "customprojectdb", Predefined: true},
					{Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false},
				},
			},
			Services: []*VariableGroup{
				{
					Name: "wordpress",
					Variables: []*Variable{
						{Key: "DB_NAME", Value: "dbwordpress", Predefined: true},
						{Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false},
						{Key: "WORDPRESS_DB", Value: "wordpress", Predefined: false},
					},
				},
				{
					Name: "rabbitmq",
					Variables: []*Variable{
						{Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false},
						{Key: "DB_NAME", Value: "rabbit", Predefined: false},
					},
				},
			},
		}
	})

	Describe("NewEnvironment", func() {
		var (
			server     *ghttp.Server
			client     *Client
			project    string
			statusCode int
			response   interface{}
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

		Context("when the project does not exist", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("shouldn't return an Environment", func() {
				env, err := NewEnvironment(client, project)

				Expect(env).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the project is found", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				response = map[string]interface{}{
					"project":     map[string]string{},
					"per_service": map[string]map[string]string{},
				}
			})

			It("should return a pointer to an Environment", func() {
				env, err := NewEnvironment(client, project)

				Expect(err).To(Not(HaveOccurred()))
				Expect(env).To(BeAssignableToTypeOf(&Environment{}))
			})
		})

		Context("when the server sends an unexpected status code", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("shouldn't return an Environment", func() {
				env, err := NewEnvironment(client, project)

				Expect(env).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UnmarshalJSON", func() {
		var (
			jsonBody []byte
			env      *Environment
		)

		JustBeforeEach(func() {
			env = &Environment{}
			json.Unmarshal(jsonBody, env)
		})

		Describe("when API response contains undesired attributes", func() {
			BeforeEach(func() {
				jsonBody = []byte(`{
				  "default": {"DB_NAME":"dbstaging"},
				  "global": {"MY_CUSTOM_GLOBAL":"one"},
				  "per_service": {
					"wordpress": {
					  "WORDPRESS_DB":"wordpress",
					  "default": {},
					  "custom": {}
					}
				  },
				  "project": {
					"default": {},
					"custom": {
					}
				  }
				}`)
			})

			It("ignores them", func() {
				Expect(env.Project.Variables).To(BeEmpty())
				Expect(env.Services[0].Variables).To(BeEmpty())
			})
		})

		Describe("when API response contains desired attributes", func() {
			BeforeEach(func() {
				jsonBody = []byte(`{
				  "per_service": {
					"wordpress": {
					  "default": {"DB_NAME":"dbwordpress"},
					  "custom": {"WORDPRESS_DB":"wordpress"}
					},
					"rabbitmq": {
					  "default": {},
					  "custom": {"DB_NAME":"rabbit"}
					}
				  },
				  "project": {
					"default": {"DB_NAME":"dbstaging"},
					"custom": {
					  "DB_NAME":"customprojectdb",
					  "MY_CUSTOM_GLOBAL":"one"
					}
				  }
				}`)
			})

			It("folds Project wide variables", func() {
				customDBVar := &Variable{
					Key: "DB_NAME", Value: "customprojectdb", Predefined: false,
				}
				Expect(env.Project.Variables).To(SatisfyAll(
					HaveLen(2),
					ContainElement(customDBVar)))
			})

			It("merges Project wide variables into each service", func() {
				expectedProjectVar := &Variable{
					Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false,
				}

				for _, service := range env.Services {
					Expect(service.Variables).To(ContainElement(expectedProjectVar))
				}
			})

			It("folds service variables", func() {
				Expect(env.Services[0].Variables).To(ContainElement(&Variable{
					Key: "DB_NAME", Value: "dbwordpress", Predefined: true,
				}))
				Expect(env.Services[1].Variables).To(ContainElement(&Variable{
					Key: "DB_NAME", Value: "rabbit", Predefined: false,
				}))
			})
		})
	})

	Describe("JSONObject", func() {
		var (
			jsonObject JSONObject
			err        error
		)

		JustBeforeEach(func() {
			jsonObject, err = environment.JSONObject()
		})

		Context("when the environment can be marshaled into JSON", func() {
			It("includes global and per_service environment variables", func() {
				Expect(jsonObject).To(SatisfyAll(
					HaveKey("global"),
					HaveKey("per_service")))
			})

			It("does not include predefined variables", func() {
				expectedJSON := map[string]string{"MY_CUSTOM_GLOBAL": "one"}

				Expect(jsonObject["global"]).To(Equal(expectedJSON))
			})

			It("includes variables for every service", func() {
				Expect(jsonObject["per_service"]).To(HaveLen(len(environment.Services)))
			})

			It("doesn't return any errors", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Context("when the environment cannot be marshaled into JSON", func() {
			BeforeEach(func() {
				environment = Environment{
					Project: &VariableGroup{
						Name:      "Project",
						Variables: []*Variable{nil},
					},
					Services: []*VariableGroup{},
				}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetServiceGroup", func() {
		var (
			groupName string
			group     *VariableGroup
			err       error
		)

		JustBeforeEach(func() {
			group, err = environment.GetServiceGroup(groupName)
		})

		Context("when the requested service exists", func() {
			BeforeEach(func() {
				groupName = environment.Services[0].Name
			})

			It("returns a pointer to the requested service", func() {
				Expect(group).To(BeIdenticalTo(environment.Services[0]))
			})

			It("doesn't return any errors", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Context("when the requested service does not exist", func() {
			BeforeEach(func() {
				groupName = "nonexistent"
			})

			It("does not return a pointer to a VariableGroup", func() {
				Expect(group).To(BeNil())
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("VariableGroup", func() {
	var variableGroup VariableGroup

	BeforeEach(func() {
		variableGroup = VariableGroup{Name: "Group"}
	})

	Describe("GetVariable", func() {
		var (
			variableName = "existing"
			preset       *Variable
			custom       *Variable
			variable     *Variable
			err          error
		)

		JustBeforeEach(func() {
			variable, err = variableGroup.GetVariable(variableName)
		})

		Context("when variable exists as predefined and custom in the group", func() {
			BeforeEach(func() {
				preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
				custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
				variableGroup.Variables = []*Variable{preset, custom}
			})

			It("returns a pointer to the custom variable", func() {
				Expect(variable).To(BeIdenticalTo(custom))
			})

			It("doesn't return any errors", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Context("when variable is defined once in the group", func() {
			BeforeEach(func() {
				preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
				variableGroup.Variables = []*Variable{preset}
			})

			It("returns a pointer to that variable", func() {
				Expect(variable).To(BeIdenticalTo(preset))
			})

			It("doesn't return any errors", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Context("when variable is not defined in the group", func() {
			BeforeEach(func() {
				preset = &Variable{Key: "not-existing-1", Value: "Preset", Predefined: true}
				custom = &Variable{Key: "not-existing-2", Value: "Custom", Predefined: false}
				variableGroup.Variables = []*Variable{preset, custom}
			})

			It("doesn't return a pointer to a variable", func() {
				Expect(variable).To(BeNil())
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("SetVariable", func() {
		var (
			variableName = "existing"
			preset       *Variable
			custom       *Variable
		)

		Context("when the variable is already defined in the group", func() {
			Context("and it is predefined", func() {
				BeforeEach(func() {
					preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
					variableGroup.Variables = []*Variable{preset}
				})

				It("doesn't update its value", func() {
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(preset.Value).To(Not(Equal("NEW_VALUE")))
				})

				It("creates a new custom variable with the given key and value", func() {
					length := len(variableGroup.Variables)
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(variableGroup.Variables).To(HaveLen(length + 1))
					Expect(variableGroup.Variables[length]).To(BeEquivalentTo(
						&Variable{Key: variableName, Value: "NEW_VALUE", Predefined: false},
					))
				})
			})

			Context("and it is custom", func() {
				BeforeEach(func() {
					custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
					variableGroup.Variables = []*Variable{custom}
				})

				It("updates its value", func() {
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(custom.Value).To(Equal("NEW_VALUE"))
				})

				It("doesn't create a new variable", func() {
					length := len(variableGroup.Variables)
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(variableGroup.Variables).To(HaveLen(length))
				})
			})

			Context("and it is both predefined and custom", func() {
				BeforeEach(func() {
					preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
					custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
					variableGroup.Variables = []*Variable{preset, custom}
				})

				It("updates the value of the custom definition only", func() {
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(custom.Value).To(Equal("NEW_VALUE"))
					Expect(preset.Value).To(Not(Equal("NEW_VALUE")))
				})

				It("doesn't create a new variable", func() {
					length := len(variableGroup.Variables)
					variableGroup.SetVariable(variableName, "NEW_VALUE")

					Expect(variableGroup.Variables).To(HaveLen(length))
				})
			})
		})

		Context("when the variable is not defined in the group", func() {
			It("creates a new custom variable with the given key and value", func() {
				length := len(variableGroup.Variables)
				variableGroup.SetVariable(variableName, "NEW_VALUE")

				Expect(variableGroup.Variables).To(HaveLen(length + 1))
				Expect(variableGroup.Variables[length]).To(BeEquivalentTo(
					&Variable{Key: variableName, Value: "NEW_VALUE", Predefined: false},
				))
			})
		})
	})

	Describe("RemoveVariable", func() {
		var (
			variableName = "existing"
			preset       *Variable
			custom       *Variable
		)

		Context("when the variable is defined in the group", func() {
			Context("and it is predefined", func() {
				BeforeEach(func() {
					preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
					variableGroup.Variables = []*Variable{preset}
				})

				It("does not remove it from the group", func() {
					variableGroup.RemoveVariable(variableName)

					Expect(variableGroup.Variables).To(ContainElement(preset))
				})

				It("returns an error", func() {
					err := variableGroup.RemoveVariable(variableName)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("and it is custom", func() {
				BeforeEach(func() {
					custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
					variableGroup.Variables = []*Variable{custom}
				})

				It("removes it from the group", func() {
					variableGroup.RemoveVariable(variableName)

					Expect(variableGroup.Variables).To(Not(ContainElement(custom)))
				})

				It("doesn't return an error", func() {
					err := variableGroup.RemoveVariable(variableName)

					Expect(err).To(Not(HaveOccurred()))
				})
			})

			Context("and it is both predefined and custom", func() {
				BeforeEach(func() {
					preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
					custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
					variableGroup.Variables = []*Variable{preset, custom}
				})

				It("removes the custom definition only", func() {
					variableGroup.RemoveVariable(variableName)

					Expect(variableGroup.Variables).To(Not(ContainElement(custom)))
					Expect(variableGroup.Variables).To(ContainElement(preset))
				})

				It("doesn't return an error", func() {
					err := variableGroup.RemoveVariable(variableName)

					Expect(err).To(Not(HaveOccurred()))
				})
			})
		})

		Context("when the variable is not defined in the group", func() {
			BeforeEach(func() {
				preset = &Variable{Key: variableName, Value: "Preset", Predefined: true}
				custom = &Variable{Key: variableName, Value: "Custom", Predefined: false}
				variableGroup.Variables = []*Variable{preset, custom}
			})

			It("doesn't remove any variables", func() {
				length := len(variableGroup.Variables)
				variableGroup.RemoveVariable("nonexistent")

				Expect(variableGroup.Variables).To(HaveLen(length))
			})

			It("doesn't return an error", func() {
				err := variableGroup.RemoveVariable("nonexistent")

				Expect(err).To(Not(HaveOccurred()))
			})
		})
	})

	Describe("CommitEnvironment", func() {
		var (
			server       *ghttp.Server
			client       *Client
			project      string
			err          error
			environment  *Environment
			expectedBody string
			statusCode   int
			response     interface{}
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			client = NewClient(server.URL(), "token")
			project = "whatever"

			environment = &Environment{
				Project: &VariableGroup{
					Name: "Project",
					Variables: []*Variable{
						{Key: "DB_NAME", Value: "dbstaging", Predefined: true},
						{Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false},
					},
				},
				Services: []*VariableGroup{
					{
						Name: "wordpress",
						Variables: []*Variable{
							{Key: "WORDPRESS_DB", Value: "wordpress", Predefined: false},
						},
					},
					{
						Name: "rabbitmq",
						Variables: []*Variable{
							{Key: "DB_NAME", Value: "rabbit", Predefined: false},
						},
					},
				},
			}
			expectedBody = `{
			  "environment":{
				"global":{
				  "MY_CUSTOM_GLOBAL":"one"
				},
				"per_service":{
				  "wordpress":{
					"WORDPRESS_DB":"wordpress"
				  },
				  "rabbitmq":{
					"DB_NAME":"rabbit"
				  }
				}
			  },
			  "format":"json"
			  }`

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/projects/"+project+"/environment/custom"),
				ghttp.VerifyJSON(expectedBody),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &response),
			))
		})

		JustBeforeEach(func() {
			err = environment.CommitEnvironment(client, project)

		})

		AfterEach(func() {
			server.Close()
		})

		It("makes a request to update the custom environment variables", func() {
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the environment cannot be marshaled into JSON", func() {
			BeforeEach(func() {
				environment = &Environment{
					Project: &VariableGroup{
						Name:      "Project",
						Variables: []*Variable{nil},
					},
					Services: []*VariableGroup{},
				}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the response is a success", func() {
			BeforeEach(func() {
				statusCode = http.StatusNoContent
			})

			It("doesn't return any errors", func() {
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Context("when the project is not found", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when request content was incorrect", func() {
			BeforeEach(func() {
				statusCode = http.StatusUnprocessableEntity
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the server sends an unexpected status code", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("shouldn't return an Environment", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
