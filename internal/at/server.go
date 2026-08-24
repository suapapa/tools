/*
Copyright © 2022 Homin Lee <i@homin.dev>
*/
package at

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A server for test and benchmark network",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: server,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntP("port", "p", 8080, "port to listen")
	serverCmd.Flags().StringP("func", "f", "echo", "server function. [echo|upper|md5]")
}

func server(cmd *cobra.Command, args []string) {
	fmt.Print("server called - ")
	switch protocol {
	case "tcp":
		tcpServer(cmd, args)
	case "udp":
		udpServer(cmd, args)
	default:
		fmt.Println("unknown protocol: ", protocol)
	}
}

func tcpServer(cmd *cobra.Command, args []string) {
	fmt.Println("tcp")
	port := cmd.Flag("port").Value.String()
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	chk(err)
	lstn, err := net.ListenTCP("tcp", addr)
	chk(err)
	defer lstn.Close()

	servFunc, err := cmd.Flags().GetString("func")
	chk(err)

	sem := make(chan struct{}, 256)
	for {
		conn, err := lstn.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		select {
		case sem <- struct{}{}:
		default:
			conn.Close()
			continue
		}
		go func(c *net.TCPConn) {
			defer func() { <-sem }()
			handleTCPConn(c, servFunc)
		}(conn)
	}
}

func udpServer(cmd *cobra.Command, args []string) {
	fmt.Println("udp")
	port := cmd.Flag("port").Value.String()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	chk(err)
	conn, err := net.ListenUDP("udp", addr)
	chk(err)
	defer conn.Close()

	servFunc, err := cmd.Flags().GetString("func")
	chk(err)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n == 0 {
			continue
		}
		log.Printf("udp recv from %s, len=%d", addr.String(), n)
		switch servFunc {
		case "echo":
			_, err = conn.WriteToUDP(buf[:n], addr)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "upper":
			_, err = conn.WriteToUDP([]byte(strings.ToUpper(string(buf[:n]))), addr)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "md5":
			_, err = conn.WriteToUDP([]byte(fmt.Sprintf("%x", md5.Sum(buf[:n]))), addr)
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Println("unknown server function: ", servFunc)
		}
	}
}

func handleTCPConn(conn *net.TCPConn, servFunc string) {
	defer conn.Close()
	idle := 60 * time.Second
	_ = conn.SetDeadline(time.Now().Add(idle))

	buf := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = conn.SetDeadline(time.Now().Add(idle))
		if n == 0 {
			continue
		}
		log.Printf("tcp recv from %s, len=%d", conn.RemoteAddr().String(), n)
		switch servFunc {
		case "echo":
			_, err = conn.Write(buf[:n])
			if err != nil {
				fmt.Println(err)
				return
			}
		case "upper":
			_, err = conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
			if err != nil {
				fmt.Println(err)
				return
			}
		case "md5":
			_, err = conn.Write([]byte(fmt.Sprintf("%x", md5.Sum(buf[:n]))))
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			fmt.Println("unknown server function: ", servFunc)
		}
	}
}
