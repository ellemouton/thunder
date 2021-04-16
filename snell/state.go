package main

import (
        "fmt"
        "net/http"
        "net/http/httputil"
        "net/url"

        "github.com/ellemouton/thunder/elle"
        elle_client "github.com/ellemouton/thunder/elle/client"
        "github.com/ellemouton/thunder/lnd"
        "github.com/ellemouton/thunder/macaroon"
)

type State struct {
        lndClient lnd.Client
        macClient macaroon.Client
        proxy *httputil.ReverseProxy
        elleClient elle.Client
}

func newState() (*State, error) {
        s := new(State)

        mc, err := macaroon.New()
        if err != nil {
                return nil, fmt.Errorf("problem creating macaroon client: %s", err)
        }
        s.macClient = mc

        lc, err := lnd.New()
        if err != nil {
                return nil, fmt.Errorf("problem creating lnd client: %s", err)
        }
        s.lndClient = lc

        ec, err := elle_client.New()
        if err != nil {
                return nil, fmt.Errorf("problem creating elle client: %s", err)
        }
        s.elleClient = ec

        serverAddr, _ := url.Parse(*forwardAddr)
        director := func(req *http.Request) {
                req.Header.Add("X-Forwarded-Host", req.Host)
                req.Header.Add("X-Origin-Host", serverAddr.Host)
                req.URL.Scheme = "http"
                req.URL.Host = serverAddr.Host
        }
        s.proxy = &httputil.ReverseProxy{Director: director}

        return s, nil
}

func (s *State) GetMacaroonClient() macaroon.Client {
        return s.macClient
}

func (s *State) GetLndClient() lnd.Client {
        return s.lndClient
}

func (s *State) GetElleClient() elle.Client {
        return s.elleClient
}

func (s *State) cleanup() {
        if err := s.macClient.Close(); err != nil {
                fmt.Errorf("error closing mac client: %v", err)
        }

        if err := s.lndClient.Close(); err != nil {
                fmt.Errorf("error closing lnd client: %v", err)
        }

        if err := s.elleClient.Close(); err != nil {
                fmt.Errorf("error closing elle client: %v", err)
        }
}

