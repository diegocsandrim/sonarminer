package strategy

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/diegocsandrim/sonarminer/git"
	"github.com/diegocsandrim/sonarminer/qualityanalyzers"
	"github.com/diegocsandrim/sonarminer/settings"
)

func Batch(namespace string, project string, config settings.Config) error {
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

	maxRunPerProject := float64(config.BatchSize)
	contributorAttractorCommitsLen := float64(len(contributorAttractorCommits))
	batchSize := int(math.Ceil(contributorAttractorCommitsLen / maxRunPerProject))

	batches := make([][]*git.ContributorAttractorCommit, 0, (len(contributorAttractorCommits)+batchSize-1)/batchSize)
	for batchSize < len(contributorAttractorCommits) {
		contributorAttractorCommits, batches = contributorAttractorCommits[batchSize:], append(batches, contributorAttractorCommits[0:batchSize:batchSize])
	}
	batches = append(batches, contributorAttractorCommits)

	day := time.Hour * 24
	fakeDate := time.Now().Add(-day * time.Duration(maxRunPerProject+1))
	for i, batch := range batches {
		analyseCommit := batch[0].Commit

		shortCommitHash := analyseCommit.Hash[0:8]
		log.Printf("Analysing commit %s (batch %d/%d) - %s\n", shortCommitHash, i+1, len(batches), fakeDate.UTC())

		totalContributors := 0
		for _, contributorAttractorCommit := range batch {
			totalContributors += len(contributorAttractorCommit.Contributors)
		}

		err = gitRepo.Checkout(analyseCommit.Hash)
		if err != nil {
			return fmt.Errorf("could not checkout to commit: %w", err)
		}

		err = qualityAnalyzer.Run(shortCommitHash, fakeDate, totalContributors)
		if err != nil {
			return fmt.Errorf("could not run analyser: %w", err)
		}
		fakeDate = fakeDate.Add(day)
	}

	return nil
}
