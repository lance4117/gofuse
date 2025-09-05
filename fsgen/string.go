package fsgen

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

// NewArticle 生成文章
// paragraphs: 段落数
// sentencesPerParagraph: 每段句子数
// wordsPerSentence: 每句词数
func NewArticle(paragraphs, sentencesPerParagraph, wordsPerSentence int) string {
	err := gofakeit.Seed(0)
	if err != nil {
		return ""
	}

	var article strings.Builder

	for p := 0; p < paragraphs; p++ {
		for s := 0; s < sentencesPerParagraph; s++ {
			article.WriteString(gofakeit.Sentence(wordsPerSentence) + " ")
		}
		article.WriteString("\n") // 段落间换行
	}

	return article.String()
}
