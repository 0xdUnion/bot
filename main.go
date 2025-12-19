package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
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
			// 时间戳 /timestamp
			if i.ApplicationCommandData().Name == "timestamp" {
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
	command := &discordgo.ApplicationCommand{
		Name:        "timestamp",
		Description: "获取当前时间戳",
	}
	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, conf.GuildID, command)
	if err != nil {
		log.Fatalf("无法创建命令: %v", err)
	}

	fmt.Println("▶️ Running")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
