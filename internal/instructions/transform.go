package instructions

import "fmt"

func BuildTransformPrompt(textContent string) (string, error) {
	base, err := load("base.txt")
	if err != nil {
		return "", err
	}

	transformRules, err := load("transform.txt")
	if err != nil {
		return "", err
	}

	combinedPrompt := fmt.Sprintf(`
	Based on the following article, populate and return a report object:
	
	Following article:
	%s`, textContent)

	final := inject(base, map[string]string{
		"ruleGoesHere":   transformRules,
		"promptGoesHere": combinedPrompt,
	})

	return final, nil
}
