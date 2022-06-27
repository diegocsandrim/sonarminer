package git

import "time"

type Commit struct {
	Id          int
	Hash        string
	ParentHash  string
	Date        time.Time
	Contributor *Contributor
	HasGoCode   bool
}

func NewCommit(id int, hash string, parentHash string, date time.Time, contributor *Contributor, hasGoCode bool) *Commit {
	c := Commit{
		Id:          id,
		Hash:        hash,
		ParentHash:  parentHash,
		Date:        date,
		Contributor: contributor,
		HasGoCode:   hasGoCode,
	}
	return &c
}
