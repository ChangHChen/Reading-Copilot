var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}


document.addEventListener('DOMContentLoaded', () => {
    const readingContentDiv = document.getElementById('reading-content');
    const bookDirectory = readingContentDiv.getAttribute('data-book-directory');
    let currentPage = 1;

    function loadPage(pageNumber) {
        let pageUrl = `${window.location.origin}/${bookDirectory}/page_${pageNumber}.txt`;
        fetch(pageUrl)
            .then(response => response.text())
            .then(text => {
                readingContentDiv.innerText = text;
                currentPage = pageNumber;
                document.getElementById('page-number').value = pageNumber;
                document.getElementById('total-pages').textContent = readingContentDiv.getAttribute('data-total-pages');
            })
            .catch(err => {
                console.error('Error loading page:', err);
                readingContentDiv.innerText = 'Failed to load page content.';
            });
    }

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

    document.getElementById('first-page').addEventListener('click', () => loadPage(1));
    document.getElementById('prev-page').addEventListener('click', () => changePage(-1));
    document.getElementById('next-page').addEventListener('click', () => changePage(1));
    document.getElementById('last-page').addEventListener('click', () => {
        const totalPages = parseInt(document.getElementById('total-pages').textContent);
        loadPage(totalPages);
    });
    document.getElementById('go-page').addEventListener('click', jumpToPage);
    document.getElementById('page-number').addEventListener('keypress', function(event) {
        if (event.key === "Enter") {
            jumpToPage();
            event.preventDefault();
        }
    });

    loadPage(currentPage);


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
                const messageWithPage = message + " (page: " + currentPage + ")"; 
                websocket.send(messageWithPage);
                displayMessage(messageWithPage, 'user');
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


    if (window.location.pathname.startsWith('/book/view/')) {
        const pathParts = window.location.pathname.split('/');
        const bookId = pathParts[pathParts.length - 1];
        initializeBookChat(bookId);
    }
});