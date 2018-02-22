package squarescale_test

import (
	"fmt"
	"github.com/onsi/gomega/types"
	. "github.com/squarescale/squarescale-cli/squarescale"
)

var _ = Describe("Environment", func() {
	var environment *Environment

	BeforeEach(func() {
		environment = &Environment{
			Project: &VariableGroup{
				Name: "Project",
				Variables: []*Variable{
					{Key: "MY_CUSTOM_GLOBAL", Value: "one", Predefined: false},
					{Key: "DB_NAME", Value: "custom_db", Predefined: false},
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
	})

	Describe("QueryVars", func() {
		var (
			options QueryOptions
			result  QueryResult
			err     error
		)

		BeforeEach(func() {
			options = QueryOptions{
				DisplayAll:   false,
				ServiceName:  "",
				VariableName: "",
			}
		})

		JustBeforeEach(func() {
			result, err = environment.QueryVars(options)
		})

		Context("when option DisplayAll is set to true", func() {
			BeforeEach(func() {
				options.DisplayAll = true
			})

			It("returns the entire environment", func() {
				Expect(result).To(BeEquivalentTo(environment))
			})

			Context("when option ServiceName is not empty", func() {
				BeforeEach(func() {
					options.ServiceName = environment.Services[0].Name
				})

				It("returns the entire environment", func() {
					Expect(result).To(BeEquivalentTo(environment))
				})

				Context("when option VariableName is not empty", func() {
					BeforeEach(func() {
						options.VariableName = environment.Services[0].Variables[0].Key
					})

					It("returns the entire environment", func() {
						Expect(result).To(BeEquivalentTo(environment))
					})
				})
			})

			Context("when option VariableName is not empty", func() {
				BeforeEach(func() {
					options.VariableName = environment.Services[0].Variables[0].Key
				})

				It("returns the entire environment", func() {
					Expect(result).To(BeEquivalentTo(environment))
				})
			})
		})

		Context("when option DisplayAll is set to false", func() {
			var (
				expectedVariableGroup *VariableGroup
				expectedVariable      *Variable
			)

			BeforeEach(func() {
				options.DisplayAll = false
			})

			Context("when option ServiceName is empty", func() {
				Context("when option VariableName is empty", func() {
					It("returns the Project group", func() {
						Expect(result).To(BeEquivalentTo(environment.Project))
					})
				})

				Context("when option VariableName corresponds to an existing variable", func() {
					BeforeEach(func() {
						expectedVariable = environment.Project.Variables[0]
						options.VariableName = expectedVariable.Key
					})

					It("returns the variable from the project group", func() {
						Expect(result).To(BeEquivalentTo(expectedVariable))
					})
				})

				Context("when option VariableName doesn't correspond to an existing variable", func() {
					BeforeEach(func() {
						options.VariableName = "nonexistent"
					})

					It("doesn't return any variables", func() {
						Expect(result).To(BeNil())
					})

					It("returns an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("when option ServiceName is not empty", func() {
				BeforeEach(func() {
					options.ServiceName = environment.Services[1].Name
					expectedVariableGroup = environment.Services[1]
				})

				Context("when option VariableName is empty", func() {
					It("returns the given service group", func() {
						Expect(result).To(BeEquivalentTo(expectedVariableGroup))
					})
				})

				Context("when option VariableName corresponds to an existing variable", func() {
					BeforeEach(func() {
						options.VariableName = environment.Services[1].Variables[0].Key
						expectedVariable = environment.Services[1].Variables[0]
					})

					It("returns the given variable", func() {
						Expect(result).To(BeEquivalentTo(expectedVariable))
					})
				})

				Context("when option VariableName doesn't correspond to an existing variable", func() {
					BeforeEach(func() {
						options.VariableName = "nonexistent"
					})

					It("doesn't return any variables", func() {
						Expect(result).To(BeNil())
					})

					It("returns an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("String", func() {
		It("writes the Project variable group name to the buffer", func() {
			Expect(environment.String()).To(ContainSubstring(environment.Project.Name))
		})

		It("writes each service variable group name to the buffer", func() {
			matchers := []types.GomegaMatcher{}
			for _, service := range environment.Services {
				matchers = append(matchers, ContainSubstring(service.Name))
			}

			Expect(environment.String()).To(SatisfyAll(matchers...))
		})

		Context("when the Project variable group contains variables", func() {
			It("writes those variables as key, value pairs to the buffer", func() {
				matchers := []types.GomegaMatcher{}
				for _, variable := range environment.Project.Variables {
					varRegexp := fmt.Sprintf("%s\\s+=\\s+%s", variable.Key, variable.Value)
					matchers = append(matchers, MatchRegexp(varRegexp))
				}

				Expect(environment.String()).To(SatisfyAll(matchers...))
			})

			It("doesn't write none to the buffer", func() {
				Expect(environment.String()).To(Not(MatchRegexp("none\n")))
			})
		})

		Context("when the Project variable group doesn't contain variables", func() {
			BeforeEach(func() {
				environment.Project.Variables = []*Variable{}
			})

			It("writes none to the buffer", func() {
				none := fmt.Sprintf("^%s\n  none\n", environment.Project.Name)
				Expect(environment.String()).To(MatchRegexp(none))
			})
		})

		Context("when a service variable group contains variables", func() {
			var service *VariableGroup

			BeforeEach(func() {
				service = environment.Services[0]
			})

			It("writes those variables as key, value pairs to the buffer", func() {
				matchers := []types.GomegaMatcher{}
				for _, variable := range service.Variables {
					varRegexp := fmt.Sprintf("%s\\s+=\\s+%s", variable.Key, variable.Value)
					matchers = append(matchers, MatchRegexp(varRegexp))
				}

				Expect(environment.String()).To(SatisfyAll(matchers...))
			})

			It("doesn't write none to the buffer", func() {
				Expect(environment.String()).To(Not(MatchRegexp("none\n")))
			})
		})

		Context("when a service variable group doesn't contain variables", func() {
			var service *VariableGroup

			BeforeEach(func() {
				service = environment.Services[0]
				service.Variables = []*Variable{}
			})

			It("writes none to the buffer", func() {
				none := fmt.Sprintf("\n%s\n  none\n", service.Name)
				Expect(environment.String()).To(MatchRegexp(none))
			})
		})
	})
})

var _ = Describe("VariableGroup", func() {
	var variableGroup VariableGroup

	Describe("String", func() {
		Context("when the group contains variables", func() {
			BeforeEach(func() {
				variableGroup = VariableGroup{
					Name: "Not empty",
					Variables: []*Variable{
						{Key: "KEY_1", Value: "value_1", Predefined: false},
						{Key: "KEY_2", Value: "value_2", Predefined: false},
						{Key: "KEY_1", Value: "value_3", Predefined: false},
					},
				}
			})

			It("writes each variable as a key, value pair on a new line to the buffer", func() {
				allVars := `^KEY_1=value_1\nKEY_2=value_2\nKEY_1=value_3\n$`
				Expect(variableGroup.String()).To(MatchRegexp(allVars))
			})
		})

		Context("when the group doesn't contain any variables", func() {
			BeforeEach(func() {
				variableGroup = VariableGroup{
					Name:      "Empty",
					Variables: []*Variable{},
				}
			})

			It("writes none on a new line to the buffer", func() {
				Expect(variableGroup.String()).To(MatchRegexp(`^none\n$`))
			})
		})
	})
})

var _ = Describe("Variable", func() {
	var variable Variable

	BeforeEach(func() {
		variable = Variable{Key: "MY_VAR", Value: "ITS_VALUE", Predefined: false}
	})

	Describe("String", func() {
		It("writes the variable value on a new line to the buffer", func() {
			Expect(variable.String()).To(MatchRegexp(`^ITS_VALUE\n$`))
		})
	})
})
