package main

import (
	"fmt"
	"sync"
)

type result struct {
	apples       int
	oranges      int
	wg           sync.WaitGroup
	appleAction  chan func()
	orangeAction chan func()
	stop         chan struct{}
}

func newResult() *result {
	r := &result{
		appleAction:  make(chan func()),
		orangeAction: make(chan func()),
		stop:         make(chan struct{}),
	}
	go r.runApples()
	go r.runOranges()
	r.wg.Add(2)
	return r
}

func (r *result) runApples() {
	for {
		select {
		case f := <-r.appleAction:
			f()
		case <-r.stop:
			return
		}
	}
}

func (r *result) runOranges() {
	for {
		select {
		case f := <-r.orangeAction:
			f()
		case <-r.stop:
			return
		}
	}
}

func (r *result) calcApples(treeLoc, houseStart, houseEnd int, apples []int) {
	for _, a := range apples {
		tmp := treeLoc + a
		if tmp >= houseStart && tmp <= houseEnd {
			r.appleAction <- func() {
				r.apples++
			}
		}
	}
	r.appleAction <- func() {
		r.wg.Done()
	}
}

func (r *result) calcOranges(treeLoc, houseStart, houseEnd int, oranges []int) {
	for _, o := range oranges {
		tmp := treeLoc + o
		if tmp >= houseStart && tmp <= houseEnd {
			r.orangeAction <- func() {
				r.oranges++
			}
		}
	}
	r.orangeAction <- func() {
		r.wg.Done()
	}
}

func main() {
	var s, t, a, b, m, n int
	fmt.Scan(&s, &t, &a, &b, &m, &n)

	inApples := make([]int, m)
	inOranges := make([]int, n)
	for i := range inApples {
		fmt.Scan(&inApples[i])
	}
	for i := range inOranges {
		fmt.Scan(&inOranges[i])
	}

	// Calc results
	r := newResult()
	go r.calcApples(a, s, t, inApples)
	go r.calcOranges(b, s, t, inOranges)

	// Wait for results
	r.wg.Wait()
	close(r.stop)

	// Print results
	fmt.Println(r.apples)
	fmt.Println(r.oranges)
}
