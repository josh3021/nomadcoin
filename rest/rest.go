package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josh3021/nomadcoin/blockchain"
	"github.com/josh3021/nomadcoin/p2p"
	"github.com/josh3021/nomadcoin/utils"
	"github.com/josh3021/nomadcoin/wallet"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	murl := fmt.Sprintf("http://localhost:%s%s", port, u)
	return []byte(murl), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("Client Connected: http://localhost:%s%s\n", port, r.URL)
		next.ServeHTTP(rw, r)
	})
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	description := []urlDescription{
		{
			URL:         url("/"),
			Method:      http.MethodGet,
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      http.MethodGet,
			Description: "See Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      http.MethodGet,
			Description: "See All Blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      http.MethodPost,
			Description: "Add a Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      http.MethodGet,
			Description: "See a Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      http.MethodGet,
			Description: "See balance of address",
		},
		{
			URL:         url("/mempool"),
			Method:      http.MethodGet,
			Description: "Show Transactions in mempool",
		},
		{
			URL:         url("/wallet"),
			Method:      http.MethodGet,
			Description: "Show my wallet",
		},
		{
			URL:         url("/transactions"),
			Method:      http.MethodPost,
			Description: "Create Transaction",
		},
		{
			URL:         url("ws"),
			Method:      http.MethodGet,
			Description: "WS",
		},
	}
	// jsonBytes, err := json.Marshal(description)
	// utils.HandleErr(err)
	// fmt.Printf("%s", jsonBytes)
	utils.HandleErr(json.NewEncoder(rw).Encode(description))
}

func status(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Status()))
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case http.MethodPost:
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewMessage(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err != nil {
		utils.HandleErr(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	} else {
		utils.HandleErr(encoder.Encode(block))
	}
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func myBalance(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	total := r.URL.Query().Get("total")
	bc := blockchain.Blockchain()
	switch total {
	case "true":
		balance := blockchain.GetBalanceByAddress(bc, address)
		utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{address, balance}))
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(bc, address)))
	}
}

func balance(rw http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]
	total := r.URL.Query().Get("total")
	bc := blockchain.Blockchain()
	switch total {
	case "true":
		balance := blockchain.GetBalanceByAddress(bc, address)
		utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{address, balance}))
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(bc, address)))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	txs := blockchain.MempoolStatus().Txs
	utils.HandleErr(json.NewEncoder(rw).Encode(txs))
}

type myWalletResponse struct {
	Address string `json:"address"`
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.HandleErr(json.NewEncoder(rw).Encode(errorResponse{err.Error()}))
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

type addPeerPayload struct {
	Address string
	Port    string
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	case http.MethodPost:
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port, true)
		rw.WriteHeader(http.StatusCreated)
	}
}

// Start REST API Server
func Start(inputPort int) {
	port = fmt.Sprintf("%d", inputPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods(http.MethodGet)
	router.HandleFunc("/status", status).Methods(http.MethodGet)
	router.HandleFunc("/blocks", blocks).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/blocks/{hash:[0-9a-f]+}", block).Methods(http.MethodGet)
	router.HandleFunc("/balance", myBalance).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", balance).Methods(http.MethodGet)
	router.HandleFunc("/mempool", mempool).Methods(http.MethodGet)
	router.HandleFunc("/wallet", myWallet).Methods(http.MethodGet)
	router.HandleFunc("/transactions", transactions).Methods(http.MethodPost)
	router.HandleFunc("/ws", p2p.Upgrade).Methods(http.MethodGet)
	router.HandleFunc("/peers", peers).Methods(http.MethodGet, http.MethodPost)
	fmt.Printf("???? REST is Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
