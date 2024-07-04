var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

function initializeBookChat(bookId) {
    const chatLog = document.getElementById('chat-log');
    const userInput = document.getElementById('user-input');
    const sendButton = document.getElementById('send-button');

    if (!chatLog || !userInput || !sendButton) {
        console.error('Chat elements not found!');
        return; 
    }

    let websocket = new WebSocket("ws://localhost:4000/ws/book/" + bookId);
    websocket.onopen = (event) => {
        console.log("WebSocket connection opened:", event);
    };

    websocket.onmessage = (event) => {
        const message = event.data;
        displayMessage(message, 'llm');
    };

    sendButton.addEventListener('click', () => {
        const message = userInput.value;
        if (message.trim() !== '') {
            websocket.send(message);
            displayMessage(message, 'user');
            userInput.value = '';
        }
    });

    function displayMessage(message, sender) {
        const messageElement = document.createElement('p');
        messageElement.classList.add(sender);
        messageElement.textContent = message;
        chatLog.appendChild(messageElement);
        chatLog.scrollTop = chatLog.scrollHeight;
    }
}

if (window.location.pathname.startsWith('/book/view/')) {
    const pathParts = window.location.pathname.split('/');
    const bookId = pathParts[pathParts.length - 1];

    initializeBookChat(bookId); 
}