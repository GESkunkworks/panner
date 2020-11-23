package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GESkunkworks/dustcollector"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/inconshreveable/log15"
)

// global version
var version string

// loggo is the global logger
var loggo log15.Logger

// setLogger sets up logging globally for the packages involved
func setLogger(noLogFile bool, logFileS, loglevel string) {
	loggo = log15.New()
	if noLogFile && loglevel == "debug" {
		loggo.SetHandler(
			log15.LvlFilterHandler(
				log15.LvlDebug,
				log15.StreamHandler(os.Stdout, log15.LogfmtFormat())))
	} else if noLogFile {
		loggo.SetHandler(
			log15.LvlFilterHandler(
				log15.LvlInfo,
				log15.StreamHandler(os.Stdout, log15.LogfmtFormat())))
	} else if loglevel == "debug" {
		// log to stdout and file
		loggo.SetHandler(log15.MultiHandler(
			log15.StreamHandler(os.Stdout, log15.LogfmtFormat()),
			log15.LvlFilterHandler(
				log15.LvlDebug,
				log15.Must.FileHandler(logFileS, log15.JsonFormat()))))
	} else {
		// log to stdout and file
		loggo.SetHandler(log15.MultiHandler(
			log15.LvlFilterHandler(
				log15.LvlInfo,
				log15.StreamHandler(os.Stdout, log15.LogfmtFormat())),
			log15.LvlFilterHandler(
				log15.LvlInfo,
				log15.Must.FileHandler(logFileS, log15.JsonFormat()))))
	}
}

// generic error handler for easy typing
func errorhandle(err error) {
	if err != nil {
		loggo.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	var profile string
	var region string
	var dateFilter string
	var logFile string
	var logLevel string
	var outFileSnap string
	var outFileSummary string
	var outFileBars string
	var noLogFile, versionFlag, showSummary bool
	var pageSize, maxPages, volBatchSize int
	var ebsSnapRate float64
	flag.StringVar(&profile, "profile", "default", "AWS session credentials profile")
	flag.StringVar(&region, "region", "us-east-1", "AWS Region")
	flag.StringVar(&logFile, "logfile", "panner.log.json", "JSON logfile location")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level (info or debug)")
	flag.StringVar(&dateFilter, "datefilter", "2018-01-01",
		"only analyze snapshots created before this date")
	flag.StringVar(&outFileSnap, "outfile-snapshots", "out-snap.csv",
		"filename of csv output file that contains all snapshots that meet dateFilter criteria")
	flag.StringVar(&outFileBars, "outfile-bars", "out-bars.csv",
		"filename of csv output file that contains snapshots "+
			"aggregated by common volume (useful in determining "+
			"potential cost savings)")
	flag.StringVar(&outFileSummary, "outfile-summary", "out-summary.txt",
		"filename of text output file that shows summary of action plan for this account.")
	flag.BoolVar(&showSummary, "show-summary", false,
		"when set summary is output to stdout as well as being written to --outfile-summary")
	flag.IntVar(&maxPages, "max-pages", 5,
		"maximum number of pages to pull during describe snapshots call")
	flag.IntVar(&volBatchSize, "describe-volumes-batch-size", 20,
		"make this larger if you hit throttling limits. "+
			"Make this smaller if you want to speed up the program.")
	flag.IntVar(&pageSize, "pagesize", 500, "number of snapshots to pull per page")
	flag.BoolVar(&versionFlag, "v", false, "print version and exit")
	flag.Float64Var(&ebsSnapRate, "ebs-snap-rate", 0.05,
		"per GB-month cost for EBS snapshot (used in analysis summary at end of script)")
	flag.BoolVar(&noLogFile, "nologfile", false,
		"Indicates whether or not to skip writing of a filesystem log file.")
	flag.Parse()
	if versionFlag {
		fmt.Printf("panner %s\n", version)
		os.Exit(0)
	}
	setLogger(noLogFile, logFile, logLevel)
	loggo.Info("Starting panner")
	var sess *session.Session
	loggo.Info("starting session", "profile", profile)
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))
	einput := dustcollector.ExpeditionInput{
		Session:                sess,
		MaxPages:               &maxPages,
		PageSize:               &pageSize,
		VolumeBatchSize:        &volBatchSize,
		DateFilter:             &dateFilter,
		Logger:                 &loggo,
		OutfileRecommendations: &outFileSummary,
		OutfileNuggets:         &outFileSnap,
		OutfileBars:            &outFileBars,
		EbsSnapRate:            &ebsSnapRate,
	}
	exp, err := dustcollector.New(&einput)
	errorhandle(err)
	err = exp.Start()
	if err != nil {
		loggo.Error("error running expedition", "error", err.Error())
		os.Exit(1)
	}
	errorhandle(err)
	loggo.Info("Writing snapshots to file")
	err = exp.ExportNuggets()
	errorhandle(err)
	loggo.Info("Writing cost info to file")
	err = exp.ExportBars()
	errorhandle(err)
	// now show/export action plan
	if showSummary {
		for _, line := range exp.GetRecommendations() {
			fmt.Println(line)
		}
	}
	err = exp.ExportRecommendations()
	errorhandle(err)
}
