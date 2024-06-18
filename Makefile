demo:
	go build -o ./bin/demo ./cmd

clean:
	rm bin/demo

.PHONY: demo