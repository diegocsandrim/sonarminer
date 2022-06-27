package git

type MonthCommits struct {
	Month   YearPeriod
	Commits []*Commit
}

type YearPeriod struct {
	Period int
	Year   int
}
