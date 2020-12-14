package main

import (
    "github.com/joho/godotenv"
    "github.com/julienschmidt/httprouter"
    "log"
    "os"
    "fmt"
    "net/http"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

}
