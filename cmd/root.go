package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/fenrirunbound/pipeline-queue/internal/client"
)

// AccessToken string
var AccessToken string

// Hostname ...
var Hostname string

// Interval ...
var Interval time.Duration

// ProjectID ...
var ProjectID int

// PipelineID ...
var PipelineID int

func errorExit(err error) {
	fmt.Fprint(os.Stderr, "[error] ")
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func timedPrint(message string) {
	now := time.Now()
	fmt.Printf("[%v] %v", now, message)
}

func waitItOut(duration time.Duration) {
	time.Sleep(duration)
}

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Waits for older pipelines to complete before going",
	Long:  `Longer description goes here`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.New(nil, Hostname, AccessToken)
		if err != nil {
			errorExit(err)
		}
		projID := strconv.Itoa(ProjectID)
		pipeID := strconv.Itoa(PipelineID)

		areWeFirst := false

		for !areWeFirst {
			timedPrint("Checking if we're first in line...")
			areWeFirst, err = client.DetermineIfFirst(projID, pipeID)
			if err != nil {
				errorExit(err)
			}

			if !areWeFirst {
				timedPrint(fmt.Sprintf("We're not first. Trying again in %v", Interval))
				waitItOut(Interval)
			}
		}
	},
}

func init() {
	ciProjectID := os.Getenv("CI_PROJECT_ID")
	ciPipelineID := os.Getenv("CI_PIPELINE_ID")
	defaultProjectID, _ := strconv.Atoi(ciProjectID)
	defaultPipelineID, _ := strconv.Atoi(ciPipelineID)

	rootCmd.PersistentFlags().StringVarP(&AccessToken, "token", "t", os.Getenv("CI_JOB_TOKEN"), "API access token. Defaults to $CI_JOB_TOKEN")
	rootCmd.PersistentFlags().StringVarP(&Hostname, "hostname", "n", "https://gitlab.com", "Hostname of the Gitlab instance. Defaults to https://gitlab.com")
	rootCmd.PersistentFlags().DurationVarP(&Interval, "interval-time", "it", 30*time.Second, "Amount of time to wait in-between polls in time.Duration format. Default is 30s")
	rootCmd.PersistentFlags().IntVarP(&ProjectID, "project", "j", defaultProjectID, "Project ID of the pipeline to run in. Defaults to $CI_PROJECT_ID")
	rootCmd.PersistentFlags().IntVarP(&PipelineID, "pipeline", "l", defaultPipelineID, "Pipeline ID of the current pipeline. Defaults to $CI_PIPELINE_ID")
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
