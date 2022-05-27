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

var strPort string

type url string

func (u url) MarshalText() ([]byte, error) {
	murl := fmt.Sprintf("http://localhost%s%s", strPort, u)
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
		fmt.Println(r.URL)
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
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain()))
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case http.MethodPost:
		blockchain.Blockchain().AddBlock()
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
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
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
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.HandleErr(json.NewEncoder(rw).Encode(errorResponse{err.Error()}))
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// Start REST API Server
func Start(port int) {
	strPort = fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods(http.MethodGet)
	router.HandleFunc("/status", status).Methods(http.MethodGet)
	router.HandleFunc("/blocks", blocks).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/blocks/{hash:[0-9a-f]+}", block).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", balance).Methods(http.MethodGet)
	router.HandleFunc("/mempool", mempool).Methods(http.MethodGet)
	router.HandleFunc("/wallet", myWallet).Methods(http.MethodGet)
	router.HandleFunc("/transactions", transactions).Methods(http.MethodPost)
	router.HandleFunc("/ws", p2p.Upgrade).Methods(http.MethodGet)
	fmt.Printf("📃 REST is Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(strPort, router))
}
