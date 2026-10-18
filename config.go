package main

import (
    "encoding/json"
    "net/url"
    "os"
    "path/filepath"
    "strings"
)

// Configuration locale, jamais embarquée dans l'exécutable ni versionnée.
func loadAPIBase() string {
    value := strings.TrimSpace(os.Getenv("CINECAM_API_BASE"))
    if value == "" {
        executable, err := os.Executable()
        if err != nil { return "" }
        data, err := os.ReadFile(filepath.Join(filepath.Dir(executable), "cinecam.local.json"))
        if err != nil { return "" }
        var config struct { APIBase string `json:"apiBase"` }
        if json.Unmarshal(data, &config) != nil { return "" }
        value = strings.TrimSpace(config.APIBase)
    }
    parsed, err := url.Parse(value)
    if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.RawQuery != "" || parsed.Fragment != "" { return "" }
    return strings.TrimRight(value, "/") + "/"
}
