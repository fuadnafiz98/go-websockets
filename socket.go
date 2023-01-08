package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

type socketServer struct {
	subscriberMessageBuffer int
	publishLimiter          *rate.Limiter
	logf                    func(f string, v ...interface{}) // Don't know what this means.
	serveMux                http.ServeMux
	subscriberMutex         sync.Mutex
	subscribers             map[*subscriber]struct{}
}

func newSocketServer() *socketServer {
	_socketServer := &socketServer{
		subscriberMessageBuffer: 16,
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
	}
	_socketServer.serveMux.Handle("/", http.FileServer(http.Dir("./static")))

	_socketServer.serveMux.HandleFunc("/subscribe", _socketServer.susbscribeHandler)
	_socketServer.serveMux.HandleFunc("/publish", _socketServer.publishHandler)
	return _socketServer
}

func (ss *socketServer) addSubscriber(sub *subscriber) {
	log.Println("Adding subscrber!")
	ss.subscriberMutex.Lock()
	ss.subscribers[sub] = struct{}{}
	ss.subscriberMutex.Unlock()
}

func (ss *socketServer) deleteSubscriber(sub *subscriber) {
	log.Println("Deleting subscrber!")
	ss.subscriberMutex.Lock()
	delete(ss.subscribers, sub)
	ss.subscriberMutex.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, ws *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ws.Write(ctx, websocket.MessageText, msg)
}

func (ss *socketServer) susbscribe(ctx context.Context, ws *websocket.Conn) error {
	ctx = ws.CloseRead(ctx)
	sub := &subscriber{
		msgs: make(chan []byte, ss.subscriberMessageBuffer),
		closeSlow: func() {
			ws.Close(websocket.StatusPolicyViolation, "Server too slow to handle load")
		},
	}
	ss.addSubscriber(sub)
	defer ss.deleteSubscriber(sub)

	for {
		select {
		case msg := <-sub.msgs:
			err := writeTimeout(ctx, time.Second*5, ws, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ss *socketServer) susbscribeHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		ss.logf("%v", err)
		return
	}
	defer ws.Close(websocket.StatusInternalError, "Closed in defer")
	err = ss.susbscribe(r.Context(), ws)
	if errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		ss.logf("%v", err)
		return
	}
}

func (ss *socketServer) publish(msg []byte) {
	ss.subscriberMutex.Lock()
	defer ss.subscriberMutex.Unlock()

	ss.publishLimiter.Wait(context.Background()) // what is happening here 😫

	for s := range ss.subscribers {
		select {
		case s.msgs <- msg:
			log.Println("Sending message ...", string(msg))
		default:
			go s.closeSlow()
		}
	}
}

func (ss *socketServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	ss.publish(msg)
	w.WriteHeader(http.StatusAccepted)
}

func (_socketServer *socketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_socketServer.serveMux.ServeHTTP(w, r)
}
