package main

import (
    "cmp"
    "fmt"
    "log"
    "os"
    "slices"
    "strconv"
    "strings"
)

const inputFilePath string = "input.txt"

type Box struct {
    l int
    w int
    h int
}

type BoxError struct {
    errorMessage string
}

func MakeBoxError(errMsg string) BoxError {
    return BoxError{errMsg}
}

func (be *BoxError) Error() string {
    return be.errorMessage
}

func min[K cmp.Ordered](a K, b K) K {
    if a < b {
        return a;
    } else {
        return b;
    }
}

func max[K cmp.Ordered](a K, b K) K {
    if a > b {
        return a;
    } else {
        return b;
    }
}

func (b *Box)GiftWrapNeeded() int {
    p := b.l*b.w
    q := b.w*b.h
    r := b.h*b.l

    s := min(min(p, q), r)

    return 2*p + 2*q + 2*r + s
}

func (b *Box)RibbonNeededV1() int {
    s := []int{b.l, b.h, b.w}
    slices.Sort(s)
    return 2*(s[0] + s[1]) + (b.l*b.h*b.w)
}

func (b *Box)RibbonNeeded() int {
    // A more fun approach

    p := min(b.l, b.h)
    q := min(b.h, b.w)
    r := min(b.w, b.l)

    m := max(p, q)
    n := max(q, r)
    o := max(r, p)

    t := (p + q + r + m + n + o) / 3
    return 2*t + (b.l*b.w*b.h)
}

func ParseIntoBox(text string) (Box, error) {
    if len(text) == 0 {
        be := MakeBoxError("empty input")
        return Box{}, &be
    }

    dimsRaw := strings.Split(text, "x")

    dims := [3]int{0, 0, 0}
    for i, dimRaw := range dimsRaw {
        dim, err := strconv.Atoi(dimRaw)
        if err != nil {
            be := MakeBoxError("wrong box format")
            return Box{}, &be
        }
        dims[i] = dim
    }

    return Box{dims[0], dims[1], dims[2]}, nil
}

func main() {
    data, err := os.ReadFile(inputFilePath)
    if err != nil {
        log.Fatal("ERROR - Could not open file")
    }

    totalRibbonNeeded := 0
    totalWrapNeeded := 0
    text := string(data)
    for _, line := range strings.Split(text, "\n") {
        if len(line) == 0 {
            continue
        }

        box, err := ParseIntoBox(line)
        if err != nil {
            log.Fatal("ERROR - ", err.Error())
            continue
        }
        totalWrapNeeded += box.GiftWrapNeeded()
        totalRibbonNeeded += box.RibbonNeeded()
    }
    fmt.Println(totalWrapNeeded)
    fmt.Println(totalRibbonNeeded)
}
