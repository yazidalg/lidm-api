# Opponent Taunting System

The opponent taunting system allows players to send emoji reactions, animations, or taunts to each other during multiplayer quiz games.

## Socket.IO Events

### Sending a Taunt (`opponent_taunt`)

**Client sends:**
```javascript
socket.emit('opponent_taunt', {
  quiz_id: 123,
  user_id: 456,
  link_lottie: "https://example.com/animations/laugh.json", // Lottie animation URL
  type: "reaction" // Optional: "reaction", "taunt", "emoji", etc.
});
```

### Receiving a Taunt (`taunt_received`)

**Client receives:**
```javascript
socket.on('taunt_received', (data) => {
  console.log(data);
  /*
  {
    quiz_id: 123,
    sender_id: 456,
    sender_name: "John Doe",
    link_lottie: "https://example.com/animations/laugh.json",
    type: "reaction",
    timestamp: 1635789123
  }
  */
});
```

### Taunt Sent Confirmation (`taunt_sent`)

**Sender receives confirmation:**
```javascript
socket.on('taunt_sent', (data) => {
  console.log('Taunt sent successfully:', data);
  /*
  {
    quiz_id: 123,
    target_id: 789,
    link_lottie: "https://example.com/animations/laugh.json",
    type: "reaction"
  }
  */
});
```

## Usage Examples

### Example Lottie URLs for Taunts:
```javascript
const taunts = {
  laugh: "https://assets9.lottiefiles.com/packages/lf20_puciaact.json",
  cry: "https://assets5.lottiefiles.com/packages/lf20_OT1bp2.json", 
  angry: "https://assets2.lottiefiles.com/packages/lf20_qpsnmykx.json",
  shocked: "https://assets4.lottiefiles.com/packages/lf20_khzniaya.json",
  thinking: "https://assets8.lottiefiles.com/packages/lf20_khwrxejb.json",
  winner: "https://assets1.lottiefiles.com/packages/lf20_aEFaHc.json"
};
```

### Sending Different Types of Taunts:
```javascript
// Laughing reaction
socket.emit('opponent_taunt', {
  quiz_id: currentQuizId,
  user_id: currentUserId,
  link_lottie: "https://assets9.lottiefiles.com/packages/lf20_puciaact.json",
  type: "laugh"
});

// Thinking/confused reaction
socket.emit('opponent_taunt', {
  quiz_id: currentQuizId,  
  user_id: currentUserId,
  link_lottie: "https://assets8.lottiefiles.com/packages/lf20_khwrxejb.json",
  type: "thinking"
});
```

## Frontend Implementation

### React/Flutter Example:
```javascript
// Display received taunt
socket.on('taunt_received', (data) => {
  // Show Lottie animation
  showLottieAnimation(data.link_lottie);
  
  // Show sender info
  showToast(`${data.sender_name} sent a reaction!`);
  
  // Auto-hide after 3 seconds
  setTimeout(() => {
    hideLottieAnimation();
  }, 3000);
});

// Send taunt function
function sendTaunt(lottieUrl, type = 'reaction') {
  socket.emit('opponent_taunt', {
    quiz_id: currentQuizId,
    user_id: currentUserId, 
    link_lottie: lottieUrl,
    type: type
  });
}
```

## Error Handling

The system handles these error cases:
- `invalid_taunt_args`: Missing required parameters
- `session_not_found`: Quiz session doesn't exist
- `opponent_not_found`: No opponent in the session

## Features

- ✅ Real-time taunt delivery between opponents
- ✅ Support for Lottie animations 
- ✅ Sender name included in taunt
- ✅ Timestamp for taunt timing
- ✅ Confirmation when taunt is sent
- ✅ Error handling for edge cases
- ✅ Only works during active multiplayer sessions