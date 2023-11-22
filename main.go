package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
)

type TlsServerConfigIn struct {
	CaCrt     string
	ServerCrt string
	ServerKey string
}

func TlsServerConfig(in TlsServerConfigIn) (*tls.Config, error) {
	wrapErr := func(err error, s string) error {
		return fmt.Errorf("TlsServerConfig: %s: %w", s, err)
	}

	caCertFile, err := os.ReadFile(in.CaCrt)
	if err != nil {
		return nil, wrapErr(err, `os.ReadFile: `+in.CaCrt)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	serverCerts, err := tls.LoadX509KeyPair(in.ServerCrt, in.ServerKey)
	if err != nil {
		return nil, wrapErr(err, `tls.LoadX509KeyPair `+in.ServerCrt+` `+in.ServerKey)
	}

	// Create the TLS Config with the CA pool and enable Client certificate validation
	return &tls.Config{
		ClientCAs:        caCertPool,
		ClientAuth:       tls.RequireAndVerifyClientCert,
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		Certificates: []tls.Certificate{serverCerts},
	}, nil
}

type TlsClientIn struct {
	CaCrt     string
	ClientCrt string
	ClientKey string
}

func TlsClient(in TlsClientIn) (*http.Client, error) {
	wrapErr := func(err error, s string) error {
		return fmt.Errorf("TlsClient: %s: %w", s, err)
	}
	caCertFile, err := os.ReadFile(in.CaCrt)
	if err != nil {
		return nil, wrapErr(err, `os.ReadFile: `+in.CaCrt)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)
	certificate, err := tls.LoadX509KeyPair(in.ClientCrt, in.ClientKey)
	if err != nil {
		return nil, wrapErr(err, `tls.LoadX509KeyPair `+in.ClientCrt+` `+in.ClientKey)
	}

	return &http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}, nil
}

const serverListenPort = `:1443`
const serverAddr = `https://localhost` + serverListenPort

func main() {
	if len(os.Args) == 1 { // start server
		tlsCert := TlsServerConfigIn{
			CaCrt:     "./ca.crt",
			ServerCrt: "./server.crt",
			ServerKey: "./server.key",
		}
		app := fiber.New(fiber.Config{
			Immutable: true,
		})
		app.Use(logger.New())
		app.Use(recover.New())
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(M.SX{
				`hello`: `world`,
			})
		})
		tlsConfig, err := TlsServerConfig(tlsCert)
		L.PanicIf(err, `TlsServerConfig`)
		ln, err := tls.Listen("tcp", serverListenPort, tlsConfig)
		L.PanicIf(err, `tls.Listen`)
		L.PanicIf(app.Listener(ln), `app.Listener`)
	}
	switch os.Args[1] {
	case `client`: // start client
		tlsCert := TlsClientIn{
			CaCrt:     `./ca.crt`,
			ClientCrt: `./client.crt`,
			ClientKey: `./client.key`,
		}
		client, err := TlsClient(tlsCert)
		L.PanicIf(err, `TlsClient`)

		r, err := client.Get(serverAddr)
		L.PanicIf(err, `client.Get `+serverAddr)

		// Read the response body
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)
		body, err := io.ReadAll(r.Body)
		L.PanicIf(err, `io.ReadAll r.Body`)

		// Print the response body to stdout
		fmt.Printf("%s\n", body)
	default:
		L.PanicIf(fmt.Errorf(`unknown command: %s`, os.Args[1]), `main`)
	}
}
