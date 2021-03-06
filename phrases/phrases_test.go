package phrases

import (
	"reflect"
	"testing"

	"github.com/gelembjuk/keyphrases/sentences"
	"github.com/gelembjuk/keyphrases/words"
)

func TestGetPhrases(t *testing.T) {
	SetLanguage("english")
	sentences.SetLanguage("english")
	words.SetLanguage("english")

	text := `Whenever Mario Draghi clears a hurdle on his path to higher inflation, a new one appears, euro-area economy.
Just as the 19-nation economy sends encouraging signals that challenges from Brexit to terrorism won’t derail the modest recovery, a new decline in oil prices is casting a shadow over an expected pick-up in inflation. With growth not strong enough to generate price pressures, the European Central Bank president may have to revise his outlook yet again.
Inflation remains far below the ECB’s 2 percent goal after more than two years of unprecedented stimulus and isn’t seen reaching it before 2018. Staff will begin to draw up fresh forecasts in mid-August, and while officials are in no rush to adjust or expand their 1.7 trillion-euro ($1.9 trillion) quantitative-easing plan in September, economists predict Draghi will have to ease policy before the end of the year.
“Now that the euro-area economy seems to have shrugged off the Brexit vote, Mario Draghi focus will again shift on inflation, euro-area economy, before 2018, against the background of those negative news from oil prices,” said Johannes Gareis, an economist at Natixis in Frankfurt. “Yes, the ECB has managed to dispel deflation fears, but all the uncertainty means inflation will stay lower for longer -- and Draghi will have to take notice.”`

	sentenceslist, _ := sentences.SplitTextForSentencesFromNews(text)

	wordslist, _ := words.SplitSentencesForWords(sentenceslist)

	phrases, _ := getBasicPhrasesHash(sentenceslist, wordslist)

	expected := map[string]int{
		"oil prices":        2,
		"euro area economy": 3,
		"Mario Draghi":      2,
	}

	if !reflect.DeepEqual(phrases, expected) {
		t.Fatalf("Got %s, expected %s.", phrases, expected)
	}

}
