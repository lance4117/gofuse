package gen

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Article struct {
	Title    string   `json:"title,omitempty" yaml:"title"`
	Author   string   `json:"author,omitempty" yaml:"author"`
	Summary  string   `json:"summary,omitempty" yaml:"summary"`
	Contents []string `json:"contents,omitempty" yaml:"contents"`
}

// NewArticle 生成文章, 输入 段落数量,句子词数
func NewArticle(paragraphs, sentenceWords int) Article {
	// 默认3段
	if paragraphs <= 0 {
		paragraphs = 3
	}
	// 默认一句10个词
	if sentenceWords <= 0 {
		sentenceWords = 10
	}
	title := gofakeit.BookTitle() // 标题感更强
	author := gofakeit.Name()
	summary := OneLiner()
	paras := make([]string, 0, paragraphs)
	for i := 0; i < paragraphs; i++ {
		paras = append(paras, Sentence(sentenceWords))
	}
	return Article{
		Title:    title,
		Author:   author,
		Summary:  summary,
		Contents: paras,
	}
}

// OneLiner 返回一行短文案：优先用 Blurb/Comment 之类，长度可控（大致）
func OneLiner() string {
	// Blurb 更偏营销短句；Comment 偏“短评”
	choices := []string{
		gofakeit.Blurb(),
		gofakeit.Comment(),
	}
	// 去掉可能的换行与多余空白
	out := strings.TrimSpace(strings.ReplaceAll(choices[gofakeit.Number(0, len(choices)-1)], "\n", " "))
	return out
}

// Sentence 按词数粗控，输出英文句子风格文本
func Sentence(words int) string {
	if words < 6 {
		words = 6
	}

	// 备选“名词短语”与“动作/谓词”来源（用领域词替代）
	nounish := []string{
		gofakeit.Company(),   // 公司名
		gofakeit.AppName(),   // 应用名
		gofakeit.BookTitle(), // 书名（名词短语）
		gofakeit.BuzzWord(),  // 商业/技术名词
	}
	verbish := []string{
		strings.ToLower(gofakeit.ConnectiveListing()), // 充当轻谓词/连接
		"enables", "accelerates", "simplifies", "improves",
	}

	caser := cases.Title(language.English)
	title := caser.String(gofakeit.Adverb())
	// 基础片段：副词 + 形容词 + 名词短语 + 连接词 + 动作 + 形容词 + 名词短语
	parts := []string{
		title,                // e.g. Quickly
		gofakeit.Adjective(), // agile
		nounish[gofakeit.Number(0, len(nounish)-1)], // Company/App/Title/Word
		gofakeit.ConnectiveComparative(),            // versus/than/…
		verbish[gofakeit.Number(0, len(verbish)-1)], // enables/accelerates/…
		gofakeit.Adjective(),                        // scalable
		nounish[gofakeit.Number(0, len(nounish)-1)],
	}

	// 根据目标 words 扩写：追加（连接词 + 形容词 + 名词替代）
	toks := strings.Fields(strings.Join(parts, " "))
	for len(toks) < words {
		toks = append(toks, gofakeit.ConnectiveCasual(), gofakeit.Adjective(), nounish[gofakeit.Number(0, len(nounish)-1)])
	}
	if len(toks) > words {
		toks = toks[:words]
	}

	line := strings.Join(toks, " ")
	if !strings.HasSuffix(line, ".") {
		line += "."
	}
	return line
}

// Paragraph 用合成句拼装段落
func Paragraph(sentences, wordsPerSentence int) string {
	if sentences <= 0 {
		sentences = 3
	}
	if wordsPerSentence <= 0 {
		wordsPerSentence = 12
	}
	buf := make([]string, 0, sentences)
	for i := 0; i < sentences; i++ {
		buf = append(buf, Sentence(wordsPerSentence))
	}
	return strings.Join(buf, " ")
}

// Letters 生成多少个随机字母
func Letters(n uint) string {
	return gofakeit.LetterN(n)
}

// LettersAndNumbers 生成指定长度的随机字符串（数字 + 大小写字母）
// 使用 crypto/rand 保证安全级随机。
func LettersAndNumbers(n int) string {
	if n <= 0 {
		return ""
	}
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size := big.NewInt(int64(len(chars)))

	result := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, size)
		if err != nil {
			result[i] = chars[i%len(chars)] // fallback
			continue
		}
		result[i] = chars[num.Int64()]
	}
	return string(result)
}
