package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	LangCode  string `json:"language_code"`
}

type Chat struct {
	Id                  int      `json:"id"`
	Type                string   `json:"type"`
	Title               string   `json:"title"`
	UserName            string   `json:"username"`
	FirstName           string   `json:"first_name"`
	LastName            string   `json:"last_name"`
	AllMembersAreAdmins bool     `json:"all_members_are_administrators"`
	Description         string   `json:"description"`
	InviteLink          string   `json:"invite_link"`
	PinnedMessage       *Message `json:"message"`
	StickerSetName      string   `json:"sticker_set_name"`
	CanSetStickerSet    bool     `json:"can_set_sticker_set"`
}

func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

type Message struct {
	Id                   int64            `json:"id"`
	From                 *User            `json:"from"`
	Date                 int64            `json:"date"`
	Chat                 *Chat            `json:"chat"`
	ForwardFrom          *User            `json:"forward_from"`
	ForwardFromChar      *Chat            `json:"forward_from_chat"`
	ForwardFromMessageId int64            `json:"forward_from_message_id"`
	AutorSignature       string           `json:"author_signature"`
	Text                 string           `json:"text"`
	Entities             []*MessageEntity `json:"entities"`
	CaptionEntities      []*MessageEntity `json:"caption_entities"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	Self   *User  `json:"user"`
}

type BotResponseParameters struct {
	migrate_to_chat_id int64
	retry_after        int
}

type BotResponse struct {
	Ok          bool                   `json:"ok"`
	Result      json.RawMessage        `json:"result"`
	Description string                 `json:"description"`
	ErrorCode   int                    `json:"error_code"`
	Parameters  *BotResponseParameters `json:"parametrs"`
}

type Update struct {
	UpdateId          int      `json:"update_id"`
	Msg               *Message `json:"message"`
	EditedMsg         *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
}

func (t *Update) Print() {
	fmt.Printf(
		"UpdateID: %v,\nMsg: %v,\nEditedMsg: %v,\nChannelPost: %v,\nEditedChannelPost: %v\n",
		t.UpdateId, t.Msg, t.EditedMsg, t.ChannelPost, t.EditedChannelPost,
	)
}
