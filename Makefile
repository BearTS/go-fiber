# Set the default target
.DEFAULT_GOAL := build

# Compile the project
build:
	go build -o bin/app .

# Run the program
run:
	./app

# Clean up
clean:
	rm -f ./bin/app