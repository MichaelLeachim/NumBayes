package bayes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
  "strings"
	"regexp"
	"sort"
)

type BayesMemory struct {
	CategoryDocument map[string]int
	CategoryWord     map[string]map[string]int
}

func Diff(data []string, ignore string) []string {
	kex := regexp.MustCompile(fmt.Sprintf("[0-9%s]+", ignore))
	result := make([]string, 0)
	for _, v := range data {
		for _, positional := range kex.FindAllString(v, -1) {
			result = append(result, NumberPositionalBinning(positional)...)
		}
		result = append(result, v)
	}
	return result
}

func (b *BayesMemory) checkNil() {
	if b.CategoryDocument == nil {
		b.CategoryDocument = make(map[string]int, 0)
	}
	if b.CategoryWord == nil {
		b.CategoryWord = make(map[string]map[string]int, 0)
	}
}

func (b *BayesMemory) TokenizeSimple(data string,keep string)[]string{
  return Diff(strings.Split(data," "),keep)
}

func (b *BayesMemory) Train(words, cats []string) {
	b.checkNil()
	for _, category := range cats {
		b.CategoryDocument[category] += 1
		for _, word := range words {
			if b.CategoryWord[category] == nil {
				b.CategoryWord[category] = make(map[string]int, 0)
			}

			b.CategoryWord[category][word] += 1
		}
	}
}

func (b *BayesMemory) UnTrain(words, cats []string) {
	b.checkNil()
	for _, category := range cats {
		b.CategoryDocument[category] -= 1
		if b.CategoryDocument[category] < 0 {
			b.CategoryDocument[category] = 0
		}
		for _, word := range words {
			if b.CategoryWord[category] == nil {
				b.CategoryWord[category] = make(map[string]int, 0)
			}

			b.CategoryWord[category][word] -= 1
			if b.CategoryWord[category][word] < 0 {
				b.CategoryWord[category][word] = 0
			}
		}
	}
}

func (b *BayesMemory) Marshal() ([]byte, error) {
	return AnyToByte(b)
}

func (b *BayesMemory) UnMarshal(data []byte) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(b)
}

func (b *BayesMemory) documents_at_all() int {
	sum := 0
	for _, v := range b.CategoryDocument {
		sum += v
	}
	return sum
}

func (b *BayesMemory) words_in_category_at_all(cat string) int {
	sum_cat := 0
	for _, v := range b.CategoryWord[cat] {
		sum_cat += v
	}
	return sum_cat
}

type ClassifyResult struct {
	Log      float64
	Simple   float64
	Category string
}

func (b *BayesMemory) Classify(words []string) ClassifyResultSlice {
	result := ClassifyResultSlice{}
	all_documents := b.documents_at_all()
	// in case classifier is empty
	if all_documents <= 0 {
		return result
	}
	for cat, count := range b.CategoryDocument {
		// category with no documents cannot be used/selected
		if count <= 0 {
			continue
		}
		// category_prob = this-category-has-documents/all-documents
		category_prob := float64(count) / float64(all_documents)

		simple_bayes := category_prob
		log_bayes := math.Log(category_prob)

		words_in_category := b.words_in_category_at_all(cat)
		// category with no features is skipped
		if words_in_category <= 0 {
			continue
		}
		for _, word := range words {
			var categoryWord float64
			// P = this-word-in-category/all-words-in-category
			if b.CategoryWord[cat][word] == 0 {
				// 0 will render other probabilities  meaningless
				categoryWord = float64(0.0000000001)
			} else {
				categoryWord = float64(b.CategoryWord[cat][word])
			}
			word_prob := categoryWord / float64(words_in_category)
			simple_bayes *= word_prob
			log_bayes += math.Log(word_prob)
		}
		result = append(result, ClassifyResult{Log: log_bayes, Simple: simple_bayes, Category: cat})
	}
	sort.Sort(result)
	return result
}
