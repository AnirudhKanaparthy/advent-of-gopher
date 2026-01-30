package main

import "fmt"
import "log"
import "os"

const inputFilePath = "input.txt"

type Vec2i struct {
    x int
    y int
}

func GiftDelivery(instructions string) int {
    visited := make(map[Vec2i]bool)
    pos := Vec2i{0, 0}
    visited[pos] = true
    for _, c := range instructions {
        switch c {
        case '^':
            pos.y += 1
        case 'v':
            pos.y -= 1
        case '<':
            pos.x -= 1
        case '>':
            pos.x += 1
        default:
            continue
        }
        visited[pos] = true
    }
    return len(visited)
}

func GiftDeliveryWithRobot(instructions string) int {
    visited := make(map[Vec2i]bool)
    santa := Vec2i{0, 0}
    roboSanta := Vec2i{0, 0}
    visited[santa] = true
    for i, c := range instructions {
        ptr := &santa
        if i % 2 == 0 {
            ptr = &roboSanta
        }

        switch c {
        case '^':
            ptr.y += 1
        case 'v':
            ptr.y -= 1
        case '<':
            ptr.x -= 1
        case '>':
            ptr.x += 1
        default:
            continue
        }
        visited[*ptr] = true
    }
    return len(visited)
}

func main() {
    data, err := os.ReadFile(inputFilePath)
    if err != nil {
        log.Fatal("ERROR - cannot read input file")
    }

    fmt.Println(GiftDelivery(string(data)))
    fmt.Println(GiftDeliveryWithRobot(string(data)))
}
