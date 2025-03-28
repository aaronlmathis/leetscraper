package scraper

import (
    "log"
    htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func ConvertHTMLToMarkdown(htmlContent string) (string, error) {
    markdown, err := htmltomarkdown.ConvertString(htmlContent)
    if err != nil {
        log.Printf("Failed to convert HTML to Markdown: %v", err)
        return htmlContent, err
    }

    return markdown, nil
}
