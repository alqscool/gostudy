package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"sync"
	"strconv"
)

type Block struct {
	Index int			//是这个块在整个链中的位置
	TimeStamp string	//块生成时的时间戳
	BPM int				//每分钟心跳数
	Hash string			//块通过 SHA256 算法生成的散列值
	PreHash string		//代表前一个块的 SHA256 散列值
}

var Blockchain []Block

var mutex = &sync.Mutex{}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.TimeStamp + strconv.Itoa(block.BPM) + block.PreHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block,error) {
	var newBlock Block

	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t.String()
	newBlock.BPM = BPM
	newBlock.PreHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock,oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PreHash{
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlock []Block)  {
	if len(newBlock) > len(Blockchain){
		Blockchain = newBlock
	}
}

func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:			":"+httpPort,
		Handler:		mux,
		ReadTimeout:	10*time.Second,
		WriteTimeout: 	10*time.Second,
		MaxHeaderBytes:	1 << 20,
	}

	if err := s.ListenAndServe(); err!=nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler  {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/",handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/",handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request)  {
	bytes, err := json.MarshalIndent(Blockchain,""," ")
	if err !=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w,string(bytes))
	//log.Println("string(bytes):",string(bytes))
}


type Message struct {
	BPM int
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock,err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
	}
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{})  {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain,genesisBlock)
		mutex.Unlock()
	}()

	log.Fatal(run())
}