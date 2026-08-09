package recommendation

type Recommendation struct {
	UserID    string `dynamodbav:"UserID"`
	ProductID string `dynamodbav:"ProductID"`
	Score     float64 `dynamodbav:"Score"`
}
