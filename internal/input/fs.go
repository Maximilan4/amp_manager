package input

import (
    "bufio"
    "github.com/sirupsen/logrus"
    "os"
)

func ScanAmid(path string, limit, offset int) (chan string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)
    isrcChan := make(chan string)
    go func() {
        err = scanAmid(scanner, isrcChan, limit, offset)
        if err != nil {
            logrus.Error(err)
        }
        err = file.Close()
        if err != nil {
            logrus.Error(err)
        }
    }()

    return isrcChan, nil
}

func scanAmid(scanner *bufio.Scanner, input chan string, limit, offset int) error {
    defer close(input)
    counter := 0
    take := 0
    for scanner.Scan() {
        counter++
        if counter < offset {
            continue
        }

        amid := scanner.Text()
        if amid == "" {
            continue
        }

        if take == limit {
            break
        }

        input <- amid
        take++
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}
