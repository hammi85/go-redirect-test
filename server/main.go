package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/spf13/cobra"
)

// Server is the context of the server
type Server struct {
	bind string
	ctx  *Context
}

type Certificates struct {
	CertFile string
	KeyFile  string
}

// NewServer returns a new server
func NewServer() *Server {
	ctx := NewContext(&State{}, nil)
	return &Server{ctx: ctx}
}

// Run is taking the commands is running the server
func (s *Server) Run(cmd *cobra.Command, args []string) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("www.my-domain.go"), //your domain here
		Cache:      autocert.DirCache("certs"),                 //folder for storing certificates
	}

	// // read configuration
	// configfile := flag.String("config", "config.json", "Configuration file")
	// flag.Parse()

	config := FromFile(*configfile)

	mux := http.NewServeMux()
	RewriteRequest(config, mux)

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: mux,
	}

	err := server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt
	if err != nil {
		log.Println(err)
	}
}

func RewriteRequest(config Configuration, mux *http.ServeMux) {
	for _, domains := range config.Domains {
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.Host == domains.Host {
				http.Redirect(w, req, domains.Target, 301)
			}
		})

	}
}
