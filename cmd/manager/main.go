package main

import (
    "amp_manager/internal/config"
    "amp_manager/internal/input"
    "amp_manager/internal/music"
    "flag"
    "github.com/sirupsen/logrus"
)

var (
    cfgPath, filePath *string
    limit, offset     *int
)

func init() {
    cfgPath = flag.String("c", "config.json", "-c <path_to_config.json>")
    filePath = flag.String("f", "tracks.txt", "-f <path to input filePath>")
    limit = flag.Int("l", 1000, "-l <limit>")
    offset = flag.Int("o", 0, "-o <offset>")
}

func main() {
    flag.Parse()
    if *cfgPath == "" || *filePath == "" {
        flag.Usage()
        return
    }

    cfg, err := config.Load(*cfgPath)
    if err != nil {
        logrus.Fatal(err)
    }

    args := flag.Args()
    if len(args) == 0 {
        flag.Usage()
    }

    switch args[0] {
    case "add":
        err = add2Playlist(cfg, *filePath, *limit, *offset)
        if err != nil {
            logrus.Fatal(err)
        }
    case "delete":
        err = deleteFromPlaylist(cfg)
        if err != nil {
            logrus.Fatal(err)
        }
    default:
        flag.Usage()
    }

}

func add2Playlist(cfg *config.ManagerConfig, filePath string, limit, offset int) error {
    amidChan, err := input.ScanAmid(filePath, limit, offset)
    if err != nil {
        return err
    }

    err = music.Add2Playlist(cfg, amidChan)
    if err != nil {
        return err
    }

    return nil
}

func deleteFromPlaylist(cfg *config.ManagerConfig) error {
    err := music.DeleteFromPlaylist(cfg)
    if err != nil {
        return err
    }

    return nil
}
