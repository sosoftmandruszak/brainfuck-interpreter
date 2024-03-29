package main

import (
	"brainfuck-interpreter/engine"
	"brainfuck-interpreter/utils"
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	var err error
	args := os.Args[1:]

	option := ""
	if len(args) > 0 {
		option = args[0]
	}

	switch option {
	case "-c":
		err = readFromConsole()
	case "-f":
		if len(args) < 2 {
			fmt.Println("Please give a file path")
			return
		}

		file := args[1]

		err = readFromFile(file)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func readFromConsole() error {
	ma := engine.NewMemoryAccess()
	stdinReader := bufio.NewReader(os.Stdin)
	scopeService := engine.NewLocalScopeService()
	commandSelector := commandSelector{
		memoryAccess: ma,
		writer:       utils.NewConsoleWriter(),
		reader:       stdinReader,
	}


	execute(commandSelector, stdinReader, scopeService)

	return nil
}

func readFromFile(path string) error {
	filReader, err := utils.NewFileReader(path)
	if err != nil {
		return err
	}

	defer filReader.Close()

	ma := engine.NewMemoryAccess()
	scopeService := engine.NewLocalScopeService()
	commandSelector := commandSelector{
		memoryAccess: ma,
		writer:       utils.NewConsoleWriter(),
		reader:       filReader,
	}

	execute(commandSelector, filReader, scopeService)

	return nil
}

func execute(commandSelector commandSelector, reader utils.ByteReader, scopeService engine.LocalScopeService, ) error {

	for {

		c, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		symbol := string(c)
		command := commandSelector.selectCommand(symbol)
		if command != nil {
			if err := engine.InterpretCommand(command, scopeService); err != nil {
				return err
			}
		}
	}

	return nil
}