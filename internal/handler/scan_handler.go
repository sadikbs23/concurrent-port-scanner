package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"concurrent-port-scanner/internal/scanner"
)

type ScanHandler struct {
	tmpl *template.Template
}

func NewScanHandler() *ScanHandler {
	tmpl := template.Must(template.ParseFiles(filepath.Join("internal", "templates", "form.html")))
	return &ScanHandler{tmpl: tmpl}
}

func (h *ScanHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.tmpl.Execute(w, nil)
}

func (h *ScanHandler) HandleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	host := r.FormValue("host")
	startPort, err := strconv.Atoi(r.FormValue("startPort"))
	if err != nil {
		http.Error(w, "Invalid starting port", http.StatusBadRequest)
		return
	}

	endPort, err := strconv.Atoi(r.FormValue("endPort"))
	if err != nil {
		http.Error(w, "Invalid ending port", http.StatusBadRequest)
		return
	}

	ps := scanner.NewPortScanner(host, startPort, endPort)
	results := ps.ScanPorts()

	w.Header().Set("Content-Type", "text/html")
	data := struct {
		Results []scanner.ScanResult
	}{
		Results: results,
	}

	h.tmpl.Execute(w, data)
}
