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
    console.log(bookId)
    console.log("wss://localhost:4000/ws/book/" + bookId)
    let websocket = new WebSocket("wss://localhost:4000/ws/book/" + bookId);
    websocket.onopen = (event) => {
        console.log("WebSocket connection opened:", event);
    };

    websocket.onmessage = (event) => {
        const message = event.data;
        displayMessage(message, 'llm');
    };

    
    sendButton.addEventListener('click', () => {
        sendMessage();
    });

    userInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            sendMessage();
        }
    });

    function sendMessage() {
        const message = userInput.value;
        if (message.trim() !== '') {
            websocket.send(message);
            displayMessage(message, 'user');
            userInput.value = '';
        }
    }

    function displayMessage(message, sender) {
        const messageElement = document.createElement('p');
        messageElement.classList.add(sender);
        messageElement.textContent = message;
        chatLog.appendChild(messageElement);
        chatLog.scrollTop = chatLog.scrollHeight;
    }
}

function loadBookContent() {
    const readingContentDiv = document.getElementById('reading-content');
    if (!readingContentDiv) {
        console.error('Reading content div not found!');
        return;
    }
    let bookPath = readingContentDiv.getAttribute('data-book-url');
    if (!bookPath) {
        console.error('Book path not found!');
        return;
    }
    if (!bookPath.startsWith('http')) {
        bookPath = window.location.origin + "/" + bookPath;
    }
    fetch(bookPath)
        .then(response => response.text())
        .then(text => {
            readingContentDiv.innerText = text;
        })
        .catch(err => {
            console.error('Error loading the book content:', err);
            readingContentDiv.innerText = 'Failed to load the book content.';
        });
}

if (window.location.pathname.startsWith('/book/view/')) {
    const pathParts = window.location.pathname.split('/');
    const bookId = pathParts[pathParts.length - 1];

    initializeBookChat(bookId);
    loadBookContent();
}