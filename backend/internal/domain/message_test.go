package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewMessage(t *testing.T) {
	msgType := MessageTypeChatMessage
	senderID := uuid.New()
	roomID := uuid.New()
	payload := ChatPayload{Message: "Hello", UserName: "Test User"}

	message := NewMessage(msgType, senderID, roomID, payload)

	if message.Type != msgType {
		t.Errorf("Expected Type %s, got %s", msgType, message.Type)
	}
	if message.SenderUserID != senderID {
		t.Errorf("Expected SenderUserID %s, got %s", senderID, message.SenderUserID)
	}
	if message.RoomID != roomID {
		t.Errorf("Expected RoomID %s, got %s", roomID, message.RoomID)
	}
	if message.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}
	if message.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}
}

func TestNewChatMessage(t *testing.T) {
	senderID := uuid.New()
	roomID := uuid.New()
	messageText := "Hello World"
	userName := "Test User"

	message := NewChatMessage(senderID, roomID, messageText, userName)

	if message.Type != MessageTypeChatMessage {
		t.Errorf("Expected Type %s, got %s", MessageTypeChatMessage, message.Type)
	}
	if message.SenderUserID != senderID {
		t.Errorf("Expected SenderUserID %s, got %s", senderID, message.SenderUserID)
	}
	if message.RoomID != roomID {
		t.Errorf("Expected RoomID %s, got %s", roomID, message.RoomID)
	}

	// Payloadの検証
	payload, ok := message.Payload.(ChatPayload)
	if !ok {
		t.Error("Expected Payload to be ChatPayload")
	}
	if payload.Message != messageText {
		t.Errorf("Expected Message %s, got %s", messageText, payload.Message)
	}
	if payload.UserName != userName {
		t.Errorf("Expected UserName %s, got %s", userName, payload.UserName)
	}
}

func TestNewWebRTCOffer(t *testing.T) {
	senderID := uuid.New()
	targetID := uuid.New()
	roomID := uuid.New()
	sdp := "test-sdp-offer"

	message := NewWebRTCOffer(senderID, targetID, roomID, sdp)

	if message.Type != MessageTypeWebRTCOffer {
		t.Errorf("Expected Type %s, got %s", MessageTypeWebRTCOffer, message.Type)
	}
	if message.SenderUserID != senderID {
		t.Errorf("Expected SenderUserID %s, got %s", senderID, message.SenderUserID)
	}
	if message.TargetUserID != targetID {
		t.Errorf("Expected TargetUserID %s, got %s", targetID, message.TargetUserID)
	}
	if message.RoomID != roomID {
		t.Errorf("Expected RoomID %s, got %s", roomID, message.RoomID)
	}

	// Payloadの検証
	payload, ok := message.Payload.(WebRTCPayload)
	if !ok {
		t.Error("Expected Payload to be WebRTCPayload")
	}
	if payload.SDP != sdp {
		t.Errorf("Expected SDP %s, got %s", sdp, payload.SDP)
	}
	if payload.Type != "offer" {
		t.Errorf("Expected Type 'offer', got %s", payload.Type)
	}
}

func TestNewWebRTCAnswer(t *testing.T) {
	senderID := uuid.New()
	targetID := uuid.New()
	roomID := uuid.New()
	sdp := "test-sdp-answer"

	message := NewWebRTCAnswer(senderID, targetID, roomID, sdp)

	if message.Type != MessageTypeWebRTCAnswer {
		t.Errorf("Expected Type %s, got %s", MessageTypeWebRTCAnswer, message.Type)
	}
	if message.TargetUserID != targetID {
		t.Errorf("Expected TargetUserID %s, got %s", targetID, message.TargetUserID)
	}

	// Payloadの検証
	payload, ok := message.Payload.(WebRTCPayload)
	if !ok {
		t.Error("Expected Payload to be WebRTCPayload")
	}
	if payload.Type != "answer" {
		t.Errorf("Expected Type 'answer', got %s", payload.Type)
	}
}

func TestNewICECandidate(t *testing.T) {
	senderID := uuid.New()
	targetID := uuid.New()
	roomID := uuid.New()
	candidate := "test-candidate"
	sdpMid := "test-mid"
	sdpMLineIndex := 0

	message := NewICECandidate(senderID, targetID, roomID, candidate, sdpMid, sdpMLineIndex)

	if message.Type != MessageTypeICECandidate {
		t.Errorf("Expected Type %s, got %s", MessageTypeICECandidate, message.Type)
	}

	// Payloadの検証
	payload, ok := message.Payload.(ICECandidatePayload)
	if !ok {
		t.Error("Expected Payload to be ICECandidatePayload")
	}
	if payload.Candidate != candidate {
		t.Errorf("Expected Candidate %s, got %s", candidate, payload.Candidate)
	}
	if payload.SDPMid != sdpMid {
		t.Errorf("Expected SDPMid %s, got %s", sdpMid, payload.SDPMid)
	}
	if payload.SDPMLineIndex != sdpMLineIndex {
		t.Errorf("Expected SDPMLineIndex %d, got %d", sdpMLineIndex, payload.SDPMLineIndex)
	}
}

func TestMessage_IsValid(t *testing.T) {
	senderID := uuid.New()
	targetID := uuid.New()
	roomID := uuid.New()

	tests := []struct {
		name     string
		message  *Message
		expected bool
	}{
		{
			name: "Valid chat message",
			message: &Message{
				Type:         MessageTypeChatMessage,
				SenderUserID: senderID,
				RoomID:       roomID,
			},
			expected: true,
		},
		{
			name: "Valid WebRTC offer",
			message: &Message{
				Type:         MessageTypeWebRTCOffer,
				SenderUserID: senderID,
				TargetUserID: targetID,
				RoomID:       roomID,
			},
			expected: true,
		},
		{
			name: "Invalid - missing type",
			message: &Message{
				SenderUserID: senderID,
				RoomID:       roomID,
			},
			expected: false,
		},
		{
			name: "Invalid - missing room ID",
			message: &Message{
				Type:         MessageTypeChatMessage,
				SenderUserID: senderID,
			},
			expected: false,
		},
		{
			name: "Invalid chat message - missing sender",
			message: &Message{
				Type:   MessageTypeChatMessage,
				RoomID: roomID,
			},
			expected: false,
		},
		{
			name: "Invalid WebRTC offer - missing target",
			message: &Message{
				Type:         MessageTypeWebRTCOffer,
				SenderUserID: senderID,
				RoomID:       roomID,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.IsValid()
			if result != tt.expected {
				t.Errorf("Expected IsValid() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMessage_IsDirectMessage(t *testing.T) {
	senderID := uuid.New()
	targetID := uuid.New()
	roomID := uuid.New()

	// Direct message (with target)
	directMessage := &Message{
		Type:         MessageTypeWebRTCOffer,
		SenderUserID: senderID,
		TargetUserID: targetID,
		RoomID:       roomID,
	}

	if !directMessage.IsDirectMessage() {
		t.Error("Expected message with TargetUserID to be direct message")
	}

	// Broadcast message (no target)
	broadcastMessage := &Message{
		Type:         MessageTypeChatMessage,
		SenderUserID: senderID,
		RoomID:       roomID,
	}

	if broadcastMessage.IsDirectMessage() {
		t.Error("Expected message without TargetUserID to not be direct message")
	}
}

func TestMessage_IsBroadcastMessage(t *testing.T) {
	roomID := uuid.New()

	tests := []struct {
		name        string
		messageType MessageType
		expected    bool
	}{
		{
			name:        "Chat message is broadcast",
			messageType: MessageTypeChatMessage,
			expected:    true,
		},
		{
			name:        "User joined is broadcast",
			messageType: MessageTypeUserJoined,
			expected:    true,
		},
		{
			name:        "User left is broadcast",
			messageType: MessageTypeUserLeft,
			expected:    true,
		},
		{
			name:        "WebRTC offer is not broadcast",
			messageType: MessageTypeWebRTCOffer,
			expected:    false,
		},
		{
			name:        "Mute user is not broadcast",
			messageType: MessageTypeMuteUser,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				Type:   tt.messageType,
				RoomID: roomID,
			}

			result := message.IsBroadcastMessage()
			if result != tt.expected {
				t.Errorf("Expected IsBroadcastMessage() to return %v for %s, got %v", tt.expected, tt.messageType, result)
			}
		})
	}
}

func TestMessage_ToJSON_FromJSON(t *testing.T) {
	senderID := uuid.New()
	roomID := uuid.New()
	originalMessage := NewChatMessage(senderID, roomID, "Test message", "Test User")

	// ToJSON
	jsonData, err := originalMessage.ToJSON()
	if err != nil {
		t.Errorf("Expected no error from ToJSON, got %v", err)
	}

	// FromJSON
	parsedMessage, err := FromJSON(jsonData)
	if err != nil {
		t.Errorf("Expected no error from FromJSON, got %v", err)
	}

	// 基本フィールドの比較
	if parsedMessage.ID != originalMessage.ID {
		t.Errorf("Expected ID %s, got %s", originalMessage.ID, parsedMessage.ID)
	}
	if parsedMessage.Type != originalMessage.Type {
		t.Errorf("Expected Type %s, got %s", originalMessage.Type, parsedMessage.Type)
	}
	if parsedMessage.SenderUserID != originalMessage.SenderUserID {
		t.Errorf("Expected SenderUserID %s, got %s", originalMessage.SenderUserID, parsedMessage.SenderUserID)
	}
	if parsedMessage.RoomID != originalMessage.RoomID {
		t.Errorf("Expected RoomID %s, got %s", originalMessage.RoomID, parsedMessage.RoomID)
	}

	// Timestampの比較（JSONでの精度の問題を考慮）
	if !parsedMessage.Timestamp.Truncate(time.Second).Equal(originalMessage.Timestamp.Truncate(time.Second)) {
		t.Errorf("Expected Timestamp %v, got %v", originalMessage.Timestamp, parsedMessage.Timestamp)
	}
}

func TestFromJSON_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)

	_, err := FromJSON(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestMessage_TimestampIsRecent(t *testing.T) {
	before := time.Now()
	message := NewMessage(MessageTypeChatMessage, uuid.New(), uuid.New(), nil)
	after := time.Now()

	if message.Timestamp.Before(before) || message.Timestamp.After(after) {
		t.Errorf("Expected Timestamp to be between %v and %v, got %v", before, after, message.Timestamp)
	}
}
