package hw03_frequency_analysis //nolint:golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseText = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

// new case
var equalWords = `Слово: "Слово" - слово, слово. Слово?! Слово, слово!`

// new case
var oneWord = "Word."

// new case
var unicode = `Ā  Ƀȗҥҥȕ - ҭȁҡ ҙɓȁљȕ ćȁӎӱѥ љӱчшӱѥ, ćȁӎӱѥ ԭөѣƥӱѥ ӎәԭɓәԭȕҵӱ
ɓ  ҙөөљөгȗчәćҡөӎ  ćȁԭӱ,  ҡөҭөƥӱѥ  өчәҥь-өчәҥь  љѥѣȗљ  Кƥȕćҭөфәƥ
Рөѣȕҥ.  Ā  өҥȁ  өчәҥь-өчәҥь  љѥѣȕљȁ  әгө. Еә љȗ ҥȁҙɓȁљȕ Ƀȗҥҥȕ ɓ
чәćҭь Пӱхȁ, ȗљȕ Пӱхȁ ҥȁҙɓȁљȕ ɓ әә чәćҭь - ҭәпәƥь ӱжә ҥȕҡҭө  ҥә
ҙҥȁәҭ,  ԭȁжә пȁпȁ Кƥȗćҭөфәƥȁ Рөѣȕҥȁ. Көгԭȁ-ҭө өҥ ҙҥȁљ, ȁ ҭәпәƥь
ҙȁѣыљ.`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		assert.Len(t, Top10(""), 0)
	})

	t.Run("base test", func(t *testing.T) {
		expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
		assert.Subset(t, expected, Top10(baseText))
	})

	// new case
	t.Run("words equals by regexp", func(t *testing.T) {
		expected := []string{"слово"}
		assert.ElementsMatch(t, expected, Top10(equalWords))
	})

	// new case
	t.Run("one word", func(t *testing.T) {
		expected := []string{"word"}
		assert.ElementsMatch(t, expected, Top10(oneWord))
	})

	// new case
	t.Run("any unicode", func(t *testing.T) {
		expected := []string{`ɓ`, `чәćҭь`, `ćȁӎӱѥ`, `пӱхȁ`, `ҥȁҙɓȁљȕ`, `өчәҥь-өчәҥь`, `ā`, `ҭәпәƥь`, `ƀȗҥҥȕ`,
			`кƥȗćҭөфәƥȁ`, `өҥȁ`, `рөѣȕҥȁ`, `көгԭȁ-ҭө`, `ҙҥȁљ`, `ҙɓȁљȕ`, `љӱчшӱѥ`, `ӎәԭɓәԭȕҵӱ`, `ҡөҭөƥӱѥ`,
			`љѥѣȗљ`, `љѥѣȕљȁ`, `әгө`, `ȁ`, `ԭөѣƥӱѥ`, `өҥ`, `ҥә`, `ćȁԭӱ`, `әә`, `ҙөөљөгȗчәćҡөӎ`, `ӱжә`,
			`ҙҥȁәҭ`, `ҙȁѣыљ`, `ҥȕҡҭө`, `ҭȁҡ`, `кƥȕćҭөфәƥ`, `рөѣȕҥ`, `еә`, `љȗ`, `ȗљȕ`, `ԭȁжә`, `пȁпȁ`}
		assert.Subset(t, expected, Top10(unicode))
	})
}
