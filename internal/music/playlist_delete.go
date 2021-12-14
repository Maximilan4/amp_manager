package music

import (
    "amp_manager/internal/config"
    "context"
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "golang.org/x/sync/errgroup"
    "io"
    "net/http"
    "net/url"
    "strconv"
)

func DeleteFromPlaylist(cfg *config.ManagerConfig) error {
    hrefs := make(chan string)
    go func() {
        defer close(hrefs)
        err := getInfo(cfg, hrefs, 0)
        if err != nil {
            logrus.WithError(err).Error("unable to get")
        }
    }()

    group, _ := errgroup.WithContext(context.Background())
    for i := 0; i < 3; i++ {
        group.Go(func() error {
            for href := range hrefs {
                logrus.Infof("deleting %s", href)
                err := deleteRequest(cfg, href)
                if err != nil {
                    logrus.WithError(err).Errorf("unable to delete %s", href)
                }
            }

            return nil
        })
    }

    if err := group.Wait(); err != nil {
        return err
    }

    return nil
}

func deleteRequest(cfg *config.ManagerConfig, href string) error {
    parsedUrl, err := url.Parse("https://amp-api.music.apple.com" + href)
    if err != nil {
        return err
    }

    request, err := http.NewRequest("DELETE", parsedUrl.String(), nil)
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
        body, _ := io.ReadAll(response.Body)
        response.Body.Close()
        logrus.Error(string(body))
        return fmt.Errorf("delete response with status %d", response.StatusCode)
    }

    return nil
}

func getInfo(cfg *config.ManagerConfig, hrefs chan string, offset int) error {
    parsedUrl, err := url.Parse(
        fmt.Sprintf("https://amp-api.music.apple.com/v1/me/library/playlists/%s/tracks", cfg.PlaylistId),
    )
    if err != nil {
        return err
    }

    query := url.Values{}
    query.Add("fields", "href")
    query.Add("offset", strconv.Itoa(offset))
    parsedUrl.RawQuery = query.Encode()

    request, err := http.NewRequest("GET", parsedUrl.String(), nil)
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

    if response.StatusCode != http.StatusOK {
        return fmt.Errorf("info response with status %d", response.StatusCode)
    }

    defer response.Body.Close()

    var info PlaylistTracks
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&info)
    if err != nil {
        return err
    }

    for _, track := range info.Data {
        hrefs <- track.Href
    }

    if info.Meta.Total > len(info.Data) {
        left := info.Meta.Total - (offset + len(info.Data))
        var newOffset int
        if left == 0 {
            return nil
        } else {
            newOffset = offset + 100
        }

        err = getInfo(cfg, hrefs, newOffset)
        if err != nil {
            return err
        }
    }

    return nil
}

type PlaylistTracks struct {
    Data []struct {
        Href string `json:"href"`
    } `json:"data"`
    Meta struct {
        Total int `json:"total"`
    } `json:"meta"`
}
