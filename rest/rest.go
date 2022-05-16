package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/josh3021/nomadcoin/blockchain"
	"github.com/josh3021/nomadcoin/utils"
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

func documentation(rw http.ResponseWriter, r *http.Request) {
	description := []urlDescription{
		{
			URL:         url("/"),
			Method:      http.MethodGet,
			Description: "See Documentation",
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
			URL:         url("/blocks/{id}"),
			Method:      http.MethodGet,
			Description: "See a Block",
		},
	}
	// jsonBytes, err := json.Marshal(description)
	// utils.HandleErr(err)
	// fmt.Printf("%s", jsonBytes)
	json.NewEncoder(rw).Encode(description)
}

type addBlockBody struct {
	Data string
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(rw).Encode(blockchain.BlockChain().GetAllBlocks())
	case http.MethodPost:
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.BlockChain().AddBlock(addBlockBody.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.BlockChain().GetBlock(height)
	encoder := json.NewEncoder(rw)
	if err != nil {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(port int) {
	strPort = fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods(http.MethodGet)
	router.HandleFunc("/blocks", blocks).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/blocks/{height:[0-9+]}", block).Methods(http.MethodGet)
	fmt.Printf("📃 REST is Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(strPort, router))
}
