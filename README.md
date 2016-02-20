# Numerical bayes is an implementation of NB classifier that takes into account continuous data.
Useful when you are dealing with data that is difficult to use before calculat
# Feutures:
* Reasonbly, well tested
* Build for development. Serialization is simple.
** Storage is up to developer.
** Tokenization is up to developer. Better use Snowball stemming, or smth. like this.

## Try to use some numeric data. For example, ages.

```golang
b := BayesMemory{}
splitLocal := func(data string) []string {
	return b.TokenizeSimple(data, "")
}
atestLocal := func(data, result string) {
	assert.Equal(t, b.Classify(splitLocal(data))[0].Category, result)
}
  // Train some age
b.Train(splitLocal("Hi. I've been born in 1945"), []string{"old"})
b.Train(splitLocal("I've been born in 1980"), []string{"middle"})
b.Train(splitLocal("1990 "), []string{"young"})
b.Train(splitLocal("23"), []string{"young"})
b.Train(splitLocal("60"), []string{"quite old"})
b.Train(splitLocal("80"), []string{"old"})
b.Train(splitLocal("100"), []string{"very old"})
b.Train(splitLocal("1940"), []string{"old"})
b.Train(splitLocal("1960"), []string{"old"})
  // Test some age
atestLocal("1948", "old")
atestLocal("1985", "middle")
atestLocal("1996", "young")
atestLocal("67", "quite old")
atestLocal("89", "old")
atestLocal("1967", "old")
```

## Try to test some duration data
```golang
b := BayesMemory{}
splitLocal := func(data string) []string {
	return b.TokenizeSimple(data, ":")
}
atestLocal := func(data, result string) {
	assert.Equal(t, b.Classify(splitLocal(data))[0].Category, result)
}
 
b.Train(splitLocal("An audio  12:23:34"), []string{"very long"})
b.Train(splitLocal("An audio 03:43:45"), []string{"very long"})
b.Train(splitLocal("00:23:56 some data"), []string{"short"})
b.Train(splitLocal("it is duration 00:45:53"), []string{"quite long"})
b.Train(splitLocal("00:03:12 but it works"), []string{"short"})
b.Train(splitLocal("00:08:12"), []string{"short"})
b.Train(splitLocal("00:12:12"), []string{"short"})
b.Train(splitLocal("02:14:14"), []string{"very long"})
b.Train(splitLocal("05:15:67"), []string{"very long"})
// ========================
atestLocal("00:05:04", "short")
atestLocal("00:25:04", "short")
atestLocal("00:44:23", "quite long")
atestLocal("03:45:20", "very long")
atestLocal("15:05:50", "very long")
atestLocal("12:31:56", "very long")
atestLocal("10:31:56", "very long")
atestLocal("05:31:56", "very long")
// TODO: add geospatial data
```

## Now, serialize it(and, maybe, save into Bolt.DB)
```
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
```

## License
MIT License

