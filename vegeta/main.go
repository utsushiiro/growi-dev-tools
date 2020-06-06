package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	vegeta "github.com/tsenart/vegeta/lib"
	"github.com/utsushiiro/growi-dev-tools/vegeta/growi"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".envファイルの読み込みに失敗しました")
	}

	var (
		_rate     = flag.Int("rate", -1, "request per second")
		_duration = flag.Int("duration", -1, "how many seconds")
	)
	flag.Parse()
	if *_rate < 1 || *_duration < 1 {
		log.Fatal("引数が不正です")
	}
	rate := vegeta.Rate{Freq: *_rate, Per: time.Second}
	duration := time.Duration(*_duration) * time.Second

	targeter, err := growi.NewRandomPageAccessTargeter()
	if err != nil {
		log.Fatalf("%+v\n", err)
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
