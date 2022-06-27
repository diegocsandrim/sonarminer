package git

type ContributorAttractorCommit struct {
	Commit       *Commit
	Contributors []*Contributor
}

func NewContributorAttractorCommit(commit *Commit) *ContributorAttractorCommit {
	c := ContributorAttractorCommit{
		Commit:       commit,
		Contributors: make([]*Contributor, 0),
	}
	return &c
}

func (c *ContributorAttractorCommit) AddAttractedContributor(contributor *Contributor) {
	c.Contributors = append(c.Contributors, contributor)
}
