package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guild_id"`
}

func main() {

	// 读取config.yaml配置
	conf := &Config{}
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 创建会话
	dg, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatalf("创建会话失败: %v", err)
	}

	// 注册回调函数
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			cmd := i.ApplicationCommandData().Name
			// 时间戳 /timestamp
			switch cmd {
			case "timestamp":
				timestamp := time.Now().Unix()
				content := fmt.Sprintf("当前时间戳为: `%d`", timestamp)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				if err != nil {
					log.Printf("响应失败: %v", err)
				}
			case "info":
				// 硬件信息 /info
				percents, _ := cpu.Percent(500*time.Millisecond, false)
				cpuPercent := 0.0
				if len(percents) > 0 {
					cpuPercent = percents[0]
				}

				// CPU
				cores, _ := cpu.Counts(true)
				info, _ := cpu.Info()
				model := "未知"
				if len(info) > 0 {
					model = info[0].ModelName
				}
				cpuContent := fmt.Sprintf("`%s`\n核心: `%d` 核\n负载: `%.2f%%`", model, cores, cpuPercent)

				// RAM
				v, err := mem.VirtualMemory()
				var memContent string
				if err == nil {
					totalGB := float64(v.Total) / 1024 / 1024 / 1024
					usedGB := float64(v.Used) / 1024 / 1024 / 1024
					memContent = fmt.Sprintf("使用率: `%.2f%%` \n内存: `%.1fGB / %.1fGB`", v.UsedPercent, usedGB, totalGB)
				} else {
					memContent = "Error"
				}

				// Embed
				embed := &discordgo.MessageEmbed{
					Title:       "硬件监控",
					Color:       0x00eaac,
					Description: fmt.Sprintf("Bot所在服务器的硬件信息"),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "💻 CPU",
							Value:  cpuContent,
							Inline: false,
						},
						{
							Name:   "📟 RAM",
							Value:  memContent,
							Inline: false,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "🕘",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}

				// Send embed
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
			}

		}
	})

	// 打开连接
	err = dg.Open()
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer dg.Close()

	// 注册命令
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "timestamp",
			Description: "获取当前时间戳",
		},
		{
			Name:        "info",
			Description: "获取硬件信息",
		},
		// add more commands here
	}
	for _, cmd := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, conf.GuildID, cmd)
		if err != nil {
			fmt.Printf("无法创建命令 %s: %w", cmd.Name, err)
		}

	}

	fmt.Println("▶️ Running")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
