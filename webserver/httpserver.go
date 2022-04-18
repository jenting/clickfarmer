// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package webserver

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"storj.io/clickfarmer/pb"
	"storj.io/common/sync2"
)

func Run(ctx context.Context, httpAddr, rpcAddr, webdir string, cacheInterval time.Duration) error {
	apiServer, err := NewAPIServer(httpAddr, rpcAddr, webdir, cacheInterval)
	if err != nil {
		return err
	}
	defer apiServer.Close()

	var group errgroup.Group
	group.Go(func() error {
		return apiServer.ServeHTTP(ctx)
	})
	group.Go(func() error {
		return apiServer.cache.loop.Run(ctx, apiServer.refreshCache)
	})

	return group.Wait()
}

type Cache struct {
	mu     sync.Mutex
	loop   *sync2.Cycle
	values JSONClicks
}

type APIServer struct {
	grpcConn    *grpc.ClientConn
	clickFarmer pb.ClickFarmerClient
	webdir      string
	httpAddr    string
	cache       Cache
}

func NewAPIServer(httpAddr, rpcAddr, webdir string, cacheInterval time.Duration) (*APIServer, error) {
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	// make a grpc proto-specific client
	client := pb.NewClickFarmerClient(conn)

	return &APIServer{
		grpcConn:    conn,
		clickFarmer: client,
		webdir:      webdir,
		httpAddr:    httpAddr,
		cache: Cache{
			loop: sync2.NewCycle(cacheInterval),
		},
	}, nil
}

func (a *APIServer) Close() error {
	return a.grpcConn.Close()
}

func (a *APIServer) refreshCache(ctx context.Context) error {
	fmt.Println("refreshing cache")

	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	_, err := a.clickFarmer.SetClicks(ctx, &pb.SetClicksRequest{
		ClickCounts: &pb.ClickCounts{
			Red:   a.cache.values.Red,
			Green: a.cache.values.Green,
			Blue:  a.cache.values.Blue,
		},
	})
	if err != nil {
		return err
	}

	getRes, err := a.clickFarmer.GetClicks(ctx, &pb.GetClicksRequest{})
	if err != nil {
		return err
	}
	a.cache.values.Red = getRes.ClickCounts.Red
	a.cache.values.Green = getRes.ClickCounts.Green
	a.cache.values.Blue = getRes.ClickCounts.Blue

	return nil
}

// Reminder: we rely on GET /api/clicks to return JSON like
// {redClicks: 0, greenClicks: 1, blueClicks: 0} for testing your submission.
// If you change the below code, make sure the API continues to work. See
// more in the assignment README.md.
type JSONClicks struct {
	Red   int64 `json:"redClicks"`
	Green int64 `json:"greenClicks"`
	Blue  int64 `json:"blueClicks"`
}

func (a *APIServer) ServeHTTP(ctx context.Context) error {
	// Reminder: we rely on the following API endpoints for testing your
	// submission. If you change the below code, make sure the API continues
	// to work. See more in the assignment README.md.
	http.HandleFunc("/api/clicks", func(w http.ResponseWriter, r *http.Request) {
		a.getClicksHandler(ctx, w, r)
	})
	http.HandleFunc("/api/clicks/red", func(w http.ResponseWriter, r *http.Request) {
		a.clickColorHandler(ctx, w, r, "red")
	})
	http.HandleFunc("/api/clicks/green", func(w http.ResponseWriter, r *http.Request) {
		a.clickColorHandler(ctx, w, r, "green")
	})
	http.HandleFunc("/api/clicks/blue", func(w http.ResponseWriter, r *http.Request) {
		a.clickColorHandler(ctx, w, r, "blue")
	})

	// static files
	fs := http.FileServer(http.Dir(filepath.Join(a.webdir, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// home page
	http.HandleFunc("/", a.indexHandler)

	return http.ListenAndServe(a.httpAddr, nil)
}
