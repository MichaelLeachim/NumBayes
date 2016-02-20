package bayes

import (
	"github.com/kljensen/snowball"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func splitmin(data string) []string {
	rex := regexp.MustCompile("\\w+")
	d := rex.FindAllString(data, -1)
	for k, v := range d {
		stemmed, err := snowball.Stem(v, "english", true)
		if err == nil {
			d[k] = strings.ToLower(stemmed)
		} else {
			d[k] = strings.ToLower(v)
		}
	}
	return d
}
func atest(t *testing.T, b BayesMemory, data, result string) {
	assert.Equal(t, b.Classify(Diff(splitmin(data), ""))[0].Category, result)
}

func TestSerialization(t *testing.T) {
	b := BayesMemory{}

	b.Train(splitmin("Dogs are awesome, cats too. I love my dog"), []string{"dog"})
	b.Train(splitmin("Cats are more preferred by software developers. I never could stand cats. I have a dog"), []string{"cat"})
	b.Train(splitmin("My dog's name is Willy. He likes to play with my wife's cat all day long. I love dogs"), []string{"dog"})
	b.Train(splitmin("Cats are difficult animals, unlike dogs, really annoying, I hate them all"), []string{"cat"})
	b.Train(splitmin("So which one should you choose? A dog, definitely."), []string{"dog"})
	b.Train(splitmin("The favorite food for cats is bird meat, although mice are good, but birds are a delicacy"), []string{"cat"})
	b.Train(splitmin("A dog will eat anything, including birds or whatever meat"), []string{"dog"})
	b.Train(splitmin("My cat's favorite place to purr is on my keyboard"), []string{"cat"})
	b.Train(splitmin("My dog's favorite place to take a leak is the tree in front of our house"), []string{"dog"})

	blabo, err := b.Marshal()
	if err != nil {
		t.Error(err)
	}
	c := BayesMemory{}
	err = c.UnMarshal(blabo)
	if err != nil {
		t.Error(err)
	}
	if len(c.CategoryWord) == 0 {
		t.Error(err)
	}
	// test serialization

	atest(t, c, "This test is about cats", "cat")
	atest(t, c, "This test is about cats", "cat")
	atest(t, c, "I hate ...", "cat")
	atest(t, c, "The most annoying animal on earth.", "cat")
	atest(t, c, "The preferred company of software developers.", "cat")
	atest(t, c, "My precious, my favorite!", "dog")
	atest(t, c, "Get off my keyboard!", "cat")
	atest(t, c, "Kill that bird!", "cat")

	atest(t, c, "This test is about dogs.", "dog")
	atest(t, c, "Cats or Dogs?", "dog")
	atest(t, c, "What pet will I love more?", "dog")
	atest(t, c, "Willy, where the heck are you?", "dog")
	atest(t, c, "I like big buts and I cannot lie.", "cat")
	atest(t, c, "Why is the front door of our house open?", "dog")
	atest(t, c, "Who is eating my meat?", "dog")
}

func TestMemoryBayes(t *testing.T) {
	b := BayesMemory{}
	b.Train(splitmin("Dogs are awesome, cats too. I love my dog"), []string{"dog"})
	b.Train(splitmin("Cats are more preferred by software developers. I never could stand cats. I have a dog"), []string{"cat"})
	b.Train(splitmin("My dog's name is Willy. He likes to play with my wife's cat all day long. I love dogs"), []string{"dog"})
	b.Train(splitmin("Cats are difficult animals, unlike dogs, really annoying, I hate them all"), []string{"cat"})
	b.Train(splitmin("So which one should you choose? A dog, definitely."), []string{"dog"})
	b.Train(splitmin("The favorite food for cats is bird meat, although mice are good, but birds are a delicacy"), []string{"cat"})
	b.Train(splitmin("A dog will eat anything, including birds or whatever meat"), []string{"dog"})
	b.Train(splitmin("My cat's favorite place to purr is on my keyboard"), []string{"cat"})
	b.Train(splitmin("My dog's favorite place to take a leak is the tree in front of our house"), []string{"dog"})

	// test serialization

	atest(t, b, "This test is about cats", "cat")
	atest(t, b, "This test is about cats", "cat")
	atest(t, b, "I hate ...", "cat")
	atest(t, b, "The most annoying animal on earth.", "cat")
	atest(t, b, "The preferred company of software developers.", "cat")
	atest(t, b, "My precious, my favorite!", "dog")
	atest(t, b, "Get off my keyboard!", "cat")
	atest(t, b, "Kill that bird!", "cat")

	atest(t, b, "This test is about dogs.", "dog")
	atest(t, b, "Cats or Dogs?", "dog")
	atest(t, b, "What pet will I love more?", "dog")
	atest(t, b, "Willy, where the heck are you?", "dog")
	atest(t, b, "I like big buts and I cannot lie.", "cat")
	atest(t, b, "Why is the front door of our house open?", "dog")
	atest(t, b, "Who is eating my meat?", "dog")
}
