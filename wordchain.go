package main

import "log"

var Env Environment;

func main() {
    Env, err := LoadEnvironment();
    if err != nil {
        log.Fatalf("Failed to load environment: %s\n", err);
    }

    log.Printf("Environment: %+v", Env);
}
