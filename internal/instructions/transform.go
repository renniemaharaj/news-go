package instructions

import (
	"fmt"
	"strings"
)

func BuildTransformPrompt(textContent string, images []string) (string, error) {
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
		- Also select images from the list that best represent the article.
	
	Following article:
	%s
	
	Images to choose from:
	%s
	`, textContent, strings.Join(images, ", "))

	final := inject(base, map[string]string{
		"ruleGoesHere":   transformRules,
		"promptGoesHere": combinedPrompt,
	})

	return final, nil
}
