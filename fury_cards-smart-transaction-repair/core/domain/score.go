package domain

type ScoreInput struct {
	ScoreData TransactionData
}

type ScoreOutput struct {
	Result string
	Score  float64
}
