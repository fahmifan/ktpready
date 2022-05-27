package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/fahmifan/ktpready"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
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

	err = nameChecker.LoadDirtyWords(dirtyWordFile)
	if err != nil {
		return fmt.Errorf("load dirty words: %w", err)
	}

	err = nameChecker.LoadBannedWords(bannedWordFile)
	if err != nil {
		return fmt.Errorf("load dirty words: %w", err)
	}

	if len(args) == 1 {
		return nil
	}

	err = nameChecker.Check(args[1])
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}

	fmt.Println("name seems legit!")
	return nil
}
