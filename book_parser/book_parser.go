package book_parser

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/html"
	"knigavuhe/dom"
	"regexp"
	"strings"
)

type BookInfo struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Series   string `json:"series"`
	Genre    string `json:"genre"`
	Reader   string `json:"reader"`
	CoverUrl string `json:"cover_url"`
}

func ExtractBookInfo(htmlCode string) (result *BookInfo) {
	doc, err := html.Parse(strings.NewReader(htmlCode))
	if err != nil {
		return nil
	}

	result = &BookInfo{}

	if div := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("div"), dom.HasClass("book_genre_pretitle")}); div != nil {
		if a := dom.FindFirst(div, []func(*html.Node) bool{dom.IsTag("a")}); a != nil {
			if found, content := dom.GetContent(a); found {
				result.Genre = content
			}
		}
	}

	if div := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("div"), dom.HasClass("book_cover")}); div != nil {
		if img := dom.FindFirst(div, []func(*html.Node) bool{dom.IsTag("img")}); img != nil {
			if found, content := dom.GetAttrValue(img, "src"); found {
				result.CoverUrl = content
			}
		}
	}

	if span := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("span"), dom.HasClass("book_title_name")}); span != nil {
		if found, content := dom.GetContent(span); found {
			result.Title = content
		}
	}

	if span := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("span"), dom.HasAttrWithValue("itemprop", "author")}); span != nil {
		if a := dom.FindFirst(span, []func(*html.Node) bool{dom.IsTag("a")}); a != nil {
			if found, content := dom.GetContent(a); found {
				result.Author = content
			}
		}
	}

	if span := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("span"), dom.HasClass("book_title_elem"), dom.HasImmediateChild("a")}); span != nil {
		if a := dom.FindFirst(span, []func(*html.Node) bool{dom.IsTag("a")}); a != nil {
			if found, content := dom.GetContent(a); found {
				result.Reader = content
			}
		}
	}

	if div := dom.FindFirst(doc, []func(*html.Node) bool{dom.IsTag("div"), dom.HasClass("book_serie_block_title")}); div != nil {
		if a := dom.FindFirst(div, []func(*html.Node) bool{dom.IsTag("a")}); a != nil {
			if found, content := dom.GetContent(a); found {
				result.Series = content
			}
		}
	}

	return
}

type BookChapter struct {
	Duration      int     `json:"duration"`
	DurationFloat float64 `json:"duration_float"`
	Error         int     `json:"error"`
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	Url           string  `json:"url"`
}

func ExtractBookChapters(htmlCode string) (error, []BookChapter) {
	re := regexp.MustCompile(`var player = new BookPlayer\((.*)\);`)
	m := re.FindStringSubmatch(htmlCode)
	if m == nil {
		return errors.New("failed to find player data in the body"), nil
	}

	jsonStr := []byte("[" + string(m[1]) + "]")

	var data []json.RawMessage

	if err := json.Unmarshal(jsonStr, &data); err != nil {
		return err, nil
	}

	var result []BookChapter

	if err := json.Unmarshal(data[1], &result); err != nil {
		return err, nil
	}

	return nil, result
}
