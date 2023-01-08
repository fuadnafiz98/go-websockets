package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type messageType string

const (
	WELCOME_MESSAGE messageType = "WELCOME_MESSAGE"
	LEAVE_MESSAGE   messageType = "LEAVE_MESSAGE"
	MESSAGE         messageType = "MESSAGE"
)

type message struct {
	Message     string      `json:"message"`
	MessageType messageType `json:"messageType"`
	Created     time.Time   `json:"created"`
}

type subscriber struct {
	id        string
	username  string
	msgs      chan message
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
	log.Printf("Adding subscrber: %v", sub.username)
	ss.subscriberMutex.Lock()
	ss.subscribers[sub] = struct{}{}
	ss.subscriberMutex.Unlock()
}

func (ss *socketServer) deleteSubscriber(sub *subscriber) {
	log.Printf("Deleting subscrber: %v", sub.username)
	ss.subscriberMutex.Lock()
	delete(ss.subscribers, sub)
	ss.subscriberMutex.Unlock()
	ss.publish([]byte("User Logged out: "+sub.username), LEAVE_MESSAGE)
}

func writeTimeout(ctx context.Context, timeout time.Duration, ws *websocket.Conn, msg message) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	messageData, _ := json.Marshal(msg)
	return ws.Write(ctx, websocket.MessageText, messageData)
}

func (ss *socketServer) susbscribe(ctx context.Context, ws *websocket.Conn) error {
	ctx = ws.CloseRead(ctx)
	sub := &subscriber{
		id:       uuid.New().String(),
		username: genRandomUsername(),
		msgs:     make(chan message, ss.subscriberMessageBuffer),
		closeSlow: func() {
			ws.Close(websocket.StatusPolicyViolation, "Server too slow to handle load")
		},
	}
	ss.addSubscriber(sub)
	ss.publish([]byte("Welcome new user: "+sub.username), WELCOME_MESSAGE)
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

func (ss *socketServer) publish(msg []byte, _messageType messageType) {
	ss.subscriberMutex.Lock()
	defer ss.subscriberMutex.Unlock()

	ss.publishLimiter.Wait(context.Background()) // what is happening here 😫

	messageData := &message{
		Message:     string(msg),
		MessageType: _messageType,
		Created:     time.Now(),
	}

	for s := range ss.subscribers {
		select {
		case s.msgs <- *messageData:
			log.Println("Sending message ...", messageData.Message)
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

	ss.publish(msg, MESSAGE)
	w.WriteHeader(http.StatusAccepted)
}

func (_socketServer *socketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_socketServer.serveMux.ServeHTTP(w, r)
}
