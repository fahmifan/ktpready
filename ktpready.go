package ktpready

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	DefaultMaxWords = 60
	DefaultMinWords = 2
)

type Dictionary map[string]struct{}

type NameChecker struct {
	MinWords    int // default to no minimum
	MaxChar     int
	DirtyWords  Dictionary
	BannedWords Dictionary
}

// NewNameChecker
func NewNameChecker(minWord, maxWord int) *NameChecker {
	if maxWord < 0 {
		maxWord = DefaultMaxWords
	}
	if minWord < 0 {
		minWord = DefaultMinWords
	}
	return &NameChecker{
		MinWords: minWord,
		MaxChar:  maxWord,
	}
}

func (n *NameChecker) LoadDirtyWords(src io.Reader) error {
	dict, err := n.loadDictionary(src)
	if err != nil {
		return fmt.Errorf("load dirty words dictionary: %w", err)
	}
	n.DirtyWords = dict
	return nil
}

func (n *NameChecker) LoadBannedWords(src io.Reader) error {
	dict, err := n.loadDictionary(src)
	if err != nil {
		return fmt.Errorf("load banned words dictionary: %w", err)
	}
	n.BannedWords = dict
	return nil
}

func (n *NameChecker) loadDictionary(src io.Reader) (Dictionary, error) {
	bt, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("read all src: %w", err)
	}
	words := strings.Split(string(bt), "\n")
	dict := make(Dictionary, len(words))
	for _, word := range words {
		if strings.TrimSpace(word) == "" {
			continue
		}
		dict[strings.ToLower(word)] = struct{}{}
	}
	return dict, nil
}

// Check checks if the name is KTP ready or not
func (n *NameChecker) Check(name string) error {
	if name == "" {
		return errors.New("check: name is empty")
	}

	nameParts := strings.Split(name, " ")
	countWord := 0
	for _, part := range nameParts {
		part = strings.ToLower(part)
		if hasNonAlphabet(part) {
			return fmt.Errorf("check: name has non alphabet character (%s)", part)
		}
		if _, isDirty := n.DirtyWords[part]; isDirty {
			return fmt.Errorf("check: name has dirty word (%s)", part)
		}
		if _, isBanned := n.BannedWords[part]; isBanned {
			return fmt.Errorf("check: name has banned word (%s)", part)
		}

		countWord += len(part)
		if countWord > n.MaxChar {
			return fmt.Errorf("check: name too long, max characters are %d", n.MaxChar)
		}
	}

	if partLen := len(nameParts); partLen <= 0 || partLen < n.MinWords {
		return fmt.Errorf("check: name too short, minimum are %d", n.MinWords)
	}

	return nil
}

var rgxNonWord = regexp.MustCompile(`(\W)`)

func hasNonAlphabet(in string) bool {
	return rgxNonWord.MatchString(in)
}
