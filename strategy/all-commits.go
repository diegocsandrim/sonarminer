package strategy

import (
	"fmt"
	"log"
	"time"

	"github.com/diegocsandrim/sonarminer/git"
	"github.com/diegocsandrim/sonarminer/qualityanalyzers"
	"github.com/diegocsandrim/sonarminer/settings"
)

func AllCommits(namespace string, project string, config settings.Config) error {
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

	commits := gitRepo.Commits()
	day := time.Hour * 24
	fakeDate := time.Now().Add(-day * time.Duration(len(commits)))

	for i, commit := range commits {
		shortCommitHash := commit.Hash[0:8]

		log.Printf("Analysing commit %s (%d/%d) as if it was in %s\n", shortCommitHash, i+1, len(commits), commit.Date.String())
		contributors := 1

		err = qualityAnalyzer.Run(shortCommitHash, fakeDate, contributors)
		if err != nil {
			return err
		}
		fakeDate = fakeDate.Add(day)
	}

	return nil
}
