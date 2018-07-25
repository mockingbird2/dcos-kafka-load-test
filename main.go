package main

func main() {
	config := ParseInput()
	createdMessages := make(chan []byte, 100)
	creator := MessageCreator(*config, createdMessages)
	creator.StartCreators()
	creator.StopCreators()
}
