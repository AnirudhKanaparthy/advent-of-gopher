package main

import "fmt"
import "log"
import "os"

const inputFilePath string = "input2.txt"

func calcBasementPosition(instructions string) int {
    floor := 0
    for i, c := range instructions {
        if c == '(' {
            floor += 1
        } else if c == ')' {
            floor -= 1
        }

        if floor == -1 {
            return i + 1
        }
    }
    return floor
}

func main() {
    content, err := os.ReadFile(inputFilePath)
    if err != nil {
        log.Fatal("could not read input file")
    }
    fmt.Println(calcBasementPosition(string(content)))
}
