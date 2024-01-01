package entity

type Question struct {
	ID              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswer   uint
	Difficulty      QuestionDifficulty
	CategoryID      uint
}

type PossibleAnswer struct {
	ID     uint
	Text   string
	Choice PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

func (c PossibleAnswerChoice) isValid() bool {
	if c >= PossibleAnswerA && c <= PossibleAnswerD {

		return true
	}

	return false

}

type QuestionDifficulty uint8

const (
	QuestionDifficultyEasy QuestionDifficulty = iota + 1
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (d QuestionDifficulty) IsValid() bool {
	if d >= QuestionDifficultyEasy && d <= QuestionDifficultyHard {

		return true
	}

	return false
}
