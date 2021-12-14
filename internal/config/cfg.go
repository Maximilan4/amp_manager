package config

import (
    "encoding/json"
    "os"
)

type ManagerConfig struct {
    PlaylistId string `json:"playlist_id"`
    AuthToken  string `json:"auth_token"`
    MediaToken string `json:"media_token"`
}

func Load(configPath string) (*ManagerConfig, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }

    decoder := json.NewDecoder(file)
    var cfg ManagerConfig
    err = decoder.Decode(&cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}
