package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sahib002/deployflow/internal/core/ports"
)

type ProxyServer struct {
	repo ports.RepositoryPort
}

func NewProxyServer(repo ports.RepositoryPort) *ProxyServer {
	return &ProxyServer{repo: repo}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Get the Subdomain (e.g., "cat-dog" from "cat-dog.localhost:8000")
	host := r.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		http.Error(w, "Invalid domain format. Use project-name.localhost", http.StatusBadRequest)
		return
	}
	projectName := parts[0]

	// 2. Find the Project
	fmt.Printf("🔎 Proxy looking for project: %s\n", projectName)
	project, err := p.repo.GetProjectByName(r.Context(), projectName)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// 3. Find the Running Container
	deployment, err := p.repo.GetActiveDeployment(r.Context(), project.ID)
	if err != nil {
		http.Error(w, "Service not running (No active deployment)", http.StatusBadGateway)
		return
	}

	// 4. The Magic: Reverse Proxy
	// We tell Go: "Forward everything to localhost:32769"
	targetURL := fmt.Sprintf("http://localhost:%d", deployment.Port)
	target, _ := url.Parse(targetURL)

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Update headers so the Flask app knows the real host
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = target.Host

	// Go standard library handles the data stream automatically
	proxy.ServeHTTP(w, r)
}
