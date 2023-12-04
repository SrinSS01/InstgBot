package commands

import (
	"InstgBot/config"
	"InstgBot/info"
	"InstgBot/session"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AccCommand struct {
	Command *discordgo.ApplicationCommand
	Config  *config.Config
}

var Acc = AccCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "acc",
		Description: "Create instagram account",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "session",
				Description: "Session ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	},
}

type GetInfoFromSessionResponse struct {
	Message      string `json:"message"`
	ErrorTitle   string `json:"error_title"`
	ErrorBody    string `json:"error_body"`
	LogoutReason int    `json:"logout_reason"`
	User         struct {
		Biography                   string      `json:"biography"`
		PrimaryProfileLinkType      int         `json:"primary_profile_link_type"`
		ShowFbLinkOnProfile         bool        `json:"show_fb_link_on_profile"`
		ShowFbPageLinkOnProfile     bool        `json:"show_fb_page_link_on_profile"`
		CanHideCategory             bool        `json:"can_hide_category"`
		SmbSupportPartner           interface{} `json:"smb_support_partner"`
		CanAddFbGroupLinkOnProfile  bool        `json:"can_add_fb_group_link_on_profile"`
		IsQuietModeEnabled          bool        `json:"is_quiet_mode_enabled"`
		LastSeenTimezone            string      `json:"last_seen_timezone"`
		AccountCategory             string      `json:"account_category"`
		AllowedCommenterType        string      `json:"allowed_commenter_type"`
		FbidV2                      string      `json:"fbid_v2"`
		FullName                    string      `json:"full_name"`
		Gender                      int         `json:"gender"`
		IsHideMoreCommentEnabled    bool        `json:"is_hide_more_comment_enabled"`
		IsMutedWordsCustomEnabled   bool        `json:"is_muted_words_custom_enabled"`
		IsMutedWordsGlobalEnabled   bool        `json:"is_muted_words_global_enabled"`
		IsMutedWordsSpamscamEnabled bool        `json:"is_muted_words_spamscam_enabled"`
		IsPrivate                   bool        `json:"is_private"`
		HasNmeBadge                 bool        `json:"has_nme_badge"`
		Pk                          int64       `json:"pk"`
		PkId                        string      `json:"pk_id"`
		ReelAutoArchive             string      `json:"reel_auto_archive"`
		StrongId                    string      `json:"strong_id__"`
		BiographyWithEntities       struct {
			RawText  string        `json:"raw_text"`
			Entities []interface{} `json:"entities"`
		} `json:"biography_with_entities"`
		CanLinkEntitiesInBio                       bool          `json:"can_link_entities_in_bio"`
		ExternalUrl                                string        `json:"external_url"`
		Category                                   interface{}   `json:"category"`
		IsCategoryTappable                         bool          `json:"is_category_tappable"`
		IsBusiness                                 bool          `json:"is_business"`
		ProfessionalConversionSuggestedAccountType int           `json:"professional_conversion_suggested_account_type"`
		AccountType                                int           `json:"account_type"`
		DisplayedActionButtonPartner               interface{}   `json:"displayed_action_button_partner"`
		SmbDeliveryPartner                         interface{}   `json:"smb_delivery_partner"`
		SmbSupportDeliveryPartner                  interface{}   `json:"smb_support_delivery_partner"`
		DisplayedActionButtonType                  interface{}   `json:"displayed_action_button_type"`
		IsCallToActionEnabled                      interface{}   `json:"is_call_to_action_enabled"`
		NumOfAdminedPages                          interface{}   `json:"num_of_admined_pages"`
		PageId                                     interface{}   `json:"page_id"`
		PageName                                   interface{}   `json:"page_name"`
		AdsPageId                                  interface{}   `json:"ads_page_id"`
		AdsPageName                                interface{}   `json:"ads_page_name"`
		BioLinks                                   []interface{} `json:"bio_links"`
		AccountBadges                              []interface{} `json:"account_badges"`
		AllMediaCount                              int           `json:"all_media_count"`
		Birthday                                   string        `json:"birthday"`
		BirthdayTodayVisibilityForViewer           string        `json:"birthday_today_visibility_for_viewer"`
		CustomGender                               string        `json:"custom_gender"`
		Email                                      string        `json:"email"`
		HasAnonymousProfilePicture                 bool          `json:"has_anonymous_profile_picture"`
		HdProfilePicUrlInfo                        struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"hd_profile_pic_url_info"`
		HdProfilePicVersions []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Url    string `json:"url"`
		} `json:"hd_profile_pic_versions"`
		InteropMessagingUserFbid               int64  `json:"interop_messaging_user_fbid"`
		IsMv4BBizAssetProfileLocked            bool   `json:"is_mv4b_biz_asset_profile_locked"`
		HasLegacyBbPendingProfilePictureUpdate bool   `json:"has_legacy_bb_pending_profile_picture_update"`
		HasMv4BPendingProfilePictureUpdate     bool   `json:"has_mv4b_pending_profile_picture_update"`
		IsMv4BMaxProfileEditReached            bool   `json:"is_mv4b_max_profile_edit_reached"`
		IsShowingBirthdaySelfie                bool   `json:"is_showing_birthday_selfie"`
		IsSupervisionFeaturesEnabled           bool   `json:"is_supervision_features_enabled"`
		IsVerified                             bool   `json:"is_verified"`
		LikedClipsCount                        int    `json:"liked_clips_count"`
		HasActiveMv4BApplication               bool   `json:"has_active_mv4b_application"`
		PhoneNumber                            string `json:"phone_number"`
		ProfilePicId                           string `json:"profile_pic_id"`
		ProfilePicUrl                          string `json:"profile_pic_url"`
		ProfileEditParams                      struct {
			Username struct {
				ShouldShowConfirmationDialog bool   `json:"should_show_confirmation_dialog"`
				IsPendingReview              bool   `json:"is_pending_review"`
				ConfirmationDialogText       string `json:"confirmation_dialog_text"`
				DisclaimerText               string `json:"disclaimer_text"`
			} `json:"username"`
			FullName struct {
				ShouldShowConfirmationDialog bool   `json:"should_show_confirmation_dialog"`
				IsPendingReview              bool   `json:"is_pending_review"`
				ConfirmationDialogText       string `json:"confirmation_dialog_text"`
				DisclaimerText               string `json:"disclaimer_text"`
			} `json:"full_name"`
		} `json:"profile_edit_params"`
		ShowConversionEditEntry bool   `json:"show_conversion_edit_entry"`
		ShowTogetherPog         bool   `json:"show_together_pog"`
		TrustedUsername         string `json:"trusted_username"`
		TrustDays               int    `json:"trust_days"`
		Username                string `json:"username"`
	} `json:"user"`
	Status string `json:"status"`
}

var followerServices = []struct {
	ID    string
	Price float64
}{
	{"1167", 0.18},
	{"644", 0.24},
	{"946", 0.32},
	{"1436", 0.32},
	{"1166", 0.40},
	{"897", 0.43},
	{"945", 0.75},
}

var SessionObj = session.Session{
	Sessions: []string{},
}
var InfoMap = map[string]*info.Info{}
var postCodeRegex = regexp.MustCompile("\"code\":\"(?P<code>\\w*)\"")
var timeOut = 5 * time.Minute

func (ac *AccCommand) ExecuteDash(s *discordgo.Session, m *discordgo.MessageCreate, sessionId string) {
	authorId := m.Author.ID
	key := authorId + m.GuildID
	if authorId == s.State.User.ID {
		return
	}
	if sessionId == "" {
		if len(SessionObj.Sessions) == 0 {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please provide a session to continue...", m.Reference())
			return
		}
		sessionId = SessionObj.Sessions[0]
		SessionObj.Sessions = SessionObj.Sessions[1:]
	}
	response, err := GetInfoFromSession(sessionId)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
	username := response.User.Username
	i := InfoMap[key]
	if i == nil {
		InfoMap[key] = &info.Info{
			Question: 0,
			URL:      CodeFormat("NA"),
		}
	}
	i = InfoMap[key]
	i.Question = 1
	i.Session = sessionId
	i.Bio = response.User.Biography
	if response.User.ExternalUrl != "" {
		i.URL = response.User.ExternalUrl
	}
	i.Email = response.User.Email
	i.Phone = response.User.PhoneNumber
	msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("The current username is **@%s**, the new username should be?", username), m.Reference())
	i.Timer = time.AfterFunc(timeOut, func() {
		delete(InfoMap, key)
		_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
	})
}
func (ac *AccCommand) QuestionAnswerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	authorId := m.Author.ID
	if authorId == s.State.User.ID {
		return
	}
	key := authorId + m.GuildID
	i := InfoMap[key]
	if i == nil {
		return
	}
	content := m.Content
	switch i.Question {
	case 1:
		if matched, _ := regexp.MatchString("^[a-zA-Z0-9._]+$", content); !matched {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please enter a valid username", m.Reference())
			return
		}
		i.Username = content
		i.Timer.Stop()
		err := PostChanges(i)
		if err != nil {
			msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
			i.Timer = time.AfterFunc(timeOut, func() {
				delete(InfoMap, key)
				_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
			})
			return
		}
		i.Question = 2
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Provide the Name?", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 2:
		i.Timer.Stop()
		if content != "-skip" {
			i.Name = content
			err := PostChanges(i)
			if err != nil {
				msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
				i.Timer = time.AfterFunc(timeOut, func() {
					delete(InfoMap, key)
					_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
				})
				return
			}
		}
		i.Question = 3
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Could you provide your bio information?", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 3:
		i.Timer.Stop()
		if content != "-skip" {
			i.Bio = content
			err := PostChanges(i)
			if err != nil {
				msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
				i.Timer = time.AfterFunc(timeOut, func() {
					delete(InfoMap, key)
					_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
				})
				return
			}
		}
		i.Question = 4
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "What’s the URL you would like to use?", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 4:
		i.Timer.Stop()
		if content != "-skip" {
			if res, _ := regexp.MatchString("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)", content); !res {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please enter a valid url", m.Reference())
				return
			}
			i.URL = content
			err := PostChanges(i)
			if err != nil {
				msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
				i.Timer = time.AfterFunc(timeOut, func() {
					delete(InfoMap, key)
					_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
				})
				return
			}
		}
		i.Question = 5
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "Upload the Posts you want to add.", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
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
		i.Timer.Stop()

		msg := m.Message

		for _, attachment := range attachments {
			code, err := PostMediaFromAttachment(attachment, i.Session)
			if err != nil {
				msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
				i.Timer = time.AfterFunc(timeOut, func() {
					delete(InfoMap, key)
					_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
				})
				return
			}
			i.Posts = append(i.Posts, "https://instagram.com/p/"+code)
			msg, _ = s.ChannelMessageSendReply(msg.ChannelID, attachment.Filename+" posted successfully", msg.Reference())
		}
		i.Question = 6

		msg, _ = s.ChannelMessageSendReply(m.ChannelID, "Upload a profile photo you want to set.", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 6:
		attachments := m.Attachments
		if len(attachments) == 0 {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please provide an image to set as pfp", m.Reference())
			return
		}
		i.Timer.Stop()
		pfp := attachments[0]
		i.Pfp = pfp.URL
		err := PostAvatar(pfp.URL, pfp.Filename, i.Session)
		if err != nil {
			msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Failed to upload pfp\n```\n%s\n```\nPlease try again", err.Error()), m.Reference())
			i.Timer = time.AfterFunc(timeOut, func() {
				delete(InfoMap, key)
				_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
			})
			return
		}
		i.Question = 7

		msg, _ := s.ChannelMessageSendReply(m.ChannelID, "What is the instagram post url you want to add followers to?", m.Reference())
		i.Timer = time.AfterFunc(timeOut, func() {
			delete(InfoMap, key)
			_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
		})
		return
	case 7:
		i.Timer.Stop()
		if content != "-skip" {
			if res, _ := regexp.MatchString("https?://(www\\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)", content); !res {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Please enter a valid url", m.Reference())
				return
			}
			err := ac.AddFollowers(content)
			if err != nil {
				msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
				i.Timer = time.AfterFunc(timeOut, func() {
					delete(InfoMap, key)
					_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
				})
				return
			}
		}
		msg, _ := s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Your Account is ready to be used.",
			Image: &discordgo.MessageEmbedImage{
				URL: i.Pfp,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Name",
					Value: CodeFormat(i.Name),
				},
				{
					Name:  "New Username",
					Value: CodeFormat("@" + i.Username),
				},
				{
					Name:  "Bio",
					Value: CodeFormat(i.Bio),
				},
				{
					Name:  "Url",
					Value: i.URL,
				},
			},
		}, m.Reference())
		_, _ = s.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("%s\n%s", "https://instagram.com/"+i.Username, i.Posts), msg.Reference())
		delete(InfoMap, key)
	}
}
func CodeFormat(str string) string {
	return fmt.Sprintf("`%s`", str)
}
func (ac *AccCommand) ExecuteSlash(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	authorId := interaction.Message.Author.ID
	key := authorId + interaction.GuildID
	data := interaction.ApplicationCommandData()
	sessionId := ""
	if len(data.Options) == 0 {
		if len(SessionObj.Sessions) == 0 {
			_ = s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please provide a session to continue...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		sessionId = SessionObj.Sessions[0]
		SessionObj.Sessions = SessionObj.Sessions[1:]
	} else {
		sessionId = data.Options[0].StringValue()
	}
	_ = s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	response, err := GetInfoFromSession(sessionId)
	if err != nil {
		content := err.Error()
		_, _ = s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}
	username := response.User.Username
	i := InfoMap[key]
	if i == nil {
		InfoMap[key] = &info.Info{
			Question: 0,
			URL:      CodeFormat("NA"),
		}
	}
	i = InfoMap[key]
	i.Question = 1
	i.Session = sessionId
	i.Bio = response.User.Biography
	if response.User.ExternalUrl != "" {
		i.URL = response.User.ExternalUrl
	}
	i.Email = response.User.Email
	i.Phone = response.User.PhoneNumber
	content := fmt.Sprintf("The current username is **@%s**, the new username should be?", username)
	msg, _ := s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	i.Timer = time.AfterFunc(timeOut, func() {
		delete(InfoMap, key)
		_, _ = s.ChannelMessageSendReply(msg.ChannelID, "Question expired", msg.Reference())
	})
}

func (ac *AccCommand) AddFollowers(instagramPostURL string) error {
	APIKey := ac.Config.ApiKey
	APIUrl := ac.Config.ApiUrl
	rand.New(rand.NewSource(time.Now().UnixNano()))
	followersCount := rand.Intn(601) + 200
	selectedFollowersService := followerServices[rand.Intn(len(followerServices))]
	resp, err := resty.New().R().
		SetFormData(map[string]string{
			"key":      APIKey,
			"action":   "add",
			"service":  selectedFollowersService.ID,
			"link":     instagramPostURL,
			"quantity": fmt.Sprintf("%d", followersCount),
		}).Post(APIUrl)
	if err != nil {
		return err
	}
	var apiResponse struct {
		Order int `json:"order"`
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	return nil
}

func PostMediaFromAttachment(attachment *discordgo.MessageAttachment, sessionId string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	request := resty.New().R()
	response, err := request.Get(attachment.URL)
	if err != nil {
		return "", err
	}
	photo := response.Body()
	req, err := http.NewRequest("POST", "https://i.instagram.com/rupload_igphoto/fb_uploader_"+timestamp, bytes.NewBuffer(photo))
	if err != nil {
		return "", err
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
		return "", err
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
		return "", err
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
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	b, err := io.ReadAll(resps.Body)
	if err != nil {
		return "", err
	}
	str := string(b)
	if strings.Contains(str, "\"status\":\"ok\"") {
		matches := postCodeRegex.FindStringSubmatch(str)
		if len(matches) != 0 {
			return matches[postCodeRegex.SubexpIndex("code")], nil
		}
	}
	return "", fmt.Errorf("erro uploading photo")
}
func PostMediaFromFile(file []byte, sessionId string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	req, err := http.NewRequest("POST", "https://i.instagram.com/rupload_igphoto/fb_uploader_"+timestamp, bytes.NewBuffer(file))
	if err != nil {
		return "", err
	}

	req.Header.Set("Host", "i.instagram.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Entity-Type", "image/jpeg")
	req.Header.Set("X-ASBD-ID", "198387")
	req.Header.Set("X-IG-App-ID", "936619743392459")
	req.Header.Set("X-Instagram-AJAX", "1007055396")
	req.Header.Set("X-Entity-Length", strconv.Itoa(len(file)))
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
		return "", err
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
		return "", err
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
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	b, err := io.ReadAll(resps.Body)
	if err != nil {
		return "", err
	}
	str := string(b)
	if strings.Contains(str, "\"status\":\"ok\"") {
		matches := postCodeRegex.FindStringSubmatch(str)
		if len(matches) != 0 {
			return matches[postCodeRegex.SubexpIndex("code")], nil
		}
	}
	return "", fmt.Errorf("erro uploading photo")
}
func PostAvatar(pfpUrl, filename, sessionId string) error {
	response, err := http.DefaultClient.Get(pfpUrl)
	if err != nil {
		return err
	}
	photo := response.Body
	apiUrl := "https://i.instagram.com/accounts/web_change_profile_picture/"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("profile_pic", filename)
	_, err = io.Copy(part, photo)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", apiUrl, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("x-csrftoken", "ARD4xPQFO8DvEghBR6ylbiSx7NSVwMs5")
	req.Header.Set("X-Instagram-AJAX", "7a3a3e64fa87")
	req.Header.Set("Cookie", fmt.Sprintf("mid=YGB8ogALAAESePCJAlGFMopcXIgR; csrftoken=ARD4xPQFO8DvEghBR6ylbiSx7NSVwMs5; sessionid=%s", sessionId))
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	respBody, _ := io.ReadAll(resp.Body)
	if bytes.Contains(respBody, []byte("has_profile_pic\":true")) {
		return nil
	} else {
		return fmt.Errorf("error: Did Not Change")
	}
}

func PostChanges(info *info.Info) error {
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
	if info.URL == CodeFormat("NA") {
		info.URL = ""
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
			if message, ok := JSON["message"].(map[string]interface{}); ok {
				return fmt.Errorf("%s", message["errors"])
			}

			if message, ok := JSON["message"].(string); ok {
				return fmt.Errorf("%s", message)
			}
		}
		return nil
	}
	return fmt.Errorf("unable to unmarshal ```json\n%s\n```", body)
}

func GetInfoFromSession(session string) (*GetInfoFromSessionResponse, error) {
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
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(string(response.Body()))
	}
	var res GetInfoFromSessionResponse
	body := response.Body()
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	//if status, ok := res["status"].(string); ok {
	if res.Status != "ok" {
		errTitle := res.ErrorTitle
		errBody := res.ErrorBody
		message := res.ErrorBody
		builder := strings.Builder{}
		if errTitle != "" {
			builder.WriteString(errTitle)
			builder.WriteByte('\n')
		}
		if errBody != "" {
			builder.WriteString(errBody)
			builder.WriteByte('\n')
		}
		if message != "" {
			builder.WriteString(message)
			builder.WriteByte('\n')
		}
		return nil, fmt.Errorf("%s", builder.String())
	}
	return &res, nil
	//}
	//return nil, fmt.Errorf("unable to unmarshal ```json\n%s\n```", body)
}
