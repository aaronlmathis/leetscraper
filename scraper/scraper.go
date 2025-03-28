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

func GetProblemOfTheDay() (DailyChallengeResponse, error) {
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
        return DailyChallengeResponse{}, err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
    if err != nil {
        log.Printf("Error creating request: %s", err)
        return DailyChallengeResponse{}, err
    }

    req.Header.Set("Content-Type", "application/json")

    // Make the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %s", err)
        return DailyChallengeResponse{}, err
    }
    defer res.Body.Close()

    // Decode the response
    var result DailyChallengeResponse
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        log.Printf("Error decoding JSON: %s", err)
        return DailyChallengeResponse{}, err
    }

    return result, nil
}

type ProblemDetailsResponse struct {
    Data struct {
        Question struct {
            Content string  `json:"content"`
            CodeSnippets []struct {
                Lang        string  `json:"lang"`
                LangSlug    string  `json:"langSlug"`
                Code        string  `json:"code"`
            } `json:"codeSnippets"`
            Title       string  `json:"Title"`
            Difficulty  string  `json:"difficulty"`
        } `json:"question"`
    } `json:"data"`
}

func GetProblemDetails(titleSlug string) (ProblemDetailsResponse, error) {
    payload := GraphQLPayload {
        Query: `query questionData($titleSlug: String!) {
            question(titleSlug: $titleSlug) {
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
    }
    
    // Marshal the payload into JSON
    jsonBody, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling payload: %s", err)
        return ProblemDetailsResponse{}, err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
    if err != nil {
        log.Printf("Error creating request: %s", err)
        return ProblemDetailsResponse{}, err
    }

    req.Header.Set("Content-Type", "application/json")

    // Make the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %s", err)
        return ProblemDetailsResponse{}, err
    }
    defer res.Body.Close()

    // Decode the response
    var result ProblemDetailsResponse
    err = json.NewDecoder(res.Body).Decode(&result)
    if err != nil {
        log.Printf("Error decoding JSON: %s", err)
        return ProblemDetailsResponse{}, err
    }

    return result, nil
}

