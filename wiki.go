package main

import (
	"regexp"
	"strings"
)

var (
	spaces = regexp.MustCompile(`[\s|\|]+`)

	categoryLink = regexp.MustCompile(`\[\[Category:.*?\]\]`)
	internalLink = regexp.MustCompile(`\[\[(.*?)\]\]`)

	cleanRegexps = []*regexp.Regexp{
		regexp.MustCompile(`\[(.*?)\]`),
		regexp.MustCompile("(?s:{{.*?}})"),
		regexp.MustCompile(`(?s:\<!--.*--\>)`),
		regexp.MustCompile(`(?s:\<.*?\>.*?\<.*?\>)`),
		regexp.MustCompile(`(?s:\<.*?\>)`)} //,

	nonWord = regexp.MustCompile(`[^0-9A-Za-z_\s]`)

	newLines = regexp.MustCompile(`\n+`)
)

// Here is an example article from the Wikipedia XML dump
//
// <page>
// 	<title>Apollo 11</title>
//      <redirect title="Foo bar" />
// 	...
// 	<revision>
// 	...
// 	  <text xml:space="preserve">
// 	  {{Infobox Space mission
// 	  |mission_name=&lt;!--See above--&gt;
// 	  |insignia=Apollo_11_insignia.png
// 	...
// 	  </text>
// 	</revision>
// </page>
//
// Note how the tags on the fields of Page and Redirect below
// describe the XML schema structure.

type Page struct {
	Title string `xml:"title"`
	//Redir Redirect `xml:"redirect"`
	Text string `xml:"revision>text"`

	PlainText string
}

func ToPlainText(wikitext string) string {
	if strings.HasPrefix(wikitext, "#REDIRECT") {
		return ""
	}
	return strings.Replace(strings.Replace(wikitext, "=", "", -1), ",", "", -1)
	wikitext = string(categoryLink.ReplaceAll([]byte(wikitext), []byte("")))
	//wikitext = string(internalLink.ReplaceAll([]byte(wikitext), []byte("$1")))

	for _, regex := range cleanRegexps {
		wikitext = string(regex.ReplaceAll([]byte(wikitext), []byte("")))
	}
	wikitext = string(newLines.ReplaceAll([]byte(wikitext), []byte("\n")))
	return wikitext
}

func PlainTextToWords(plainText string) []string {
	return strings.Split(plainText, " ")
}
