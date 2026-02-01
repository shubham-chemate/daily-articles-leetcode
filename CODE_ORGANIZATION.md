# LeetCode Daily Articles Fetcher - Code Organization

## File Structure

The codebase is now organized into separate modules for better maintainability:

```
.
├── main.go              # Main entry point and orchestration
├── types.go             # Data structures (Article, Author, Tag, etc.)
├── leetcode_client.go   # LeetCode GraphQL API client
├── email_sender.go      # SendGrid email functionality
├── file_writer.go       # File output functionality
├── utils.go             # Utility functions (timestamp handling)
└── go.mod               # Go module definition
```

## Module Responsibilities

### 📄 `main.go`
**Purpose**: Application entry point and workflow orchestration

**Responsibilities**:
- Read configuration from environment variables
- Coordinate data fetching, email sending, and file writing
- Handle error management and exit codes
- Print progress and summary information

**Key Functions**:
- `main()`: Main execution flow

---

### 🏗️ `types.go`
**Purpose**: Define all data structures

**Structures**:
- `Article`: LeetCode discuss article
- `Author`: Article author information
- `Tag`: Article tags
- `Reaction`: Article reactions
- `ArticlesResponse`: GraphQL API response wrapper

---

### 🌐 `leetcode_client.go`
**Purpose**: Interface with LeetCode's GraphQL API

**Responsibilities**:
- Fetch articles from LeetCode
- Handle pagination
- Parse GraphQL responses
- Filter articles by creation time

**Key Functions**:
- `fetchArticlesAfterTime(cutoffTime)`: Fetch all articles after a specific time
- `fetchDiscussArticlesWithSkip(count, skip)`: Fetch a batch of articles with pagination

**Constants**:
- `leetcodeGraphQLURL`: LeetCode API endpoint
- `discussTopicsQuery`: GraphQL query for fetching articles

---

### 📧 `email_sender.go`
**Purpose**: Handle email delivery via SendGrid

**Responsibilities**:
- Generate HTML email templates
- Send emails via SendGrid API
- Handle HTML escaping and text truncation

**Key Functions**:
- `sendEmailViaSendGrid(...)`: Send email using SendGrid API
- `generateHTMLEmail(articles, ist)`: Create beautiful HTML email from articles
- `escapeHTML(s)`: Escape special HTML characters
- `truncateText(s, maxLen)`: Truncate text with ellipsis

**Constants**:
- `sendGridAPIURL`: SendGrid API endpoint

**Data Structures**:
- `SendGridEmail`: Email payload structure
- `EmailAddress`: Email address with optional name
- `Content`: Email content (HTML/text)

---

### 📝 `file_writer.go`
**Purpose**: Write articles to text files

**Responsibilities**:
- Format articles for text output
- Create organized text files
- Handle file creation and writing

**Key Functions**:
- `writeArticlesToFile(articles, filename)`: Format and write articles to file

**Output Format**:
- Header with article count and timestamp
- Detailed article information (UUID, title, author, tags, reactions)
- Clean, readable text format

---

### 🛠️ `utils.go`
**Purpose**: Shared utility functions

**Responsibilities**:
- Timestamp persistence (read/write)
- Timestamp formatting (ISO to IST)

**Key Functions**:
- `readLastProcessedTimestamp()`: Read last processed timestamp from file
- `writeLastProcessedTimestamp(t)`: Write timestamp to file
- `formatStringTimestamp(ts)`: Convert ISO timestamp to IST readable format

**Constants**:
- `lastTimestampFile`: Filename for timestamp persistence

---

## Data Flow

```
main()
  │
  ├─► readLastProcessedTimestamp() [utils.go]
  │
  ├─► fetchArticlesAfterTime() [leetcode_client.go]
  │     └─► fetchDiscussArticlesWithSkip() [leetcode_client.go]
  │
  ├─► generateHTMLEmail() [email_sender.go]
  │     └─► sendEmailViaSendGrid() [email_sender.go]
  │
  ├─► writeArticlesToFile() [file_writer.go]
  │
  └─► writeLastProcessedTimestamp() [utils.go]
```

## Benefits of This Organization

### ✅ **Separation of Concerns**
Each file has a single, clear responsibility

### ✅ **Easy to Test**
Functions can be tested independently

### ✅ **Better Maintainability**
Changes to email logic don't affect file writing or API calls

### ✅ **Reusability**
Modules can be used in other projects (e.g., email_sender, leetcode_client)

### ✅ **Clear Dependencies**
Easy to see what depends on what

### ✅ **Easier Debugging**
Issues can be isolated to specific modules

## Building and Running

```bash
# Build
go build

# Run with file output only
ENABLE_FILE_OUTPUT="true" go run .

# Run with email (requires configuration)
SENDGRID_API_KEY="key" FROM_EMAIL="from@example.com" TO_EMAILS="to@example.com" go run .

# Run with both
SENDGRID_API_KEY="key" FROM_EMAIL="from@example.com" TO_EMAILS="to@example.com" ENABLE_FILE_OUTPUT="true" go run .
```

## Adding New Features

### Adding a new output format (e.g., JSON):
1. Create `json_writer.go`
2. Implement `writeArticlesToJSON(articles, filename)`
3. Add call in `main.go` after fetching articles

### Adding a new data source:
1. Create `another_api_client.go`
2. Implement fetch functions following the same pattern
3. Update `main.go` to fetch from multiple sources

### Adding new email provider:
1. Create `new_email_sender.go`
2. Implement similar interface to `sendEmailViaSendGrid()`
3. Update `main.go` to choose provider based on environment variable

## Code Style

- **Exported functions**: Start with capital letter (e.g., `Article`)
- **Private functions**: Start with lowercase (e.g., `fetchArticlesAfterTime`)
- **Error handling**: Return errors, don't panic
- **Comments**: Document all exported types and functions
- **Imports**: Standard library first, then third-party (we use stdlib only)
