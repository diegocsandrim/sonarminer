package strategy

import (
	"fmt"
	"log"
	"time"

	"github.com/diegocsandrim/sonarminer/git"
	"github.com/diegocsandrim/sonarminer/qualityanalyzers"
	"github.com/diegocsandrim/sonarminer/settings"
)

func AnalyseByPeriod(namespace string, project string, config settings.Config) error {
	gitRepo := git.NewGitRepo(namespace, project)

	err := gitRepo.Clone()
	if err != nil {
		return fmt.Errorf("could not clone repo: %w", err)
	}

	err = gitRepo.LoadCommits()
	if err != nil {
		return fmt.Errorf("could not load commits: %w", err)
	}

	qualityAnalyzer, err := qualityanalyzers.CreateSonnarAnalyser(
		qualityanalyzers.FormatProjectKey(namespace, project),
		config.SonarKey,
		config.SonarURL,
		gitRepo.ProjectDir(),
	)
	if err != nil {
		return err
	}

	defer qualityAnalyzer.Close()

	period := config.PeriodInterval
	monthlyCommits := gitRepo.CodeCommitsByPeriod(period)

	for i, monthCommits := range monthlyCommits {
		commit := getEarlyCommit(monthCommits.Commits)
		contributors := uniqueContributors(monthCommits.Commits)

		err = gitRepo.Checkout(commit.Hash)
		if err != nil {
			return err
		}

		shortCommitHash := commit.Hash[0:8]

		startTimestamp := time.Date(monthCommits.Month.Year, time.Month(monthCommits.Month.Period*period+1), 1, 0, 0, 0, 0, time.UTC)

		log.Printf("Analysing commit %s (batch %d/%d) - %s\n", shortCommitHash, i+1, len(monthlyCommits), startTimestamp)

		err = qualityAnalyzer.Run(shortCommitHash, startTimestamp, len(contributors))
		if err != nil {
			return err
		}
	}

	return nil
}

func uniqueContributors(commits []*git.Commit) []*git.Contributor {
	contruibutors := make([]*git.Contributor, 0, len(commits))
	hash := make(map[string]interface{}, 0)

	for _, commit := range commits {
		if !commit.HasGoCode {
			continue
		}
		_, exists := hash[commit.Contributor.Id]
		if !exists {
			hash[commit.Contributor.Id] = nil
			contruibutors = append(contruibutors, commit.Contributor)
		}
	}
	return contruibutors
}

func getEarlyCommit(commits []*git.Commit) *git.Commit {
	//TODO: improve using the first commit in the master branch if possible
	var early *git.Commit

	for _, commit := range commits {
		if early == nil {
			early = commit
			continue
		}

		if commit.Date.Before(early.Date) {
			early = commit
		}
	}
	return early
}
