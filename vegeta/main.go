package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	vegeta "github.com/tsenart/vegeta/lib"
	"github.com/utsushiiro/growi-dev-tools/vegeta/growi"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, ".envファイルの読み込みに失敗しました\n")
		os.Exit(1)
	}

	var (
		_rate                  = flag.Int("rate", -1, "request per second")
		_duration              = flag.Int("duration", -1, "how many seconds")
		isRandomPageAccessMode = flag.Bool("random-page-access", false, "random page access mode")
		isRandomPageUpdateMode = flag.Bool("random-page-update", false, "random page update mode")
	)
	flag.Parse()

	if *_rate < 1 || *_duration < 1 {
		fmt.Fprintf(os.Stderr, "rate, durationの値が不正です\n")
		os.Exit(1)
	}
	rate := vegeta.Rate{Freq: *_rate, Per: time.Second}
	duration := time.Duration(*_duration) * time.Second

	if *isRandomPageAccessMode == *isRandomPageUpdateMode {
		fmt.Fprintf(os.Stderr, "random-page-accessとrandom-page-updateのどちらかを指定してください\n")
		os.Exit(1)
	}

	factory, err := growi.NewGrowiTargeterFactory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var targeter vegeta.Targeter
	if *isRandomPageAccessMode {
		targeter, err = factory.NewRandomPageAccessTargeter()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	} else if *isRandomPageUpdateMode {
		targeter, err = factory.NewRandomPageUpdateTargeter()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "存在しないmodeです\n")
		os.Exit(1)
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "GROWI") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	fmt.Printf("## Vegeta Metrics\n\n")
	reporter(os.Stdout)
}
