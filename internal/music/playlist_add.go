package music

import (
    "amp_manager/internal/config"
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "net/http"
)

type PlaylistSong struct {
    Id   string `json:"id"`
    Type string `json:"type"`
}

type Add2PlaylistRequestData struct {
    Data []PlaylistSong `json:"data"`
}

func (apr *Add2PlaylistRequestData) Add(id ...string) {
    for _, amid := range id {
        apr.Data = append(apr.Data, PlaylistSong{
            Id:   amid,
            Type: "song",
        })
    }
}

func NewAdd2PlaylistRequest() *Add2PlaylistRequestData {
    return &Add2PlaylistRequestData{Data: make([]PlaylistSong, 0)}
}

func Add2Playlist(cfg *config.ManagerConfig, amidChan chan string) error {
    bulk := AmidBulk{
        MaxItems: 100,
        Values:   make([]string, 0, 100),
    }

    for amid := range amidChan {
        if bulk.CanAdd() {
            bulk.Add(amid)
            continue
        }

        logrus.Infof("Adding %d tracks to playlist %s", len(bulk.Values), cfg.PlaylistId)
        request := NewAdd2PlaylistRequest()
        request.Add(bulk.Flush()...)
        err := add(cfg, request)
        if err != nil {
            return err
        }
    }

    request := NewAdd2PlaylistRequest()
    request.Add(bulk.Flush()...)
    logrus.Infof("Adding last %d tracks to playlist %s", len(bulk.Values), cfg.PlaylistId)
    err := add(cfg, request)
    if err != nil {
        return err
    }

    return nil
}

func add(cfg *config.ManagerConfig, data *Add2PlaylistRequestData) error {
    url := fmt.Sprintf("https://amp-api.music.apple.com/v1/me/library/playlists/%s/tracks", cfg.PlaylistId)
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }

    request, err := http.NewRequest("POST", url, bytes.NewReader(body))
    if err != nil {
        return err
    }

    request.Header.Add("content-type", "application/json")
    request.Header.Add("authorization", fmt.Sprintf("Bearer %s", cfg.AuthToken))
    request.Header.Add("media-user-token", cfg.MediaToken)

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return err
    }

    if response.StatusCode != http.StatusNoContent {
        return fmt.Errorf("wrong response status, expected 204, got %d", response.StatusCode)
    }

    return nil
}
