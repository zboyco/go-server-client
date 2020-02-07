package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	countStr := ""
	for {
		fmt.Print("请输入客户端数量: ")
		if reader.Scan() {
			countStr = reader.Text()
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			log.Fatalln("请输入正确的数量")
		}
		if count <= 0 {
			log.Fatalln("数量错误!")
		}

		var wg sync.WaitGroup

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(clientNo int) {
				defer wg.Done()
				tcpAddr, err := net.ResolveTCPAddr("tcp4", ":9043")
				if err != nil {
					log.Fatalf("Fatal error: %s", err.Error())
				}
				conn, err := net.DialTCP("tcp", nil, tcpAddr)
				if err != nil {
					log.Fatalf("Fatal error: %s", err.Error())
				}
				defer conn.Close()

				var headSize int
				var headBytes = make([]byte, 4)
				headBytes[0] = '$'
				headBytes[3] = '#'

				buffer := make([]byte, 512)
				times := 0
				for {
					times++
					s := fmt.Sprintf("Clinet %v say: %v", clientNo, times)
					content := []byte(s)
					headSize = len(content)
					binary.BigEndian.PutUint16(headBytes[1:], uint16(headSize))
					_, err := conn.Write(headBytes)
					if err != nil {
						log.Println(err)
						return
					}
					_, err = conn.Write(content)
					if err != nil {
						log.Println(err)
						return
					}

					n, err := conn.Read(buffer)
					if err != nil {
						log.Println("Read failed:", err)
						//break
					}

					log.Println(fmt.Sprintf("Clinet [%v] receive : %v", clientNo, string(buffer[:n])))

					time.Sleep(time.Duration(rand.Intn(11)) * time.Second)

					//s = fmt.Sprintf("hello golang - %v", clientNo)
					//content = []byte(s)
					//headSize = len(content)
					//binary.BigEndian.PutUint16(headBytes[1:], uint16(headSize))
					//conn.Write(headBytes)
					//conn.Write(content)
					//time.Sleep(time.Duration(rand.Intn(12)) * time.Second)
					//
					//s = fmt.Sprintf("hello socket - %v", clientNo)
					//content = []byte(s)
					//headSize = len(content)
					//binary.BigEndian.PutUint16(headBytes[1:], uint16(headSize))
					//conn.Write(headBytes)
					//conn.Write(content)
					//time.Sleep(time.Duration(rand.Intn(12)) * time.Second)
				}
			}(i + 1)
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		}
		wg.Wait()
	}
}