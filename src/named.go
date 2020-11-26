package main

import "fmt"

func main() {
    fmt.Println("Enter your services VM IP:")

    var services string

    fmt.Scanln(&services)

    fmt.Println(services)
}