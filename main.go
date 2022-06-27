package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/diegocsandrim/sonarminer/settings"
	"github.com/diegocsandrim/sonarminer/sonar"
	"github.com/diegocsandrim/sonarminer/strategy"
	"github.com/urfave/cli/v2"
)

func main() {
	config := settings.Config{
		SonarKey: "",
		SonarURL: "",
	}

	app := &cli.App{
		Usage: "analyse the project history using sonnar-scanner",
		Commands: []*cli.Command{
			{
				Name:    "analyse",
				Aliases: []string{"a"},
				Usage:   "analyse the repository history",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "sonarkey",
						Usage:       "Sonarqube token",
						EnvVars:     []string{"SONAR_TOKEN"},
						Destination: &(config.SonarKey),
					},
					&cli.StringFlag{
						Name:        "sonarurl",
						Usage:       "Sonarqube URL",
						EnvVars:     []string{"SONAR_URL"},
						Value:       "http://127.0.0.1:9000",
						Destination: &(config.SonarURL),
					},
					&cli.StringFlag{
						Name:        "strategy",
						Usage:       "Strategy to analyse the repositories, one of: ALL, PERIOD, BATCH, INTEREST",
						Value:       "PERIOD",
						Destination: &(config.Strategy),
					},
					&cli.IntFlag{
						Name:        "interval",
						Usage:       "When using strategy=PERIOD, set the interval in months",
						Value:       6,
						Destination: &(config.PeriodInterval),
					},
					&cli.IntFlag{
						Name:        "batch",
						Usage:       "Batch size when using strategy=BATCH",
						Value:       20,
						Destination: &(config.BatchSize),
					},
				},
				Action: func(c *cli.Context) error {
					if config.SonarKey == "" {
						token, err := sonar.NewToken(config.SonarURL)
						if err != nil {
							return fmt.Errorf("token not provided, fail to create one with default user/password: %w", err)
						}
						config.SonarKey = token
					}

					if c.Args().Len() == 0 {
						return fmt.Errorf("must provide at least one repository to analyse")
					}
					repositories := c.Args().Slice()
					for _, repository := range repositories {
						err := runRepository(config, repository)
						if err != nil {
							return cli.Exit(err.Error(), 1)
						}
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runRepository(config settings.Config, repositoryFullName string) error {
	repositoryParts := strings.Split(repositoryFullName, "/")
	if len(repositoryParts) != 2 {
		return cli.Exit("argument must be in format namespace/project", 1)
	}

	namespace := repositoryParts[0]
	project := repositoryParts[1]

	log.Printf("starting namespace %s, project %s", namespace, project)

	var err error

	switch config.Strategy {
	case "ALL":
		err = strategy.AllCommits(namespace, project, config)
	case "PERIOD":
		err = strategy.AnalyseByPeriod(namespace, project, config)
	case "BATCH":
		err = strategy.Batch(namespace, project, config)
	case "INTEREST":
		err = strategy.InterestNewContributor(namespace, project, config)
	default:
		return fmt.Errorf("unknown strategy: %s", config.Strategy)
	}

	if err != nil {
		return err
	}

	log.Printf("finished namespace %s, project %s", namespace, project)

	return nil
}
