// Chris Fowler -- June 11 2018 
// I was starting to teach myself Go before taking a break in May and June to write and defend my
// thesis. One thing I talked about was Ewens Sampling Formula, a distribution on S_n (the permutations
// on [n]) induced by a measure that weights cycle lengths. This was a straightforward exercise to generate
// such permutations in Go. A natural next step would be to verify invariable generation.

package main

import (
	"fmt"
	"math/rand"
	"bytes"
	"strconv"
	"time"
)

// Utilize the Feller Coupling to construct a random cycle structure, i.e. partition of n,
// with distribution ESF(n, a)
func ESFcycle(n int, a float64) []int {
	result := make([]int, 0, n)
	result = append(result, 1)
	for i, index := 0, -1; i < n; i++ {
		rnd := rand.Float64()
		if rnd < a/(a+float64(i)) {
			// if our random int is less than the mean we start a new cycle
			result = append(result, 1)
			index++
		} else {
			// else we continue the current cycle
			result[index]++
		}
	}
	return result
}

// Constructs a uniformly random permutation of n elements/integers
// I could have used rand.Perm(n) and then added 1 to each integer,
// but wanted to do it this way. This could be integrated into the for
// loop in ESFcycle for added efficiency.
func RandomPermutation(n int) []int {
	result := make([]int, n)
	// fill the array first with 1 through n in order
	for i := 0; i < n; i++ {
		result[i] = i + 1
	}
	
	for i := 0; i < n; i++ {
		// at each step we're equally likely to map to any of the remaining n - i elements
		next := int(rand.Int31n(int32(n-i)))
		result[i], result[i+next] = result[i+next], result[i]  // nice feature of Go
		// the value that was in result[i] gets moved to later in the array, to be placed
		// in a different slot
	}
	return result
}

func RandomESF(n int, a float64) string {
	cycles := ESFcycle(n, a)
	order := RandomPermutation(n)
	var buffer bytes.Buffer  // use a buffer to create the result quickly
	buffer.WriteString("(")
	for cycle, i := 0, 0; i < n; i++ {
		// write the next element into the cycle, then count of elements in current cycle
		// decrements by 1
		buffer.WriteString(strconv.Itoa(order[i]))
		cycles[cycle]--
		if cycles[cycle] == 0 {
			// if we've exhausted the current cycle we proceed to the next one
			cycle++
			buffer.WriteString(")")
			if i != n-1 {
				// we've reached the end of the permutation and don't need a "("
				buffer.WriteString("(")
			}
		} else {
			// if we're continuing with this cycle then put a space between integers
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
}

// Seeds the random generator and iterates through
func main() {
	// seed rand, self-explanatory
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 20; i++ {
		// make 20 ESF permutations, with differing alpha parameters
		// then print the permutations
		alpha := (2.0 - float64(i)/10.0)
		fmt.Println("the parameter is ", alpha, " and the perm is ", RandomESF(30, alpha))
	}
}