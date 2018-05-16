package main

func main() {
	c := make(chan int)
	q := make(chan int)
	c <- 1
	q <- 1
	x, ok, y, ok := <-c, <-q
}
