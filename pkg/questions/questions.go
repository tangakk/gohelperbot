package questions

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/agnivade/levenshtein"
)

type Question struct {
	Text         string   `json:"text"` //текст вопроса
	lemmatized   string   //лемматизированный вопрос
	Answer       string   `json:"answer"` //ответ
	Extra        string   `json:"extra"`  //доп. вложения
	Subquestions []string `json:"sub"`    //доп. вопросы, если потребуются
}

type Questions struct {
	Q          []Question
	lemmatizer *golem.Lemmatizer
}

func New(l *golem.Lemmatizer) *Questions {
	return &Questions{Q: []Question{}, lemmatizer: l}
}

func (q *Questions) ReadFromJSON(path string) error {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var tmp []Question
	err = json.Unmarshal(jsonData, &tmp)
	if err != nil {
		return err
	}
	for i := 0; i < len(tmp); i++ {
		var re = regexp.MustCompile(`[[:punct:]]`)
		tmp[i].Text = re.ReplaceAllString(tmp[i].Text, "")
		txt := strings.Split(tmp[i].Text, " ")
		for j := 0; j < len(txt); j++ {
			txt[j] = q.lemmatizer.Lemma(txt[j])
		}
		tmp[i].lemmatized = strings.Join(txt, " ")
	}
	q.Q = tmp
	return nil
}

func (q *Questions) Ask(what string) Question {
	//лемматизируем... ну крч приводим к нормальной форме:
	var re = regexp.MustCompile(`[[:punct:]]`)
	what = re.ReplaceAllString(what, "")
	whatSplitted := strings.Split(what, " ")
	for i := 0; i < len(whatSplitted); i++ {
		whatSplitted[i] = q.lemmatizer.Lemma(whatSplitted[i])
	}
	what = strings.Join(whatSplitted, " ")
	//сравниваем каждый вопрос с ответом
	min_w, index := 99999, -1
	for i, v := range q.Q {
		w := levenshtein.ComputeDistance(what, v.lemmatized)
		//fmt.Println(what, v.lemmatized, w)
		min_w = min(min_w, w)
		if min_w == w {
			index = i
		}
	}
	return q.Q[index]
}
