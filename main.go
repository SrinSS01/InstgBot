package main

import (
	"InstgBot/config"
	"InstgBot/session"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Info struct {
	Session  string
	Username string
	Name     string
	Bio      string
	URL      string
	Email    string
	Phone    string
	Timer    *time.Timer
	Question int
}

var (
	discord    *discordgo.Session
	cnfg       = config.Config{}
	sessionObj = session.Session{
		Sessions: []string{},
	}
	accRegex, _ = regexp.Compile("^-acc(?: +(?P<sessionid>.+))?$")
	InfoMap     = map[string]*Info{}
)

func init() {
	file, err := os.ReadFile("session.json")
	if err != nil {
		return
	}
	if err := json.Unmarshal(file, &sessionObj); err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Print("Enter bot token: ")
		if _, err := fmt.Scanln(&cnfg.Token); err != nil {
			log.Fatal("Error during Scanln(): ", err)
		}
		configJson()
		return
	}
	if err := json.Unmarshal(file, &cnfg); err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func configJson() {
	marshal, err := json.Marshal(&cnfg)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
		return
	}
	if err := os.WriteFile("config.json", marshal, 0644); err != nil {
		log.Fatal("Error during WriteFile(): ", err)
	}
}

func saveSession() {
	marshal, err := json.Marshal(&sessionObj)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
		return
	}
	if err := os.WriteFile("session.json", marshal, 0644); err != nil {
		log.Fatal("Error during WriteFile(): ", err)
	}
	return
}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Println(s.State.User.Username + " is ready")
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	authorId := m.Author.ID
	if authorId == s.State.User.ID {
		return
	}
	matches := accRegex.FindStringSubmatch(m.Content)
	if len(matches) == 0 {
		return
	}
	msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Processing...", m.Reference())
	sessionId := matches[accRegex.SubexpIndex("sessionid")]
	if sessionId == "" {
		sessionId = sessionObj.Sessions[0]
		sessionObj.Sessions = sessionObj.Sessions[1:]
	}
	response, err := getInfoFromSession(sessionId)
	username := response["username"]
	if err != nil {
		_, _ = s.ChannelMessageSendReply(msg.ChannelID, err.Error(), msg.Reference())
		return
	}
	key := authorId + m.GuildID
	info := InfoMap[key]
	if info == nil {
		InfoMap[key] = &Info{
			Question: 0,
		}
	}
	info = InfoMap[key]
	info.Question = 1
	info.Session = sessionId
	info.Email = response["email"].(string)
	info.Phone = response["phone_number"].(string)
	msg, _ = s.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("The current username @%s, the new username should be?", username), msg.Reference())
	info.Timer = time.AfterFunc(60*time.Second, func() {
		delete(InfoMap, key)
		_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
	})
}

func questionAnswerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	authorId := m.Author.ID
	if authorId == s.State.User.ID {
		return
	}
	key := authorId + m.GuildID
	info := InfoMap[key]
	if info == nil {
		return
	}
	content := m.Content
	switch info.Question {
	case 1:
		if matched, _ := regexp.MatchString("^[A-Za-z]\\w{5,29}$", content); !matched {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please enter a valid username", m.Reference())
			return
		}
		info.Username = content
		defer info.Timer.Stop()
		info.Question = 2
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Provide the Name?", m.Reference())
		info.Timer = time.AfterFunc(60*time.Second, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 2:
		if content != "-" {
			info.Name = content
		}
		defer info.Timer.Stop()
		info.Question = 3
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Could you provide your bio information?", m.Reference())
		info.Timer = time.AfterFunc(60*time.Second, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 3:
		if content != "-" {
			info.Bio = content
		}
		defer info.Timer.Stop()
		info.Question = 4
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "What’s the URL you would like to use?", m.Reference())
		info.Timer = time.AfterFunc(60*time.Second, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 4:
		if content != "-" {
			if res, _ := regexp.MatchString("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)", content); !res {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please enter a valid url", m.Reference())
				return
			}
			info.URL = content
		}
		defer info.Timer.Stop()
		info.Question = 5
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Upload the Posts you want to add.", m.Reference())
		info.Timer = time.AfterFunc(60*time.Second, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 5:
		attachments := m.Attachments
		if len(attachments) == 0 {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please provide an image to post", m.Reference())
			return
		}
		defer info.Timer.Stop()
		err := postChanges(info)
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Failed to upload changes:\n%s", err.Error()), m.Reference())
			return
		}

		msg := m.Message

		for _, attachment := range attachments {
			err := postMedia(attachment, info.Session)
			if err != nil {
				msg, _ = s.ChannelMessageSendReply(msg.ChannelID, err.Error(), msg.Reference())
			}
		}
		_, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Your Account is ready to be used.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Name",
					Value: info.Name,
				},
				{
					Name:  "New Username",
					Value: "@" + info.Username,
				},
				{
					Name:  "Bio",
					Value: info.Bio,
				},
				{
					Name:  "Url",
					Value: info.URL,
				},
			},
		}, m.Reference())
		delete(InfoMap, key)
	}
}

func postMedia(attachment *discordgo.MessageAttachment, sessionId string) error {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	request := resty.New().R()
	response, err := request.Get(attachment.URL)
	if err != nil {
		return err
	}
	photo := response.Body()
	req, err := http.NewRequest("POST", "https://i.instagram.com/rupload_igphoto/fb_uploader_"+timestamp, bytes.NewBuffer(photo))
	if err != nil {
		return err
	}

	req.Header.Set("Host", "i.instagram.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Entity-Type", "image/jpeg")
	req.Header.Set("X-ASBD-ID", "198387")
	req.Header.Set("X-IG-App-ID", "936619743392459")
	req.Header.Set("X-Instagram-AJAX", "1007055396")
	req.Header.Set("X-Entity-Length", strconv.Itoa(len(photo)))
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Entity-Name", "fb_uploader_"+timestamp)
	req.Header.Set("Origin", "https://www.instagram.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")
	req.Header.Set("Referer", "https://www.instagram.com/")
	req.Header.Set("Offset", "0")
	req.Header.Set("X-Instagram-Rupload-Params", fmt.Sprintf(`{"media_type":1,"upload_id":"%s","upload_media_height":318,"upload_media_width":318}`, timestamp))
	req.AddCookie(&http.Cookie{Name: "sessionid", Value: sessionId})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	formData := url.Values{
		"source_type":                   {"library"},
		"caption":                       {""},
		"upload_id":                     {timestamp},
		"disable_comments":              {"0"},
		"like_and_view_counts_disabled": {"0"},
		"igtv_share_preview_to_feed":    {"1"},
		"is_unified_video":              {"1"},
		"video_subtitles_enabled":       {"0"},
		"disable_oa_reuse":              {"false"},
	}

	reqs, err := http.NewRequest("POST", "https://www.instagram.com/api/v1/media/configure/", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		panic(err)
	}

	reqs.Header.Set("Host", "www.instagram.com")
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("X-ASBD-ID", "198387")
	reqs.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqs.Header.Set("X-IG-App-ID", "936619743392459")
	reqs.Header.Set("X-Instagram-AJAX", "1007055396")
	reqs.Header.Set("Accept-Language", "en-US,en;q=0.9")
	reqs.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqs.Header.Set("Origin", "https://www.instagram.com")
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")
	reqs.Header.Set("Referer", "https://www.instagram.com/y7zpsjdisjsji/")
	reqs.Header.Set("X-IG-WWW-Claim", "hmac.AR3aUUlZDRH4FHt_roaehflyECa0AZ--MzLh1brm-OcqvpSn")
	reqs.Header.Set("Content-Length", fmt.Sprintf("%d", len(formData.Encode())))
	reqs.Header.Set("Connection", "keep-alive")
	reqs.Header.Set("X-CSRFToken", "6aEqkfPKzcgL5kIZQL6mzYFc86dWtTHs")
	reqs.AddCookie(&http.Cookie{Name: "sessionid", Value: sessionId})

	clients := &http.Client{}
	resps, err := clients.Do(reqs)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	b, _ := io.ReadAll(resps.Body)
	if strings.Contains(string(b), "\"status\":\"ok\"") {
		fmt.Println("Successfully Upload Photo")
		return nil
	} else {
		return fmt.Errorf("erro uploading photo")
	}
}

func postChanges(info *Info) error {
	apiUrl := "https://i.instagram.com/api/v1/accounts/edit_profile/"
	request := resty.New().R()
	request.SetHeaders(map[string]string{
		"Host":              "i.instagram.com",
		"Cookie":            "sessionid=" + info.Session,
		"X-Ig-Capabilities": "3brTvw==",
		"User-Agent":        "Instagram 103.1.0.15.119 Android (25/7.1.2; 240dpi; 720x1280; samsung; SM-G988N; z3q; exynos8895; en_US; 164094540)",
		"Content-Type":      "application/x-www-form-urlencoded; charset=UTF-8",
	})
	randomUUID, err := uuid.NewRandom()
	randomUUID2, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	request.SetFormData(map[string]string{
		"external_url": info.URL,
		"phone_number": info.Phone,
		"username":     info.Username,
		"first_name":   info.Name,
		"_csrftoken":   "missing",
		"_uid":         "48403763438",
		"device_id":    randomUUID.String(),
		"_uuid":        randomUUID2.String(),
		"biography":    info.Bio,
		"email":        info.Email,
	})
	response, err := request.Post(apiUrl)
	if err != nil {
		return err
	}
	var JSON map[string]interface{}
	body := response.Body()
	err = json.Unmarshal(body, &JSON)
	if err != nil {
		return err
	}
	if status, ok := JSON["status"].(string); ok {
		if status != "ok" {
			return fmt.Errorf("%s", JSON["message"].(map[string]interface{})["errors"])
		}
		return nil
	}
	return fmt.Errorf("unable to unmarshal ```json\n%s\n```", body)
}

func getInfoFromSession(session string) (map[string]interface{}, error) {
	apiUrl := "https://i.instagram.com/api/v1/accounts/current_user/?edit=true"
	request := resty.New().R()
	request.SetHeaders(map[string]string{
		"Cookie":            "sessionid=" + session,
		"X-Ig-Capabilities": "3brTvw==",
		"X-Ig-App-Id":       "567067343352427",
		"User-Agent":        "Instagram 103.1.0.15.119 Android (25/7.1.2; 240dpi; 720x1280; samsung; SM-G988N; z3q; exynos8895; en_US; 164094540)",
		"Content-Type":      "application/x-www-form-urlencoded; charset=UTF-8",
	})
	response, err := request.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	body := response.Body()
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if status, ok := res["status"].(string); ok {
		if status != "ok" {
			return nil, fmt.Errorf("%s\n%s", res["error_title"], res["error_body"])
		}
		return res["user"].(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("unable to unmarshal ```json\n%s\n```", body)
}

func main() {
	var err error
	discord, err = discordgo.New("Bot " + cnfg.Token)
	if err != nil {
		log.Fatal("Error creating Discord session", err)
		return
	}
	discord.Identify.Intents |= discordgo.IntentMessageContent
	discord.AddHandler(onReady)
	discord.AddHandler(messageCreate)
	discord.AddHandler(questionAnswerHandler)
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening connection", err)
		return
	}
	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	if err := discord.Close(); err != nil {
		log.Fatal("Error closing connection", err)
		return
	}
	log.Println("Bot is shutting down")
	saveSession()
}
