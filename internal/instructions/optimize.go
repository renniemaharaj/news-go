package instructions

import (
	"fmt"
	"strings"
)

func BuildOptimizationPrompt(masterList, userList []string, userPrompt string) (string, error) {
	base, err := load("base.txt")
	if err != nil {
		return "", err
	}

	optimizeRules, err := load("optimize.txt")
	if err != nil {
		return "", err
	}

	combinedPrompt := fmt.Sprintf(`
	Optimize and return the two an optimized list object:

	Master List:
	%s

	User Preference:
	%s

	User Prompt:
	%s
`, strings.Join(masterList, ", "), strings.Join(userList, ", "), userPrompt)

	final := inject(base, map[string]string{
		"ruleGoesHere":   optimizeRules,
		"promptGoesHere": combinedPrompt,
	})

	return final, nil
}
