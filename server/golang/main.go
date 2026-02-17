package main

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    _ "github.com/microsoft/go-mssqldb"
)

type Suggestion struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Category string `json:"category"`
    Message  string `json:"message"`
    Date     string `json:"date"`
    AIReply  string `json:"aiReply"`
}

type WebhookPayload struct {
    ID      int    `json:"id"`
    Message string `json:"message"`
    Email   string `json:"email"`
}

var db *sql.DB
var webhookURL string

func initDB() {
    err := godotenv.Load("/root/.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    connectionString := os.Getenv("DATABASE_URL")
    if connectionString == "" {
        log.Fatal("DATABASE_URL not set in environment")
    }

    webhookURL = os.Getenv("WEBHOOK_URL") // Optional webhook URL

    db, err = sql.Open("sqlserver", connectionString)
    if err != nil {
        log.Fatal("Error creating connection pool: ", err.Error())
    }

    // Create table if it doesn't exist
    createTableSQL := `
    IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='suggestions' AND xtype='U')
    CREATE TABLE suggestions (
        id INT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(100),
        email NVARCHAR(100),
        category NVARCHAR(50),
        message NVARCHAR(MAX),
        date DATETIME2
    );`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal("Error creating table: ", err.Error())
    }

    fmt.Println("Database connected successfully!")
}

func getSuggestions(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT TOP 5 id, name, email, category, message, aiReply, date FROM suggestions ORDER BY date DESC")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var suggestions []Suggestion
    for rows.Next() {
        var s Suggestion
        var date time.Time
        err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Category, &s.Message, &s.AIReply, &date)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        s.Date = date.Format(time.RFC3339)
        suggestions = append(suggestions, s)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(suggestions)
}

func callWebhook(id int, message string, email string) {
    if webhookURL == "" {
        log.Println("WEBHOOK_URL not configured, skipping webhook call")
        return
    }

    // Prepare webhook payload
    payload := WebhookPayload{
        ID:      id,
        Message: message,
	Email:	 email,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshaling webhook payload: %v", err)
        return
    }

    // Call the webhook
    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        log.Printf("Error calling webhook: %v", err)
        return
    }
    defer resp.Body.Close()

    log.Printf("Webhook called successfully for suggestion ID: %d", id)
}

func createSuggestion(w http.ResponseWriter, r *http.Request) {
    var s Suggestion
    err := json.NewDecoder(r.Body).Decode(&s)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Use SQL Server's GETDATE() function instead of passing date from client
    query := `INSERT INTO suggestions (name, email, category, message, aiReply, date) OUTPUT INSERTED.id VALUES (@p1, @p2, @p3, @p4, '', GETDATE())`
    
    var insertedID int
    err = db.QueryRow(query, s.Name, s.Email, s.Category, s.Message).Scan(&insertedID)
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch the complete inserted record
    row := db.QueryRow("SELECT id, name, email, category, message, date FROM suggestions WHERE id = @p1", insertedID)
    var date time.Time
    err = row.Scan(&s.ID, &s.Name, &s.Email, &s.Category, &s.Message, &date)
    if err != nil {
        http.Error(w, "Error fetching inserted record: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    s.ID = insertedID
    s.Date = date.Format(time.RFC3339)

    // Call webhook with the new suggestion
    go callWebhook(s.ID, s.Message, s.Email)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(s)


    /*
    s.Date = time.Now().Format(time.RFC3339)

    query := `INSERT INTO suggestions (name, email, category, message, date) VALUES (@p1, @p2, @p3, @p4, @p5)`
    _, err = db.Exec(query, s.Name, s.Email, s.Category, s.Message, s.Date)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(s)
    */
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func main() {
    initDB()
    defer db.Close()

    r := mux.NewRouter()
    
    r.HandleFunc("/api/suggestions", enableCORS(getSuggestions)).Methods("GET", "OPTIONS")
    r.HandleFunc("/api/suggestions", enableCORS(createSuggestion)).Methods("POST", "OPTIONS")
    
    // Serve static files
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
