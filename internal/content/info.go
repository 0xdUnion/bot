package content

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func Info() *discordgo.MessageEmbed {
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

	return embed
}
