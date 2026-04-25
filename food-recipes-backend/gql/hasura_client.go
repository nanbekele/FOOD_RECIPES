package gql

import (
    "net/http"

    "github.com/hasura/go-graphql-client"
)

var Client *graphql.Client

func InitHasuraClient(endpoint string, adminSecret string) {
    Client = graphql.NewClient(endpoint, &http.Client{
        Transport: &adminTransport{adminSecret: adminSecret},
    })
}

type adminTransport struct{ adminSecret string }

func (t *adminTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    if t.adminSecret != "" {
        req.Header.Set("x-hasura-admin-secret", t.adminSecret)
    }
    req.Header.Set("Content-Type", "application/json")
    return http.DefaultTransport.RoundTrip(req)
}

