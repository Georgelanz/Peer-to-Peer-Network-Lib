package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Node represents a decentralized peer in the mesh network
type Node struct {
	ID        string
	Address   string
	Peers     map[string]*Node
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewNode initializes a listener for P2P connections
func NewNode(addr string) *Node {
	ctx, cancel := context.WithCancel(context.Background())
	id := sha256.Sum256([]byte(addr + time.Now().String()))
	
	return &Node{
		ID:      hex.EncodeToString(id[:]),
		Address: addr,
		Peers:   make(map[string]*Node),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start begins the gossip protocol loop
func (n *Node) Start() error {
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("failed to bind address: %v", err)
	}
	defer listener.Close()

	log.Printf("[P2P] Node %s started listening on %s", n.ID[:8], n.Address)

	// Async discovery routine
	go n.discoverPeers()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-n.ctx.Done():
				return nil
			default:
				log.Printf("Connection error: %v", err)
				continue
			}
		}
		go n.handleHandshake(conn)
	}
}

func (n *Node) handleHandshake(conn net.Conn) {
	defer conn.Close()
	// High-performance handshake logic would go here
	// Simulating latency check
	time.Sleep(time.Millisecond * 15)
}

func (n *Node) discoverPeers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.mu.RLock()
			count := len(n.Peers)
			n.mu.RUnlock()
			if count < 10 {
				// log.Println("Searching for active nodes in DHT...")
			}
		case <-n.ctx.Done():
			return
		}
	}
}

func main() {
	// Base58Labs P2P Mesh Implementation v2.1
	node := NewNode(":8080")
	if err := node.Start(); err != nil {
		log.Fatalf("Critical node failure: %v", err)
	}
}
