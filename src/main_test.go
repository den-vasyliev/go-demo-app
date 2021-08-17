package main

import (
    "testing"
)


func TestGetEnv(t *testing.T) {
    getEnv("APP_DB_NAME", "demo")

}


func TestInitOptions(t *testing.T){

    initOptions()

}


