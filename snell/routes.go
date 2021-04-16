package main

import (
        "encoding/base64"
        "encoding/hex"
        "fmt"
        "github.com/ellemouton/thunder/elle"
        "github.com/lightningnetwork/lnd/lntypes"
        "net/http"
        "regexp"
        "strings"

        "github.com/prometheus/common/log"
        "github.com/skip2/go-qrcode"
)

func startRouter(s *State) {
        http.HandleFunc("/", s.fwdHandler())

        log.Infof("serving on :9001")
        err := http.ListenAndServe(":9001", nil)
        if err != nil {
                panic(err)
        }
}

var authRegex = regexp.MustCompile("LSAT (.*?):([a-f0-9]{64})")

func (s *State) fwdHandler() func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                // Check if a payment is even required for this path.
                resp, err := s.elleClient.RequiresPayment(r.Context(), r.URL.Path)
                if err != nil {
                       fmt.Fprint(w, err)
                       return
                }

                if !resp.RequiresPayment {
                        s.proxy.ServeHTTP(w, r)
                        return
                }

                // Cool, payment is required. Before generating an invoice, first
                // check if the request has valid auth.
                // Check header for authorization and redirect to payment handler if needed.
                auth := r.Header.Get("Authorization")

                if auth == "" {
                        s.paymentHandler(w, r, resp)
                        return
                }

                if !authRegex.MatchString(auth) {
                        s.paymentHandler(w, r, resp)
                        return
                }

                matches := authRegex.FindStringSubmatch(auth)
                if len(matches) != 3 {
                        s.paymentHandler(w, r, resp)
                        return
                }

                macString, preimageHexString := matches[1], matches[2]

                macBytes, err := base64.StdEncoding.DecodeString(macString)
                if err != nil {
                        s.paymentHandler(w, r, resp)
                        return
                }

                preimage, err := hex.DecodeString(preimageHexString)
                if err != nil {
                        s.paymentHandler(w, r, resp)
                        return
                }

                valid, err := s.macClient.Verify(r.Context(), macBytes, preimage, "blog", resp.ID)
                if err != nil {
                        s.paymentHandler(w, r, resp)
                        return
                }

                if !valid {
                        s.paymentHandler(w, r, resp)
                        return
                }

                s.proxy.ServeHTTP(w, r)
        }
}

func (s * State) paymentHandler(w http.ResponseWriter, r *http.Request, details *elle.AssetDetails) {
        invoice, err := s.lndClient.AddInvoice(r.Context(), details.Price, 3600, "ellemouton.com:" +r.URL.Path)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        // construct QR code of the invoice
        png, err := qrcode.Encode(strings.ToUpper(invoice.PaymentRequest), qrcode.Medium, 256)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        encodedPngString := base64.StdEncoding.EncodeToString(png)

        var payHash lntypes.Hash
        copy(payHash[:], invoice.RHash)

        // Bake a new macaroon
        mac, err := s.macClient.Create(r.Context(), payHash, "blog", details.ID)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        macBytes, err := mac.MarshalBinary()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        macString := base64.StdEncoding.EncodeToString(macBytes)

        // Add the partial LSAT (mac + invoice) to the response header
        str := fmt.Sprintf("LSAT macaroon=\"%s\", invoice=\"%s\"", macString, invoice.PaymentRequest)
        w.Header().Set("WWW-Authenticate", str)

        c := struct {
                ID int64
                Price int64
                QrCode     string
                Invoice string
                Macaroon string
        }{
                ID: details.ID,
                Price: details.Price,
                QrCode: encodedPngString,
                Invoice: invoice.PaymentRequest,
                Macaroon: macString,
        }
        fmt.Println(c.ID)
        err = templates.ExecuteTemplate(w, "payment.html", c)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
}