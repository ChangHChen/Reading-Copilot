# Reading Copilot

Reading Copilot is a comprehensive web application designed to enhance the reading experience by providing dynamic interactions with books sourced from the Gutenberg Project. It incorporates advanced NLP services to support non-spoiler interactions and offers robust user management and book interaction features.

## Components

### Main Service

- **Security**: Utilizes `mkcert` for trusted local HTTPS development.
- **User Management**: Supports sign-up, login, password changes, and viewing personal information.
- **API Keys**: Incomplete functionality for setting up API keys for closed-end LLM integration. (Unfinished)
- **Authentication**: Stateless authentication tracking for secure user sessions.
- **Book Lists**: Fetches the most downloaded books from the Gutenberg Project via Gutendex API, available to both users and guests.
- **Search Functionality**: Allows searching for books through the Gutendex API.
- **Book Viewing**: Registered users can read books and interact with the Copilot chat service.
- **Reading Progress Tracking**: Tracks each user's reading history and progress on individual books.

### MySQL Database

- **User Data**: Stores user credentials and profile information.
- **Sessions**: Manages web service sessions.
- **Gutendex Caching**:
  - Caches top book lists for efficiency; updates daily.
  - Caches metadata for individual books to reduce repeated API calls.
- **Reading Progress**: Tracks and stores users' reading histories and progress.

### NLP Service

- **Technology**: Built with Python.
- **Models**:
  - Utilizes two LLM models: WhereIsAI/UAE-Large-V1 for embedding text chunks and one local LLM model serving as the backbone of the Copilot.
  - Closed-source models support. (Unfinished)
- **Book Processing**:
  - Upon the first view of a book, the main service requests the NLP service to initialize a Milvus collection, segment the book into chunks, embed these, and store them in the Milvus collection with metadata (e.g., page numbers).
- **Query Handling**:
  - Embeds user queries to retrieve the most relevant text segments.
  - Constructs prompts by combining text segments with user queries and specific instructions, then fetches responses from the Copilot LLM.

### Vector Database (Milvus)

- Stores embedded vectors of book text along with metadata to facilitate efficient retrieval during user interactions with the NLP service.

## Highlights

- **Non-Spoiler Interaction**: Limits the knowledge base to the pages preceding the user's current progress to prevent story spoilers.