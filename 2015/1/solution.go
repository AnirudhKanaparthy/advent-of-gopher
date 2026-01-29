package main

import "fmt"
import "log"
import "os"

const inputFilePath string = "input.txt"

func calcFloor(instructions string) int {
    floor := 0
    for _, c := range instructions {
        if c == '(' {
            floor += 1
        } else if c == ')' {
            floor -= 1
        }
    }
    return floor
}

func main() {
    content, err := os.ReadFile(inputFilePath)
    if err != nil {
        log.Fatal("could not read input file")
    }
    fmt.Println(calcFloor(string(content)))
}
