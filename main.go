package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := "your discord token here"

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error membuat sesi Discord:", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error membuka koneksi:", err)
		return
	}

	fmt.Println("Bot Discord sedang berjalan. Tekan CTRL-C untuk berhenti.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, s.State.User.ID) && strings.Contains(strings.ToLower(m.Content), "request file") {

		parts := strings.Split(strings.ToLower(m.Content), "request file")
		if len(parts) > 1 {
			fileName := strings.TrimSpace(parts[1])

			file, err := os.Open(fileName)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "File tidak ditemukan: "+fileName)
				return
			}
			defer file.Close()

			_, err = s.ChannelFileSend(m.ChannelID, fileName, file)
			if err != nil {
				fmt.Println("Error kirim file:", err)
				s.ChannelMessageSend(m.ChannelID, "Gagal mengirim file.")
			}
		}
	}
}
