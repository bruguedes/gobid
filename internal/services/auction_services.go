package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	PlaceBid MessageKind = iota
)

type MessageKind int

type Message struct {
	Message string
	Kind    MessageKind
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id         uuid.UUID
	Context    context.Context
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID]*Client

	BidsService *BidsService
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserID uuid.UUID
}

func NewAuctionRoom(id uuid.UUID, bidsService *BidsService) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Context:     context.Background(),
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		BidsService: bidsService,
	}
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserID: userID,
	}
}
