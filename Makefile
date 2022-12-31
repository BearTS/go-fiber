# Set the default target
.DEFAULT_GOAL := build

# Compile the project
build:
	go build -o bin/main .

# Run the program
run:
	./main

# Clean up
clean:
	rm -f main