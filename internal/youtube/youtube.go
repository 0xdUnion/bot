package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

const statePath = "./data/youtube.json"

func update(channelID string, videoID string) {
	// 确保./data存在
	_ = os.MkdirAll("./data", 0755)

	data := map[string]string{}
	if b, err := os.ReadFile(statePath); err == nil {
		_ = json.Unmarshal(b, &data)
	}

	data[channelID] = videoID

	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(statePath, b, 0644)
}

func get(channelID string) string {

	data := map[string]string{}
	if b, err := os.ReadFile(statePath); err == nil {
		_ = json.Unmarshal(b, &data)
	}

	return data[channelID]

}

func fetchLatestVideo(channelID string) (string, string, error) {
	url := fmt.Sprintf(
		"https://www.youtube.com/feeds/videos.xml?channel_id=%s",
		channelID,
	)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return "", "", err
	}

	channelName := feed.Title

	if len(feed.Items) == 0 {
		return "", channelName, fmt.Errorf("empty feed")
	}

	fullID := feed.Items[0].GUID
	parts := strings.Split(fullID, ":")
	return parts[len(parts)-1], channelName, nil
}

func RunCheckYouTube(s *discordgo.Session, discordChannelID string, ytChannels []string) {
	for _, ch := range ytChannels {
		latestVideoID, channelName, err := fetchLatestVideo(ch)
		if err != nil {
			log.Printf("[YouTube] 获取频道 %s 失败: %v", ch, err)
			continue
		}

		// 获取上次记录的ID
		oldID := get(ch)

		// 如果ID一样则跳过
		if oldID == latestVideoID {
			continue
		}

		// 更新
		update(ch, latestVideoID)

		// ID为空->首次运行或新添加的频道
		if oldID == "" {
			log.Printf("[YouTube] 频道 %s 初始记录已更新，跳过首条推送", ch)
			continue
		}

		// 发送
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", latestVideoID)
		message := fmt.Sprintf(
			"🌟**有新视频发布！**\n\n"+
				"**频道**：%s\n"+
				"🔗 %s\n",
			channelName, videoURL,
		)
		_, err = s.ChannelMessageSend(discordChannelID, message)
		if err != nil {
			log.Printf("[BotSys] 消息发送失败 (%s): %v", ch, err)
		} else {
			log.Printf("[YouTube] 推送成功: %s", videoURL)
		}

		time.Sleep(5 * time.Second)
	}
}
