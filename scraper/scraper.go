package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type DailyChallengeResponse struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Date     string `json:"date"`
			Link     string `json:"link"`
			Question struct {
				TitleSlug string `json:"titleSlug"`
			} `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}
type ProblemDetailsResponse struct {
	Data struct {
		Question struct {
			QuestionId   string        `json:"questionFrontendId"`
			Content      string        `json:"content"`
			CodeSnippets []CodeSnippet `json:"codeSnippets"`
			Title        string        `json:"title"`
			Difficulty   string        `json:"difficulty"`
		} `json:"question"`
	} `json:"data"`
}

type DailyCodingChallengeQuestion struct {
	Date      string
	Link      string
	TitleSlug string
}

type CodeSnippet struct {
	Lang     string
	LangSlug string
	Code     string
}

type LeetCodeProblem struct {
	Id           string
	Title        string
	TitleSlug    string
	Description  string
	Difficulty   string
	CodeSnippets []CodeSnippet
}

func toSet(slugs []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slugs))
	for _, slug := range slugs {
		set[slug] = struct{}{}
	}
	return set
}

func GetDailyCodingChallengeQuestion() (DailyCodingChallengeQuestion, error) {
	payload := GraphQLPayload{
		Query: `query questionOfToday {
                    activeDailyCodingChallengeQuestion {
                        date
                        link
                        question {
                            titleSlug
                        }
                    }
                }`,
	}

	// Marshal the payload into JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return DailyCodingChallengeQuestion{}, fmt.Errorf("error marshalling payload: %s", err)

	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return DailyCodingChallengeQuestion{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return DailyCodingChallengeQuestion{}, fmt.Errorf("error sending request: %s", err)
	}
	defer res.Body.Close()

	// Decode the response
	var result DailyChallengeResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return DailyCodingChallengeQuestion{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	// Populate dailyQuestion
	dailyQuestion := DailyCodingChallengeQuestion{
		Date:      result.Data.ActiveDailyCodingChallengeQuestion.Date,
		Link:      result.Data.ActiveDailyCodingChallengeQuestion.Link,
		TitleSlug: result.Data.ActiveDailyCodingChallengeQuestion.Question.TitleSlug,
	}

	return dailyQuestion, nil
}

func GetProblemDetails(titleSlug string, langSlugs []string) (LeetCodeProblem, error) {
	payload := GraphQLPayload{
		Query: `query questionData($titleSlug: String!) {
            question(titleSlug: $titleSlug) {
                questionFrontendId
                content
                codeSnippets {
                    lang
                    langSlug
                    code
                }
                title
                difficulty
            }
        }`,
		Variables: map[string]interface{}{
			"titleSlug": titleSlug,
		},
	}

	// Marshal the payload into JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {

		return LeetCodeProblem{}, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return LeetCodeProblem{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return LeetCodeProblem{}, fmt.Errorf("error sending request: %s", err)
	}
	defer res.Body.Close()

	// Decode the response
	var result ProblemDetailsResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return LeetCodeProblem{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	// Convert HTML content to markdown for description
	markdownContent, err := ConvertHTMLToMarkdown(result.Data.Question.Content)
	if err != nil {
		return LeetCodeProblem{}, fmt.Errorf("error converting html to markdown: %s", err)
	}

	// Filter only preferred languages

	var preferredLangSnippets []CodeSnippet
	langSlugSet := toSet(langSlugs)

	for _, snippet := range result.Data.Question.CodeSnippets {
		if _, ok := langSlugSet[snippet.LangSlug]; ok {
			preferredLangSnippets = append(preferredLangSnippets, snippet)
		}
	}

	// Build problem
	problem := LeetCodeProblem{
		Id:           result.Data.Question.QuestionId,
		Title:        result.Data.Question.Title,
		TitleSlug:    titleSlug,
		Description:  markdownContent,
		Difficulty:   result.Data.Question.Difficulty,
		CodeSnippets: preferredLangSnippets,
	}

	return problem, nil
}
