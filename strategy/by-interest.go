package strategy

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/diegocsandrim/sonarminer/git"
	"github.com/diegocsandrim/sonarminer/qualityanalyzers"
	"github.com/diegocsandrim/sonarminer/settings"
)

func InterestNewContributor(namespace string, project string, config settings.Config) error {
	gitRepo := git.NewGitRepo(namespace, project)

	err := gitRepo.Clone()
	if err != nil {
		return fmt.Errorf("could not clone repo: %w", err)
	}

	err = gitRepo.LoadCommits()
	if err != nil {
		return fmt.Errorf("could not load commits: %w", err)
	}

	contributorAttractorCommits := gitRepo.ContributorAttractorCommits()
	sort.Slice(contributorAttractorCommits, func(i, j int) bool {
		commitI := contributorAttractorCommits[i].Commit
		commitJ := contributorAttractorCommits[j].Commit
		return commitI.Id < commitJ.Id
	})

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

	day := time.Hour * 24
	analysisDate := time.Now().UTC().Add(-day * time.Duration(len(contributorAttractorCommits)))

	for i, contributorAttractorCommit := range contributorAttractorCommits {
		analysisDate = analysisDate.Add(day)
		shortCommitHash := contributorAttractorCommit.Commit.Hash[0:8]
		log.Printf("Analysing commit %s (%d/%d)\n", shortCommitHash, i+1, len(contributorAttractorCommits))
		err = gitRepo.Checkout(contributorAttractorCommit.Commit.Hash)
		if err != nil {
			return fmt.Errorf("could not checkout to commit: %w", err)
		}

		attractedContributors := len(contributorAttractorCommit.Contributors)

		err = qualityAnalyzer.Run(shortCommitHash, analysisDate, attractedContributors)
		if err != nil {
			return fmt.Errorf("could not run analyser: %w", err)
		}
	}
	return nil
}
