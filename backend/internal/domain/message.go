package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of message
type MessageType string

const (
	// WebRTC シグナリング
	MessageTypeWebRTCOffer     MessageType = "webrtc_offer"
	MessageTypeWebRTCAnswer    MessageType = "webrtc_answer"
	MessageTypeICECandidate    MessageType = "ice_candidate"

	// ルーム管理
	MessageTypeJoinRoom  MessageType = "join_room"
	MessageTypeLeaveRoom MessageType = "leave_room"
	MessageTypeUserJoined MessageType = "user_joined"
	MessageTypeUserLeft   MessageType = "user_left"

	// チャット
	MessageTypeChatMessage MessageType = "chat_message"

	// 制御
	MessageTypeMuteUser    MessageType = "mute_user"
	MessageTypeAdmitUser   MessageType = "admit_user"
	MessageTypeScreenShare MessageType = "screen_share"
)

// Message represents a real-time message
type Message struct {
	ID           uuid.UUID   `json:"id"`
	Type         MessageType `json:"type"`
	SenderUserID uuid.UUID   `json:"senderUserId,omitempty"`
	TargetUserID uuid.UUID   `json:"targetUserId,omitempty"`
	RoomID       uuid.UUID   `json:"roomId"`
	Payload      interface{} `json:"payload"`
	Timestamp    time.Time   `json:"timestamp"`
}

// ChatPayload represents chat message payload
type ChatPayload struct {
	Message string `json:"message"`
	UserName string `json:"userName"`
}

// WebRTCPayload represents WebRTC signaling payload
type WebRTCPayload struct {
	SDP  string `json:"sdp,omitempty"`
	Type string `json:"type,omitempty"`
}

// ICECandidatePayload represents ICE candidate payload
type ICECandidatePayload struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

// ControlPayload represents control message payload
type ControlPayload struct {
	Action   string    `json:"action"`
	TargetID uuid.UUID `json:"targetId,omitempty"`
	Reason   string    `json:"reason,omitempty"`
}

// NewMessage creates a new message
func NewMessage(msgType MessageType, senderID, roomID uuid.UUID, payload interface{}) *Message {
	return &Message{
		ID:           uuid.New(),
		Type:         msgType,
		SenderUserID: senderID,
		RoomID:       roomID,
		Payload:      payload,
		Timestamp:    time.Now(),
	}
}

// NewChatMessage creates a new chat message
func NewChatMessage(senderID, roomID uuid.UUID, message, userName string) *Message {
	payload := ChatPayload{
		Message:  message,
		UserName: userName,
	}
	return NewMessage(MessageTypeChatMessage, senderID, roomID, payload)
}

// NewWebRTCOffer creates a new WebRTC offer message
func NewWebRTCOffer(senderID, targetID, roomID uuid.UUID, sdp string) *Message {
	payload := WebRTCPayload{
		SDP:  sdp,
		Type: "offer",
	}
	msg := NewMessage(MessageTypeWebRTCOffer, senderID, roomID, payload)
	msg.TargetUserID = targetID
	return msg
}

// NewWebRTCAnswer creates a new WebRTC answer message
func NewWebRTCAnswer(senderID, targetID, roomID uuid.UUID, sdp string) *Message {
	payload := WebRTCPayload{
		SDP:  sdp,
		Type: "answer",
	}
	msg := NewMessage(MessageTypeWebRTCAnswer, senderID, roomID, payload)
	msg.TargetUserID = targetID
	return msg
}

// NewICECandidate creates a new ICE candidate message
func NewICECandidate(senderID, targetID, roomID uuid.UUID, candidate, sdpMid string, sdpMLineIndex int) *Message {
	payload := ICECandidatePayload{
		Candidate:     candidate,
		SDPMid:        sdpMid,
		SDPMLineIndex: sdpMLineIndex,
	}
	msg := NewMessage(MessageTypeICECandidate, senderID, roomID, payload)
	msg.TargetUserID = targetID
	return msg
}

// IsValid validates the message
func (m *Message) IsValid() bool {
	if m.Type == "" || m.RoomID == uuid.Nil {
		return false
	}

	// メッセージタイプに応じた検証
	switch m.Type {
	case MessageTypeChatMessage:
		return m.SenderUserID != uuid.Nil
	case MessageTypeWebRTCOffer, MessageTypeWebRTCAnswer, MessageTypeICECandidate:
		return m.SenderUserID != uuid.Nil && m.TargetUserID != uuid.Nil
	case MessageTypeMuteUser, MessageTypeAdmitUser:
		return m.SenderUserID != uuid.Nil
	default:
		return true
	}
}

// IsDirectMessage checks if the message is targeted to a specific user
func (m *Message) IsDirectMessage() bool {
	return m.TargetUserID != uuid.Nil
}

// IsBroadcastMessage checks if the message should be broadcast to all room participants
func (m *Message) IsBroadcastMessage() bool {
	switch m.Type {
	case MessageTypeChatMessage, MessageTypeUserJoined, MessageTypeUserLeft:
		return true
	default:
		return false
	}
}

// ToJSON converts message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON creates message from JSON
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
