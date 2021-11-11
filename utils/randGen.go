package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GetRand() string {
	rand.Seed(time.Now().Unix())
	return strconv.FormatInt(rand.Int63n(1<<62), 10)
}
