package challenge

import (
	"crypto/rand"

	"github.com/bwesterb/go-pow"
)

type Challenge struct {
	difficulty uint32
}

func WithDifficulty(difficulty uint32) Challenge {
	return Challenge{
		difficulty: difficulty,
	}
}

func (c Challenge) New() string {
	nonce := make([]byte, 64)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	return pow.NewRequest(c.difficulty, nonce)
}
