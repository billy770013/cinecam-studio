// CinéCam Studio — interface locale de caméra pour Minecraft Bedrock.
package main

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var apiBase string

//go:embed web
var webFS embed.FS

var httpClient = &http.Client{Timeout: 30 * time.Second}

// relaye une requete du navigateur local vers le serveur Minecraft (lien secret)
func relais(chemin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiBase == "" {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusServiceUnavailable)
            w.Write([]byte(`{"error":"Configurez apiBase dans cinecam.local.json, puis relancez CineCam Studio."}`))
            return
        }
        url := apiBase + chemin
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}
		var corps io.Reader
		if r.Body != nil {
			b, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
			corps = bytes.NewReader(b)
		}
		req, err := http.NewRequest(r.Method, url, corps)
		if err != nil {
			w.WriteHeader(502)
			w.Write([]byte(`{"error":"requête invalide"}`))
			return
		}
		req.Header.Set("Content-Type", "application/json")
		rep, err := httpClient.Do(req)
		if err != nil {
			w.WriteHeader(502)
			w.Write([]byte(`{"error":"serveur injoignable (internet ?)"}`))
			return
		}
		defer rep.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rep.StatusCode)
		io.Copy(w, rep.Body)
	}
}

func main() {
    apiBase = loadAPIBase()
	adresse := net.JoinHostPort("127.0.0.1", "47931")
	url := "http://" + adresse + "/"
	ln, err := net.Listen("tcp", adresse)
	if err != nil {
		ouvrirFenetre(url) // deja ouvert : on remontre la fenetre
		return
	}

	sousWeb, _ := fs.Sub(webFS, "web")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sousWeb)))
	mux.HandleFunc("/api/players", relais("players"))
	mux.HandleFunc("/api/where", relais("where"))
	mux.HandleFunc("/api/tool", relais("tool"))
	mux.HandleFunc("/api/play", relais("play"))
	mux.HandleFunc("/api/stop", relais("stop"))

	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)

	ouvrirFenetre(url)

	// la fenetre fermee -> on s'arrete (sauf si on tourne depuis peu). Simple : on attend un signal
	// ou on surveille une derniere activite via /api/ping (heartbeat du navigateur).
	dernier := time.Now()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) { dernier = time.Now(); w.Write([]byte("ok")) })

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	tick := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-sig:
			return
		case <-tick.C:
			if time.Since(dernier) > 90*time.Second {
				return // plus de heartbeat depuis 1m30 : la fenetre est fermee
			}
		}
	}
}

// evite l'erreur "imported and not used" si un helper OS n'utilise pas exec/strings
var _ = exec.Command
var _ = strings.TrimSpace
