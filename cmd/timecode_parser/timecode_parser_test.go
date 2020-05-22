package timecodeparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("when raw text contains more than 3 timecodes", func(t *testing.T) {
		rawText := `
			00:11 Human Blue - Electric Harmonie
			05:13 Chris Liberator - The Cult
			09:30 David Moleon - Mole On
			14:15 Eric Sneo - Pulses
			17:35 Guy McAffer & Rackitt - B1 Untitled [RAW 020]
			20:57 The Advent - Recreate 6 (Player Remix 1)
			Long prefix before timecode 01:30:59 Unknown
		`
		collection := Parse(rawText)

		assert.Equal(t, 7, len(collection))
		assert.Equal(t, 11, collection[0].Seconds)
		assert.Equal(t, "Human Blue - Electric Harmonie", collection[0].Description)
		assert.Equal(t, 5459, collection[6].Seconds)
		assert.Equal(t, "Long prefix before timecode", collection[6].Description)
	})

	t.Run("when raw text contains less than 3 timecodes", func(t *testing.T) {
		rawText := `
			00:11 Human Blue - Electric Harmonie
			05:13 Chris Liberator - The Cult
		`
		collection := Parse(rawText)

		assert.Equal(t, 0, len(collection))
	})
}

func Test_fetchDescription(t *testing.T) {
	t.Run("when correct description given", func(t *testing.T) {
		texts := []string{
			"",
			"- Human Blue - Electric Harmonie",
		}

		description := fetchDescription(texts)

		assert.Equal(t, "Human Blue - Electric Harmonie", description)
	})

	t.Run("when description contains special characters", func(t *testing.T) {
		texts := []string{
			"Prefix string",
			"- 日本語 - Название песни",
		}

		description := fetchDescription(texts)

		assert.Equal(t, "日本語 - Название песни", description)
	})
}

func Test_ParseSeconds(t *testing.T) {
	t.Run("when correct time string given", func(t *testing.T) {
		seconds := ParseSeconds("02:22:33")

		assert.Equal(t, 8553, seconds)
	})

	t.Run("when incorrect time string given", func(t *testing.T) {
		seconds := ParseSeconds("Seven seconds")

		assert.Equal(t, 0, seconds)
	})
}
