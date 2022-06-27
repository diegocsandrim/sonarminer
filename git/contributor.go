package git

type Contributor struct {
	Id            string
	Commits       []*Commit
	firstCommit   *Commit
	firstGoCommit *Commit
}

func NewContributor(id string) *Contributor {
	c := Contributor{
		Id:      id,
		Commits: make([]*Commit, 0),
	}
	return &c
}

func (c *Contributor) AddCommit(commit *Commit) {
	c.Commits = append(c.Commits, commit)

	if c.firstCommit == nil {
		c.firstCommit = commit
	}

	if c.firstGoCommit == nil && commit.HasGoCode {
		c.firstGoCommit = commit
	}

	if commit.Date.Before(c.firstCommit.Date) {
		c.firstCommit = commit
	}
}

func (c *Contributor) FirstCommit() *Commit {
	return c.firstCommit
}

func (c *Contributor) FirstGoCommit() *Commit {
	return c.firstGoCommit
}

func (c *Contributor) IsMainContributor() bool {
	return c.FirstCommit().ParentHash == ""
}
