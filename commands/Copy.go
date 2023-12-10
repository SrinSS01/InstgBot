package commands

import (
	"InstgBot/info"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"net/http"
	"regexp"
	"time"
)

type CopyCommand struct {
	Command *discordgo.ApplicationCommand
}

var Copy = CopyCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "copy",
		Description: "Copy user account and create a duplicate copy account",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "username",
				Description: "Username to copy details from",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

type InstgAccResponse struct {
	Data struct {
		User *struct {
			AiAgentType interface{} `json:"ai_agent_type"`
			Biography   string      `json:"biography"`
			BioLinks    []struct {
				Title    string `json:"title"`
				LynxUrl  string `json:"lynx_url"`
				Url      string `json:"url"`
				LinkType string `json:"link_type"`
			} `json:"bio_links"`
			FbProfileBiolink      interface{} `json:"fb_profile_biolink"`
			BiographyWithEntities struct {
				RawText  string        `json:"raw_text"`
				Entities []interface{} `json:"entities"`
			} `json:"biography_with_entities"`
			BlockedByViewer        bool        `json:"blocked_by_viewer"`
			RestrictedByViewer     interface{} `json:"restricted_by_viewer"`
			CountryBlock           bool        `json:"country_block"`
			EimuId                 string      `json:"eimu_id"`
			ExternalUrl            string      `json:"external_url"`
			ExternalUrlLinkshimmed string      `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
			Fbid             string `json:"fbid"`
			FollowedByViewer bool   `json:"followed_by_viewer"`
			EdgeFollow       struct {
				Count int `json:"count"`
			} `json:"edge_follow"`
			FollowsViewer         bool        `json:"follows_viewer"`
			FullName              string      `json:"full_name"`
			GroupMetadata         interface{} `json:"group_metadata"`
			HasArEffects          bool        `json:"has_ar_effects"`
			HasClips              bool        `json:"has_clips"`
			HasGuides             bool        `json:"has_guides"`
			HasChannel            bool        `json:"has_channel"`
			HasBlockedViewer      bool        `json:"has_blocked_viewer"`
			HighlightReelCount    int         `json:"highlight_reel_count"`
			HasRequestedViewer    bool        `json:"has_requested_viewer"`
			HideLikeAndViewCounts bool        `json:"hide_like_and_view_counts"`
			Id                    string      `json:"id"`
			IsBusinessAccount     bool        `json:"is_business_account"`
			IsProfessionalAccount bool        `json:"is_professional_account"`
			IsSupervisionEnabled  bool        `json:"is_supervision_enabled"`
			IsGuardianOfViewer    bool        `json:"is_guardian_of_viewer"`
			IsSupervisedByViewer  bool        `json:"is_supervised_by_viewer"`
			IsSupervisedUser      bool        `json:"is_supervised_user"`
			IsEmbedsDisabled      bool        `json:"is_embeds_disabled"`
			IsJoinedRecently      bool        `json:"is_joined_recently"`
			GuardianId            interface{} `json:"guardian_id"`
			BusinessAddressJson   interface{} `json:"business_address_json"`
			BusinessContactMethod string      `json:"business_contact_method"`
			BusinessEmail         interface{} `json:"business_email"`
			BusinessPhoneNumber   interface{} `json:"business_phone_number"`
			BusinessCategoryName  interface{} `json:"business_category_name"`
			OverallCategoryName   interface{} `json:"overall_category_name"`
			CategoryEnum          interface{} `json:"category_enum"`
			CategoryName          interface{} `json:"category_name"`
			IsPrivate             bool        `json:"is_private"`
			IsVerified            bool        `json:"is_verified"`
			IsVerifiedByMv4B      bool        `json:"is_verified_by_mv4b"`
			IsRegulatedC18        bool        `json:"is_regulated_c18"`
			EdgeMutualFollowedBy  struct {
				Count int           `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_mutual_followed_by"`
			PinnedChannelsListCount        int         `json:"pinned_channels_list_count"`
			ProfilePicUrl                  string      `json:"profile_pic_url"`
			ProfilePicUrlHd                string      `json:"profile_pic_url_hd"`
			RequestedByViewer              bool        `json:"requested_by_viewer"`
			ShouldShowCategory             bool        `json:"should_show_category"`
			ShouldShowPublicContacts       bool        `json:"should_show_public_contacts"`
			ShowAccountTransparencyDetails bool        `json:"show_account_transparency_details"`
			TransparencyLabel              interface{} `json:"transparency_label"`
			TransparencyProduct            string      `json:"transparency_product"`
			Username                       string      `json:"username"`
			ConnectedFbPage                interface{} `json:"connected_fb_page"`
			Pronouns                       []string    `json:"pronouns"`
			EdgeFelixVideoTimeline         struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_video_timeline"`
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []Post `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_saved_media"`
			EdgeMediaCollections struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_collections"`
			EdgeRelatedProfiles struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_related_profiles"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}

type Post struct {
	Node struct {
		Typename   string `json:"__typename"`
		Id         string `json:"id"`
		Shortcode  string `json:"shortcode"`
		Dimensions struct {
			Height int `json:"height"`
			Width  int `json:"width"`
		} `json:"dimensions"`
		DisplayUrl            string `json:"display_url"`
		EdgeMediaToTaggedUser struct {
			Edges []interface{} `json:"edges"`
		} `json:"edge_media_to_tagged_user"`
		FactCheckOverallRating interface{} `json:"fact_check_overall_rating"`
		FactCheckInformation   interface{} `json:"fact_check_information"`
		GatingInfo             interface{} `json:"gating_info"`
		SharingFrictionInfo    struct {
			ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
			BloksAppUrl               interface{} `json:"bloks_app_url"`
		} `json:"sharing_friction_info"`
		MediaOverlayInfo interface{} `json:"media_overlay_info"`
		MediaPreview     string      `json:"media_preview"`
		Owner            struct {
			Id       string `json:"id"`
			Username string `json:"username"`
		} `json:"owner"`
		IsVideo              bool        `json:"is_video"`
		HasUpcomingEvent     bool        `json:"has_upcoming_event"`
		AccessibilityCaption interface{} `json:"accessibility_caption"`
		DashInfo             struct {
			IsDashEligible    bool        `json:"is_dash_eligible"`
			VideoDashManifest interface{} `json:"video_dash_manifest"`
			NumberOfQualities int         `json:"number_of_qualities"`
		} `json:"dash_info"`
		HasAudio           bool   `json:"has_audio"`
		TrackingToken      string `json:"tracking_token"`
		VideoUrl           string `json:"video_url"`
		VideoViewCount     int    `json:"video_view_count"`
		EdgeMediaToCaption struct {
			Edges []struct {
				Node struct {
					Text string `json:"text"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"edge_media_to_caption"`
		EdgeMediaToComment struct {
			Count int `json:"count"`
		} `json:"edge_media_to_comment"`
		CommentsDisabled bool `json:"comments_disabled"`
		TakenAtTimestamp int  `json:"taken_at_timestamp"`
		EdgeLikedBy      struct {
			Count int `json:"count"`
		} `json:"edge_liked_by"`
		EdgeMediaPreviewLike struct {
			Count int `json:"count"`
		} `json:"edge_media_preview_like"`
		Location           interface{} `json:"location"`
		NftAssetInfo       interface{} `json:"nft_asset_info"`
		ThumbnailSrc       string      `json:"thumbnail_src"`
		ThumbnailResources []struct {
			Src          string `json:"src"`
			ConfigWidth  int    `json:"config_width"`
			ConfigHeight int    `json:"config_height"`
		} `json:"thumbnail_resources"`
		FelixProfileGridCrop interface{}   `json:"felix_profile_grid_crop"`
		CoauthorProducers    []interface{} `json:"coauthor_producers"`
		PinnedForUsers       []struct {
			Id            string `json:"id"`
			IsVerified    bool   `json:"is_verified"`
			ProfilePicUrl string `json:"profile_pic_url"`
			Username      string `json:"username"`
		} `json:"pinned_for_users"`
		ViewerCanReshare          bool   `json:"viewer_can_reshare"`
		ProductType               string `json:"product_type"`
		ClipsMusicAttributionInfo struct {
			ArtistName            string `json:"artist_name"`
			SongName              string `json:"song_name"`
			UsesOriginalAudio     bool   `json:"uses_original_audio"`
			ShouldMuteAudio       bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason string `json:"should_mute_audio_reason"`
			AudioId               string `json:"audio_id"`
		} `json:"clips_music_attribution_info"`

		EdgeSidecarToChildren struct {
			Edges []Post `json:"edges"`
		} `json:"edge_sidecar_to_children"`
	} `json:"node"`
}

var regex = regexp.MustCompile("(?P<oldUsername>[a-zA-Z0-9._]+)")

func (c *CopyCommand) ExecuteDash(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	matches := regex.FindStringSubmatch(args)
	if len(matches) == 0 {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "Invalid args", m.Reference())
		return
	}
	oldUsername := matches[regex.SubexpIndex("oldUsername")]
	//newUsername := matches[regex.SubexpIndex("newUsername")]
	if len(SessionObj.Sessions) == 0 {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "No sessions available", m.Reference())
		return
	}
	sessionId := SessionObj.Sessions[0]
	err, inf := GetInfoFromUsername(oldUsername)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
	inf.Session = sessionId
	infoFromSession, err := GetInfoFromSession(inf.Session)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
		return
	}

	// ask for username
	embed := &discordgo.MessageEmbed{
		Description: "Do you want to set a new username?",
		Color:       0x00ffff,
	}
	msg, _ := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Embeds:    []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Yes",
						Style:    discordgo.SuccessButton,
						CustomID: "yes",
					},
					discordgo.Button{
						Label:    "No",
						Style:    discordgo.DangerButton,
						CustomID: "no",
					},
				},
			},
		},
	})
	res, err := waitYesNoResponse(s, m.Author.ID, 20*time.Second, func(session *discordgo.Session, i *discordgo.InteractionCreate) {

	})
	msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    msg.ChannelID,
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{},
		ID:         msg.ID,
		Flags:      msg.Flags,
	})
	if err != nil || res == "no" {
		inf.Username = infoFromSession.User.Username
	}
	// end of ask for username

	inf.Email = infoFromSession.User.Email
	inf.Phone = infoFromSession.User.PhoneNumber
	err = PostChanges(inf)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
		return
	}
	err = PostAvatar(inf.Pfp, fmt.Sprintf("%d.png", time.Now().Unix()), inf.Session)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error()), m.Reference())
		return
	}
	SessionObj.Sessions = SessionObj.Sessions[1:]
	msg, _ = s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Your Account is ready to be used.",
		Image: &discordgo.MessageEmbedImage{
			URL: inf.Pfp,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Name",
				Value: inf.Name,
			},
			{
				Name:  "New Username",
				Value: "@" + inf.Username,
			},
			{
				Name:  "Bio",
				Value: inf.Bio,
			},
			{
				Name:  "Url",
				Value: inf.URL,
			},
		},
	}, m.Reference())
	msg, _ = s.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("%s", "https://instagram.com/"+inf.Username), msg.Reference())
	embed = &discordgo.MessageEmbed{
		Description: "Do you want copy the posts?",
		Color:       0x00ffff,
	}
	msg, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: msg.Reference(),
		Embeds:    []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Yes",
						Style:    discordgo.SuccessButton,
						CustomID: "yes",
					},
					discordgo.Button{
						Label:    "No",
						Style:    discordgo.DangerButton,
						CustomID: "no",
					},
				},
			},
		},
	})
	res, err = waitYesNoResponse(s, m.Author.ID, 20*time.Second, nil)
	msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    msg.ChannelID,
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{},
		ID:         msg.ID,
		Flags:      msg.Flags,
	})
	if err != nil || res == "no" {
		return
	}
	request := resty.New().R()
	//apiUrl := "https://instgbot-nextjs.vercel.app/api/og?"
	//requestBuilder := strings.Builder{}
	//requestBuilder.WriteString(apiUrl)
	var files []*discordgo.File
	buffers := map[string][]byte{}
	for i, post := range inf.Posts {
		if i == 6 {
			break
		}
		response, err := request.Get(post)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		filename := fmt.Sprintf("%d.png", time.Now().Unix())
		body := response.Body()
		buffers[fmt.Sprintf("%d", i+1)] = body
		reader := bytes.NewReader(body)
		files = append(files, &discordgo.File{
			Name:   filename,
			Reader: reader,
		})
	}
	msg, err = s.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Files: files,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "select_posts",
						Placeholder: "Please select the posts you want to copy",
						MenuType:    discordgo.StringSelectMenu,
						MaxValues:   6,
						Options: []discordgo.SelectMenuOption{
							{
								Label: "First",
								Value: "1",
							},
							{
								Label: "Second",
								Value: "2",
							},
							{
								Label: "Third",
								Value: "3",
							},
							{
								Label: "Fourth",
								Value: "4",
							},
							{
								Label: "Fifth",
								Value: "5",
							},
							{
								Label: "Sixth",
								Value: "6",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
	selections, err := awaitSelection(s, m.Author.ID, 20*time.Second)
	msg, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    msg.ChannelID,
		Components: []discordgo.MessageComponent{},
		ID:         msg.ID,
		Flags:      msg.Flags,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, selection := range selections {
		code, err := PostMediaFromFile(buffers[selection], sessionId)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("unable to send post\n```\n%s\n```", err.Error()))
			continue
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<https://instagram.com/p/%s> posted successfully", code))
	}
}

func (c *CopyCommand) ExecuteSlash(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if len(SessionObj.Sessions) == 0 {
		_ = s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No sessions available",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	data := interaction.ApplicationCommandData()
	username := data.Options[0].StringValue()
	sessionId := SessionObj.Sessions[0]
	SessionObj.Sessions = SessionObj.Sessions[1:]
	_ = s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	err, inf := GetInfoFromUsername(username)
	if err != nil {
		e := err.Error()
		_, _ = s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &e,
		})
		return
	}
	inf.Session = sessionId
	infoFromSession, err := GetInfoFromSession(inf.Session)
	if err != nil {
		errorStr := fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error())
		_, _ = s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &errorStr,
		})
		return
	}
	inf.Username = infoFromSession.User.Username
	inf.Email = infoFromSession.User.Email
	inf.Phone = infoFromSession.User.PhoneNumber
	err = PostChanges(inf)
	if err != nil {
		errorStr := fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error())
		_, _ = s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &errorStr,
		})
		return
	}
	err = PostAvatar(inf.Pfp, fmt.Sprintf("%d.png", time.Now().Unix()), inf.Session)
	if err != nil {
		errorStr := fmt.Sprintf("```\n%s\n```\nPlease try again", err.Error())
		_, _ = s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &errorStr,
		})
		return
	}
	m, _ := s.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Your Account is ready to be used.",
				Image: &discordgo.MessageEmbedImage{
					URL: inf.Pfp,
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Name",
						Value: CodeFormat(inf.Name),
					},
					{
						Name:  "New Username",
						Value: CodeFormat("@" + inf.Username),
					},
					{
						Name:  "Bio",
						Value: CodeFormat(inf.Bio),
					},
					{
						Name:  "Url",
						Value: inf.URL,
					},
				},
			},
		},
	})
	m, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("%s", "https://instagram.com/"+inf.Username), m.Reference())
	embed := &discordgo.MessageEmbed{
		Description: "Do you want copy the posts?",
		Color:       0x00ffff,
	}
	m, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Embeds:    []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Yes",
						Style:    discordgo.SuccessButton,
						CustomID: "yes",
					},
					discordgo.Button{
						Label:    "No",
						Style:    discordgo.DangerButton,
						CustomID: "no",
					},
				},
			},
		},
	})
	res, err := waitYesNoResponse(s, m.Author.ID, 20*time.Second, nil)
	m, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    m.ChannelID,
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{},
		ID:         m.ID,
		Flags:      m.Flags,
	})
	if err != nil || res == "no" {
		return
	}
	request := resty.New().R()
	//apiUrl := "https://instgbot-nextjs.vercel.app/api/og?"
	//requestBuilder := strings.Builder{}
	//requestBuilder.WriteString(apiUrl)
	var files []*discordgo.File
	buffers := map[string][]byte{}
	for i, post := range inf.Posts {
		if i == 6 {
			break
		}
		response, err := request.Get(post)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		filename := fmt.Sprintf("%d.png", time.Now().Unix())
		body := response.Body()
		buffers[fmt.Sprintf("%d", i+1)] = body
		reader := bytes.NewReader(body)
		files = append(files, &discordgo.File{
			Name:   filename,
			Reader: reader,
		})
	}
	m, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: files,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "select_posts",
						Placeholder: "Please select the posts you want to copy",
						MenuType:    discordgo.StringSelectMenu,
						MaxValues:   6,
						Options: []discordgo.SelectMenuOption{
							{
								Label: "First",
								Value: "1",
							},
							{
								Label: "Second",
								Value: "2",
							},
							{
								Label: "Third",
								Value: "3",
							},
							{
								Label: "Fourth",
								Value: "4",
							},
							{
								Label: "Fifth",
								Value: "5",
							},
							{
								Label: "Sixth",
								Value: "6",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
	selections, err := awaitSelection(s, m.Author.ID, 20*time.Second)
	m, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    m.ChannelID,
		Components: []discordgo.MessageComponent{},
		ID:         m.ID,
		Flags:      m.Flags,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, selection := range selections {
		code, err := PostMediaFromFile(buffers[selection], sessionId)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("unable to send post\n```\n%s\n```", err.Error()))
			continue
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<https://instagram.com/p/%s> posted successfully", code))
	}
}

func GetInfoFromUsername(username string) (error, *info.Info) {
	apiUrl := "https://www.instagram.com/api/v1/users/web_profile_info/?username=" + username
	r := resty.New().R().SetHeaders(map[string]string{
		"Host":             "i.instagram.com",
		"Connection":       "keep-alive",
		"Accept":           "*/*",
		"X-Entity-Type":    "image/jpeg",
		"X-ASBD-ID":        "198387",
		"X-IG-App-ID":      "936619743392459",
		"X-Instagram-AJAX": "1007055396",
		"X-Entity-Length":  "14911",
		"Accept-Language":  "en-US,en;q=0.9",
		"Origin":           "https://www.instagram.com",
		"Referer":          "https://www.instagram.com/",
		"Offset":           "0",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
	})
	response, err := r.Get(apiUrl)
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf(response.Status()), nil
	}
	if err != nil {
		return err, nil
	}
	var res InstgAccResponse
	err = json.Unmarshal(response.Body(), &res)
	if err != nil {
		return err, nil
	}
	if res.Status != "ok" || res.Data.User == nil {
		return fmt.Errorf("response returned error"), nil
	}
	inf := info.Info{
		URL: CodeFormat("NA"),
	}
	//inf.Username = res.Data.User.Username
	inf.Name = res.Data.User.FullName
	inf.Bio = res.Data.User.Biography
	inf.Pfp = res.Data.User.ProfilePicUrlHd
	externalUrl := res.Data.User.ExternalUrl
	if externalUrl != "" {
		inf.URL = externalUrl
	}
	inf.Posts = make([]string, 0)
	for _, post := range res.Data.User.EdgeOwnerToTimelineMedia.Edges {
		if post.Node.IsVideo {
			//inf.Posts = append(inf.Posts, post.Node.VideoUrl)
			continue
		}
		if post.Node.Typename == "GraphImage" {
			inf.Posts = append(inf.Posts, post.Node.DisplayUrl)
			continue
		}

		if post.Node.Typename == "GraphSidecar" {
			for _, p := range post.Node.EdgeSidecarToChildren.Edges {
				if p.Node.IsVideo {
					//inf.Posts = append(inf.Posts, post.Node.VideoUrl)
					continue
				}
				inf.Posts = append(inf.Posts, p.Node.DisplayUrl)
			}
		}
	}
	return nil, &inf
}

func waitYesNoResponse(s *discordgo.Session, userID string, timeout time.Duration, ifYes func(session *discordgo.Session, i *discordgo.InteractionCreate)) (string, error) {
	res := make(chan string)
	defer s.AddHandlerOnce(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		_ = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if i.Type != discordgo.InteractionMessageComponent || i.Member.User.ID != userID {
			return
		}
		if i.MessageComponentData().CustomID == "yes" && ifYes != nil {
			ifYes(session, i)
		}
		res <- i.MessageComponentData().CustomID
	})()
	timer := time.NewTimer(timeout)
	select {
	case r := <-res:
		timer.Stop()
		return r, nil
	case <-timer.C:
		return "", fmt.Errorf("timeout")
	}
}

func awaitSelection(s *discordgo.Session, userId string, timeout time.Duration) ([]string, error) {
	res := make(chan []string)
	defer s.AddHandlerOnce(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Member.User.ID != userId || i.Type != discordgo.InteractionMessageComponent || i.MessageComponentData().ComponentType != discordgo.SelectMenuComponent || i.MessageComponentData().CustomID != "select_posts" {
			return
		}
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		res <- i.MessageComponentData().Values
	})()
	timer := time.NewTimer(timeout)
	select {
	case r := <-res:
		timer.Stop()
		return r, nil
	case <-timer.C:
		return nil, fmt.Errorf("timeout")
	}
}
