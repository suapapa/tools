/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A client for test and benchmark network",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: client,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clientCmd.Flags().StringP("addr", "a", "127.0.0.1:8080", "server address to connect")
	clientCmd.Flags().BoolP("interactive", "i", false, "interactive mode")
}

func client(cmd *cobra.Command, args []string) {
	fmt.Println("client called")
	switch protocol {
	case "tcp", "udp":
		addr := cmd.Flag("addr").Value.String()
		conn, err := net.Dial(protocol, addr)
		chk(err)
		interactiveStr := cmd.Flag("interactive").Value.String()
		var interactive bool
		if interactiveStr == "true" {
			interactive = true
		}
		handleClient(conn, interactive)
	default:
		fmt.Println("unknown protocol: ", protocol)
	}
}

func handleClient(conn net.Conn, interactive bool) {
	defer conn.Close()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			fmt.Printf("from %s(len=%d): %s\n", conn.RemoteAddr(), n, string(buf[:n]))
		}
	}()

	if interactive {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">> ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text == "quit" {
				break
			}
			fmt.Fprintf(conn, "%s", text)
		}
	} else {
		for {
			fmt.Fprintf(conn, "hello")
			time.Sleep(1 * time.Second)
		}
	}
}
