package main

import (
    "fmt"
    "log"
    "os"
    "strings"
)

const inputFilePath string = "input.txt"

func IsNiceStringPart1(text string) bool {
    l := len(text)
    if l <= 3 {
        return false
    }

    twiceChars := 0
    for i := range text[:l-1] {
        s := text[i:i+2]
        switch s {
        case "ab": fallthrough
        case "cd": fallthrough
        case "pq": fallthrough
        case "xy":
            return false
        }
        if text[i] == text[i+1] {
            twiceChars += 1
        }
    }
    if twiceChars < 1 {
        return false
    }

    vowelsCount := 0
    for _, c := range text {
        switch c {
        case 'a': fallthrough
        case 'e': fallthrough
        case 'i': fallthrough
        case 'o': fallthrough
        case 'u': vowelsCount += 1
        }
    }

    return vowelsCount >= 3
}

func CheckCondition1(text string) bool {
    if len(text) < 4 {
        return false
    }

    pairMap := make(map[string]int)
    lastPair := ""
    for i := range text[:len(text)-1] {
        pair := text[i:i+2]

        flag := (lastPair == pair)
        lastPair = pair

        if flag {
            continue
        }

        val, ok := pairMap[pair]
        if !ok {
            pairMap[pair] = 0
            val = 0
        }

        pairMap[pair] = val + 1
        if pairMap[pair] >= 2 {
            return true
        }
    }
    return false
}

func CheckCondition2(text string) bool {
    if len(text) < 3 {
        return false
    }

    lastPair := text[:2]
    for i := range text[1:len(text)-1] {
        pair := text[i:i+2]

        if lastPair[0] == pair[1] && lastPair[1] == pair[0] {
            return true
        }

        lastPair = pair
    }
    return false
}

func IsNiceStringPart2(text string) bool {
    return CheckCondition1(text) && CheckCondition2(text)
}

func main() {
    data, err := os.ReadFile(inputFilePath)
    if err != nil {
        log.Fatal("ERROR - cannot read input file")
    }

    numNiceP1 := 0
    numNiceP2 := 0
    for _, line := range strings.Split(string(data), "\n") {
        if IsNiceStringPart1(line) {
            numNiceP1 += 1
        }
        if IsNiceStringPart2(line) {
            numNiceP2 += 1
        }
    }

    fmt.Printf("Nice Part 1: %v\n", numNiceP1)
    fmt.Printf("Nice Part 2: %v\n", numNiceP2)
}
