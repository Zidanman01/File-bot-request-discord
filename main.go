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
	// Ganti dengan Token Bot Discord Anda
	token := "MTQ2MTc0MTg3MDkyODM2NzY2OQ.G7GPnh.fWlDW30cu-V5Yj33CHXlSwMrN6Ifs-MG7N46aU"

	// Membuat sesi Discord baru
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error membuat sesi Discord:", err)
		return
	}

	// Menambahkan handler untuk menangkap pesan masuk
	dg.AddHandler(messageCreate)

	// Menentukan intent agar bisa membaca konten pesan
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Membuka koneksi ke Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error membuka koneksi:", err)
		return
	}

	fmt.Println("Bot Discord sedang berjalan. Tekan CTRL-C untuk berhenti.")

	// Menjaga program tetap jalan sampai ada sinyal stop
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Abaikan pesan dari bot itu sendiri
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Cek apakah bot di-mention dan mengandung kata "request file"
	// Di Discord, mention berbentuk <@!ID_BOT>
	if strings.Contains(m.Content, s.State.User.ID) && strings.Contains(strings.ToLower(m.Content), "request file") {

		// Logika mengambil nama file
		parts := strings.Split(strings.ToLower(m.Content), "request file")
		if len(parts) > 1 {
			fileName := strings.TrimSpace(parts[1])

			// Membuka file lokal
			file, err := os.Open(fileName)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "❌ File tidak ditemukan: "+fileName)
				return
			}
			defer file.Close()

			// Mengirim file ke Discord
			_, err = s.ChannelFileSend(m.ChannelID, fileName, file)
			if err != nil {
				fmt.Println("Error kirim file:", err)
				s.ChannelMessageSend(m.ChannelID, "❌ Gagal mengirim file.")
			}
		}
	}
}
