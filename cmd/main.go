package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DeepSeekResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type ChatPageData struct{}

func callDeepSeek(userMessage string) (string, error) {
	payload := map[string]interface{}{
		"model": "deepseek-r1",
		"messages": []map[string]string{
			{"role": "user", "content": userMessage},
		},
		"stream": false,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var dsResp DeepSeekResponse
	if err := json.Unmarshal(body, &dsResp); err != nil {
		return "", err
	}
	return dsResp.Message.Content, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	userMsg := r.FormValue("message")
	if userMsg == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}
	answer, err := callDeepSeek(userMsg)
	if err != nil {
		http.Error(w, "Error calling DeepSeek: "+err.Error(), http.StatusInternalServerError)
		return
	}
	respData := map[string]string{"response": answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, ChatPageData{})
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/chat", chatHandler)
	port := "8080"
	log.Printf("Server started on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
