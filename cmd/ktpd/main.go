package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/ktpready"
	"github.com/fahmifan/ktpready/https"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}

func run(args []string) error {
	dirtyWordFile, err := os.Open("./corpus/dirty_words.txt")
	if err != nil {
		return fmt.Errorf("open badwords: %w", err)
	}
	defer dirtyWordFile.Close()

	bannedWordFile, err := os.Open("./corpus/banned_words.txt")
	if err != nil {
		return fmt.Errorf("open badwords: %w", err)
	}
	defer bannedWordFile.Close()

	nameChecker := &ktpready.NameChecker{MinWords: 2}
	nameChecker.LoadBannedWords(bannedWordFile)
	nameChecker.LoadDirtyWords(dirtyWordFile)

	server := https.NewServer(nameChecker)
	err = server.Run()
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
