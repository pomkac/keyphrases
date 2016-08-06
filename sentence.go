package keyphrases

import (
	"regexp"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

func (obj *TextPhrases) splitTextForSentences(text string) ([]string, error) {
	// prepare tokenizer
	sentenceslist := []string{}

	langfile := "data/" + obj.Language + ".json"

	b, err := data.Asset(langfile)

	if err != nil {
		return sentenceslist, err
	}

	// load the training data
	training, err := sentences.LoadTraining(b)

	if err != nil {
		return sentenceslist, err
	}

	// create the default sentence tokenizer
	tokenizer := sentences.NewSentenceTokenizer(training)

	text, _ = obj.cleanTextAfterHTML(text)

	if obj.NewsSource {
		// this text is from news sources. It can have specific news format
		// clean a text from standard news message formatting , and specific language
		text, _, _ = obj.langobj.CleanNewsMessage(text)
	}

	sentencesobjs := tokenizer.Tokenize(text)

	for _, s := range sentencesobjs {
		sentence := s.Text

		// remove last symbol of a sentence if it is a dot or so
		if len(sentence) < 3 {
			continue
		}

		sentence, _ = obj.cleanAndNormaliseSentence(sentence)

		sentenceslist = append(sentenceslist, sentence)
	}

	return sentenceslist, nil
}

func (obj *TextPhrases) cleanAndNormaliseSentence(sentence string) (string, error) {

	sentence, _ = obj.langobj.CleanAndNormaliseSentence(sentence)

	replace := [][]string{
		{"[\\[\\]}{]", ""},
		{"[:;-]", " "},
		{"[.?!):]", " "},
		{"\\s\\s+", " "},
		{"^\\s+", ""},
		{"\\s+$", ""},
	}

	for _, template := range replace {
		r := regexp.MustCompile(template[0])

		sentence = r.ReplaceAllString(sentence, template[1])
	}

	return sentence, nil
}
