package main

import (
  "fmt"
  "time"
)

func worker(die chan bool) {
    fmt.Println("Begin: This is Worker")
    for {
        select {
        //case xx：
        //做事的分支
        case <-die:
            fmt.Println("Done: This is Worker")
            die <- true
          time.Sleep( time.Second * 5 )
            return
        }
    }
}

func main() {
  die := make(chan bool, 5)

  for i := 1; i <= 1000; i++ {
    go worker(die)
  }

  die <- true
  <-die
  fmt.Println("Worker goroutine has been terminated")
}
