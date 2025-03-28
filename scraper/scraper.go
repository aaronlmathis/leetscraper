package scraper

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
)

type GraphQLPayload struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables,omitempty"`
}

type DailyChallengeResponse struct {
    Data struct {
        ActiveDailyCodingChallengeQuestion struct {
            Date    string `json:"date"`
            Link    string `json:"link"`
            Question struct {
                TitleSlug string `json:"titleSlug"`
            } `json:"question"`
        } `json:"activeDailyCodingChallengeQuestion"`
    } `json:"data"`
}

type ProblemDetailsResponse struct {
    Data struct {
        Question struct {
            QuestionId string `json:"questionFrontendId"`
            Content string  `json:"content"`
            CodeSnippets []struct {
                Lang        string  `json:"lang"`
                LangSlug    string  `json:"langSlug"`
                Code        string  `json:"code"`
            } `json:"codeSnippets"`
            Title       string  `json:"title"`
            Difficulty  string  `json:"difficulty"`
        } `json:"question"`
    } `json:"data"`
}

type DailyCodingChallengeQuestion struct {
	Date string
	Link string
	TitleSlug string
}

type LeetCodeProblem struct {
	Id string
	Title string
	Description string
	Difficulty string
	CodeSnippets []struct {
        Lang		string `json:"lang"`
        LangSlug	string `json:"langSlug"`
        Code		string `json:"code"`
	}
} 

func GetProblemOfTheDay() (DailyCodingChallengeQuestion, error) {
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
        log.Printf("Error marshalling payload: %s", err)
        return DailyCodingChallengeQuestion{}, err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
    if err != nil {
        log.Printf("Error creating request: %s", err)
        return DailyCodingChallengeQuestion{}, err
    }

    req.Header.Set("Content-Type", "application/json")

    // Make the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %s", err)
        return DailyCodingChallengeQuestion{}, err
    }
    defer res.Body.Close()

    // Decode the response
    var result DailyChallengeResponse
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        log.Printf("Error decoding JSON: %s", err)
        return DailyCodingChallengeQuestion{}, err
    }
	
	dailyQuestion := DailyCodingChallengeQuestion {
		Date: result.Data.ActiveDailyCodingChallengeQuestion.Date,
		Link: result.Data.ActiveDailyCodingChallengeQuestion.Link,
		TitleSlug: result.Data.ActiveDailyCodingChallengeQuestion.Question.TitleSlug,
	}
		
	
    return dailyQuestion, nil
}

func GetProblemDetails(titleSlug string) (LeetCodeProblem, error) {
    payload := GraphQLPayload {
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
        log.Printf("Error marshalling payload: %s", err)
        return LeetCodeProblem{}, err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
    if err != nil {
        log.Printf("Error creating request: %s", err)
        return LeetCodeProblem{}, err
    }

    req.Header.Set("Content-Type", "application/json")

    // Make the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %s", err)
        return LeetCodeProblem{}, err
    }
    defer res.Body.Close()

    // Decode the response
    var result ProblemDetailsResponse
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        log.Printf("Error decoding JSON: %s", err)
        return LeetCodeProblem{}, err
    }
	
    markdownContent, err := ConvertHTMLToMarkdown(result.Data.Question.Content)
    if err != nil {
        log.Printf("Error converting html to markdown: %s", err)
        return LeetCodeProblem{}, err
    }

	problem := LeetCodeProblem{
		Id: result.Data.Question.QuestionId,
		Title: result.Data.Question.Title,
        Description: markdownContent,
        Difficulty: result.Data.Question.Difficulty,
        CodeSnippets: result.Data.Question.CodeSnippets,
    }

    return problem, nil
}

