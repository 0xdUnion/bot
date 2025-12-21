package main

import (
	"bot/internal/content"
	"bot/internal/youtube"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Token         string   `yaml:"token"`
	GuildID       string   `yaml:"guild_id"`
	FlowChannelID string   `yaml:"flow_channel_id"`
	YTChannels    []string `yaml:"yt_channels"`
}

func main() {

	// 读取config.yaml配置
	conf := &Config{}
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("[BotSys]❌ 读取配置文件失败: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("[BotSys]❌ 解析配置文件失败: %v", err)
	}

	// 创建会话
	dg, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatalf("[BotSys]❌ 创建会话失败: %v", err)
	}

	// 注册回调函数
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			cmd := i.ApplicationCommandData().Name
			switch cmd {
			case "timestamp":
				// 时间戳 /timestamp

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content.TimeStamp(),
					},
				})
				if err != nil {
					log.Printf("[BotSys]❌ 响应失败: %v", err)
				}
			case "info":
				// 硬件信息 /info

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{content.Info()},
					},
				})
				if err != nil {
					log.Printf("[BotSys]❌ 响应失败: %v", err)
				}
			}

		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("▶️ 启动 YouTube 检查任务...")

		go func() {
			// 间隔10分钟
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()

			// 首次运行
			youtube.RunCheckYouTube(s, conf.FlowChannelID, conf.YTChannels)

			for {
				select {
				case <-ticker.C:
					youtube.RunCheckYouTube(s, conf.FlowChannelID, conf.YTChannels)
				}
			}
		}()
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
