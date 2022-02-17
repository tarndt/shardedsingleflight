package shardedsingleflight

import (
	"math/big"
)

func nextPrime(n uint64) uint64 {
	switch n {
	case 0, 1:
		return 2
	case 2, 3:
		return n
	}

	if n%2 == 0 {
		n++
	}

	//There are very small gaps between consecutive small primes, brute force it
	i, one := big.NewInt(int64(n)), big.NewInt(1)
	for {
		if i.ProbablyPrime(0) { //ProbablyPrime is 100% accurate for inputs less than 2^64
			return i.Uint64()
		}
		i.Add(i, one)
	}
}
