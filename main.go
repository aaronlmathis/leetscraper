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
    date := problemOfTheDay.Data.ActiveDailyCodingChallengeQuestion.Date
    titleSlug := problemOfTheDay.Data.ActiveDailyCodingChallengeQuestion.Question.TitleSlug
    url := "https://leetcode.com" + problemOfTheDay.Data.ActiveDailyCodingChallengeQuestion.Link
    fmt.Printf(`Success! The Problem of the day for %s is at %s`, date, url)

    problemDetails, err := scraper.GetProblemDetails(titleSlug)
    content := problemDetails.Data.Question.Content

    if err != nil {
        log.Printf("Error getting problem details: %v", err)
    }
    fmt.Printf("%s\n", content)

    
}

