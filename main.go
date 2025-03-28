package main

import (
    "fmt"
    "log"
    "github.com/aaronlmathis/leetscraper/scraper"
)

func main() {
    problemOfTheDay, err := scraper.GetProblemOfTheDay()
    if err != nil {
        log.Printf("Error scraping problem of the day: %v", err)
    }
    date := problemOfTheDay.Date
    titleSlug := problemOfTheDay.TitleSlug
    url := "https://leetcode.com" + problemOfTheDay.Link
    fmt.Printf("The Problem of the day for %v is at %v\n\n\n", date, url)

    leetCodeProblem, err := scraper.GetProblemDetails(titleSlug)
    if err != nil {
        log.Printf("Error getting problem details: %v", err)
    }
    fmt.Printf("#%v - %v\n\n", leetCodeProblem.Id, leetCodeProblem.Title)
    fmt.Printf("%v \n\n", leetCodeProblem.Description)
    
    for _, snippet := range leetCodeProblem.CodeSnippets {
        if snippet.LangSlug == "python3" || snippet.LangSlug == "golang" {
            fmt.Printf("%v\n", snippet.Code)
        }
    }    
}

