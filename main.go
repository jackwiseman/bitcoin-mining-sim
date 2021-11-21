package main

func main() {

	ch := make(chan Block)
	go run logger(ch)

}
