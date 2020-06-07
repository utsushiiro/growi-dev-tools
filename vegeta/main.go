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
		_rate              = flag.Int("rate", -1, "request per second")
		_duration          = flag.Int("duration", -1, "how many seconds")
		isRandomAccessMode = flag.Bool("random-access", false, "random access mode")
		isRandomUpdateMode = flag.Bool("random-update", false, "random update mode")
	)
	flag.Parse()

	if *_rate < 1 || *_duration < 1 {
		fmt.Fprintf(os.Stderr, "rate, durationの値が不正です\n")
		os.Exit(1)
	}
	rate := vegeta.Rate{Freq: *_rate, Per: time.Second}
	duration := time.Duration(*_duration) * time.Second

	if *isRandomAccessMode == *isRandomUpdateMode {
		fmt.Fprintf(os.Stderr, "random-accessとrandom-updateのどちらかを指定してください\n")
		os.Exit(1)
	}

	factory, err := growi.NewGrowiTargeterFactory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var targeter vegeta.Targeter
	if *isRandomAccessMode {
		targeter, err = factory.NewRandomPageAccessTargeter()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	} else if *isRandomUpdateMode {
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
