var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}


document.addEventListener('DOMContentLoaded', () => {
    const flashMessage = document.querySelector('.flash');
    if (flashMessage) {
        setTimeout(() => {
            flashMessage.classList.add('flash-hide');
        }, 3500);
    }
    if (window.location.pathname.startsWith('/book/view/')) {
        const pathParts = window.location.pathname.split('/');
        const bookId = pathParts[pathParts.length - 1];

        const readingContentDiv = document.getElementById('reading-content');
        const authenticated = readingContentDiv.getAttribute('authenticated') === 'true';
        let currentPage = readingContentDiv.getAttribute('current-page');
        currentPage = parseInt(currentPage, 10);
        const bookDirectory = readingContentDiv.getAttribute('data-book-directory');
        if (readingContentDiv) {
            initializePaginationAndChat(bookId, bookDirectory, authenticated, currentPage);
        } else {
            console.error('Reading content division not found.');
        }
    }
});

function initializePaginationAndChat(bookId, bookDirectory, authenticated, currentPage) {

    if (authenticated) {
        console.log("authenticated")
        const chatLog = document.getElementById('chat-log');
        const userInput = document.getElementById('user-input');
        const sendButton = document.getElementById('send-button');
    
        if (!chatLog || !userInput || !sendButton) {
            console.error('Chat elements not found!');
            return; 
        }

        let websocket = new WebSocket(`wss://localhost:4000/ws/book/${bookId}`);
        websocket.onopen = (event) => {
            console.log("WebSocket connection opened:", event);
        };
    
        websocket.onmessage = (event) => {
            const message = event.data;
            displayMessage(message, 'llm');
        };
    
        sendButton.addEventListener('click', () => sendMessage());
        userInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                sendMessage();
            }
        });
    
        function sendMessage() {
            const messageText = userInput.value;
            if (messageText.trim() !== '') {
                const messageData = {
                    type: 'chat',
                    message: messageText,
                    page: currentPage
                };
                websocket.send(JSON.stringify(messageData));
                displayMessage(messageText, 'user');
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
        
        function updateProgress(pageNumber) {
            if (websocket && websocket.readyState === WebSocket.OPEN) {
                const progressData = { type: 'progress', page: pageNumber };
                websocket.send(JSON.stringify(progressData));
            } else {
                console.error('WebSocket is not initialized or not open.');
            }
        }
    }


    function loadPage(pageNumber) {
        let pageUrl = `${window.location.origin}/${bookDirectory}/page_${pageNumber}.txt`;
        fetch(pageUrl)
            .then(response => response.text())
            .then(text => {
                document.getElementById('reading-content').innerText = text;
                currentPage = pageNumber;
                document.getElementById('page-number').value = pageNumber;
                document.getElementById('total-pages').textContent = document.getElementById('reading-content').getAttribute('data-total-pages');
                if (authenticated){
                    updateProgress(currentPage);
                }
                
            })
            .catch(err => {
                console.error('Error loading page:', err);
                document.getElementById('reading-content').innerText = 'Failed to load page content.';
            });
    }


    document.getElementById('first-page').addEventListener('click', () => loadPage(1));
    document.getElementById('prev-page').addEventListener('click', () => changePage(-1));
    document.getElementById('next-page').addEventListener('click', () => changePage(1));
    document.getElementById('last-page').addEventListener('click', () => {
        const totalPages = parseInt(document.getElementById('total-pages').textContent);
        loadPage(totalPages);
    });
    document.getElementById('go-page').addEventListener('click', () => jumpToPage());
    document.getElementById('page-number').addEventListener('keypress', function(event) {
        if (event.key === "Enter") {
            jumpToPage();
            event.preventDefault();
        }
    });

    function changePage(direction) {
        const nextPage = currentPage + direction;
        const totalPages = parseInt(document.getElementById('total-pages').textContent);
        if (nextPage > 0 && nextPage <= totalPages) {
            loadPage(nextPage);
        }
    }

    function jumpToPage() {
        const pageNumber = parseInt(document.getElementById('page-number').value);
        const totalPages = parseInt(document.getElementById('total-pages').textContent);
        if (pageNumber >= 1 && pageNumber <= totalPages) {
            loadPage(pageNumber);
        }
    }

    loadPage(currentPage);
}