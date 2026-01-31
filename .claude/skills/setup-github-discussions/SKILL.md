---
name: setup-github-discussions
description: Set up and manage GitHub Discussions for community engagement and Q&A
argument-hint: "[--init|--list|--enable CATEGORY] [--pin MESSAGE] [--welcome TEXT]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_discussions.py *)
---

# Setup GitHub Discussions

Enable and manage GitHub Discussions for community engagement, Q&A, announcements, and feature discussions.

## Usage

```
/setup-github-discussions --init                   # Enable Discussions feature
/setup-github-discussions --list                   # List discussion categories
/setup-github-discussions --enable "General"       # Enable specific category
/setup-github-discussions --pin "v1.0.0 Release"   # Pin discussion
/setup-github-discussions --welcome "Welcome!"     # Set welcome message
```

## Arguments

- `$0`: Action (--init to enable, --list to show, --enable/--pin/--welcome for management)
- `$1+`: Options (category name, message text, etc.)

## Discussion Categories

| Category | Purpose | Type |
|----------|---------|------|
| General | General discussion | Discussion |
| Q&A | Questions and answers | Question |
| Announcements | Project announcements | Announcement |
| Feature Requests | Feature ideas and requests | Discussion |
| Ideas | General ideas and brainstorming | Discussion |
| Show and Tell | Showcase projects/work | Discussion |

## Prerequisites

- GitHub repository settings access
- Admin permissions for repository
- GitHub CLI (`gh`) installed and authenticated

## Task

Set up GitHub Discussions for community engagement and support.

### Step 1: Enable Discussions

**Enable feature**:
```bash
python scripts/automation/github_discussions.py --init
```

**What it does**:
- Enables Discussions in repository settings
- Creates default categories
- Sets up moderation settings
- Configures permissions

### Step 2: List Categories

**View configured categories**:
```bash
python scripts/automation/github_discussions.py --list
```

**Shows**:
- Category names
- Descriptions
- Emoji icons
- Discussion count

### Step 3: Enable Specific Categories

**Enable Q&A category**:
```bash
python scripts/automation/github_discussions.py --enable "Q&A"
```

**Enable Feature Requests**:
```bash
python scripts/automation/github_discussions.py --enable "Feature Requests"
```

### Step 4: Configure Welcome Message

**Set welcome banner**:
```bash
python scripts/automation/github_discussions.py --welcome "Welcome to our community! Please read guidelines before posting."
```

### Step 5: Pin Important Discussions

**Pin announcement**:
```bash
python scripts/automation/github_discussions.py --pin "v1.0.0 Released"
```

**Pin community guidelines**:
```bash
python scripts/automation/github_discussions.py --pin "Community Guidelines"
```

## Default Categories

### General

**Purpose**: General project discussion

**Description**: General conversations about the project

**Icon**: 💬

**Usage**:
- Feature discussions
- Architecture decisions
- Project direction
- General talk

### Q&A

**Purpose**: Questions and answers

**Description**: Ask questions and get answers

**Icon**: 🤔

**Usage**:
- How-to questions
- Troubleshooting
- Documentation questions
- API usage questions

### Announcements

**Purpose**: Important announcements

**Description**: Latest news and announcements

**Icon**: 📢

**Usage**:
- Release announcements
- Breaking changes
- Important updates
- Roadmap updates

### Feature Requests

**Purpose**: Request new features

**Description**: Suggest new features

**Icon**: ✨

**Usage**:
- Feature ideas
- Enhancement requests
- Prioritization discussions
- Feature voting

### Ideas

**Purpose**: Brainstorming and ideas

**Description**: Share ideas and concepts

**Icon**: 💡

**Usage**:
- Design ideas
- Architecture suggestions
- Integration ideas
- Community proposals

### Show and Tell

**Purpose**: Showcase projects

**Description**: Share projects using our software

**Icon**: 🎉

**Usage**:
- Community projects
- Integration examples
- Custom extensions
- Success stories

## Examples

**Initialize Discussions**:
```bash
/setup-github-discussions --init
```

**List all categories**:
```bash
/setup-github-discussions --list
```

**Enable Q&A**:
```bash
/setup-github-discussions --enable "Q&A"
```

**Pin release announcement**:
```bash
/setup-github-discussions --pin "v1.0.0 Release"
```

**Set welcome message**:
```bash
/setup-github-discussions --welcome "Welcome to Revenge! Check docs before posting."
```

## Managing Discussions

### Creating Discussions Programmatically

**Via GitHub API**:
```bash
gh api graphql -f query='
  mutation {
    createDiscussion(input: {
      repositoryId: "R_kgDOXXXXXX"
      categoryId: "DIC_category"
      title: "New Discussion"
      body: "Discussion content"
    }) {
      discussion {
        id
        url
      }
    }
  }
'
```

### Moderating Discussions

**Pin important discussions**:
- Announcements
- Guidelines
- FAQ

**Hide spam/off-topic**:
- GitHub UI: Discussion → ⋯ → Hide discussion
- Moderates without deletion

**Delete if needed**:
- GitHub UI: Discussion → ⋯ → Delete discussion
- Removes from public view

### Answering Questions

**Mark answer in Q&A**:
1. Click "Mark as answer"
2. Response marked as solution
3. Appears at top of Q&A
4. Shows resolution to viewers

## Discussion Workflows

### For Users Asking Questions

1. **Search existing Q&A**
   - May already be answered
   - Faster solution

2. **Create new Q&A discussion**
   - Category: Q&A
   - Title: Clear question
   - Body: Detailed description
   - Include: Version, environment, error messages

3. **Engage with responders**
   - Comment with additional info
   - Clarify questions
   - Mark answer when solved

### For Maintainers Responding

1. **Monitor Q&A category**
   - Daily or weekly
   - Respond to new questions

2. **Provide helpful answers**
   - Direct to documentation
   - Provide code examples
   - Explain concepts clearly

3. **Mark answers**
   - When response solves issue
   - Mark as solution
   - Helps future visitors

### For Feature Requests

1. **Encourage discussion**
   - Get community input
   - Understand use cases
   - Gather requirements

2. **Provide feedback**
   - Feasibility assessment
   - Timeline if planned
   - Alternative solutions

3. **Convert to issues**
   - When ready to implement
   - Create GitHub Issue
   - Link to discussion

## Community Guidelines

**Create pinned guidelines discussion**:

```markdown
# Community Guidelines

1. **Be respectful** - Treat others with respect
2. **Stay on topic** - Keep discussions relevant
3. **Search first** - Check for existing discussions
4. **Provide context** - Include versions, error messages
5. **No spam** - No promotional or spam content
6. **Use right category** - Choose appropriate category
7. **No cross-posting** - Post once, not multiple times
8. **Report issues** - Use Issue Tracker for bugs

Violating these guidelines may result in:
- Comment deletion
- Discussion hiding
- User warnings
- Muting/banning
```

## Analytics and Insights

### Discussion Statistics

**View metrics**:
1. Settings → Discussions
2. Analytics tab
3. Shows:
   - Total discussions
   - Category breakdown
   - Trending topics
   - Activity by category

### Engagement Tracking

**Most active discussions**:
- By comment count
- By view count
- By last activity

**Category popularity**:
- Q&A most viewed
- Feature requests most discussed
- General second most active

## Integration with Other Features

### Link Issues to Discussions

**Reference in issue**:
```markdown
Discussed in #123 Discussion
```

**Link from discussion**:
```markdown
Related issue: #456
```

### Link PRs to Discussions

**Reference in PR**:
```markdown
Implements discussion #789
```

### Publish Release Notes

**In Announcements category**:
```markdown
# v1.0.0 Released

Changes:
- Feature X
- Fix for Y

See CHANGELOG.md for details
```

## Troubleshooting

**"Discussions not enabled"**:
1. Check repository settings
2. Discussions might be disabled
3. Run: `/setup-github-discussions --init`
4. Refresh page

**"Can't create discussion"**:
1. Verify you have permission
2. Check category is enabled
3. Try different category
4. Check for GitHub outages

**"Discussion not appearing**:
1. May be pending moderation
2. Check GitHub spam filters
3. Verify content follows guidelines
4. Try creating again

**"Comments not visible**:
1. May be hidden as spam
2. Check moderation queue
3. Refresh browser
4. Check if you're blocked

## Tips

1. **Start with few categories**:
   - General, Q&A, Announcements
   - Add more as community grows
   - Don't overwhelm with too many

2. **Promote discussions**:
   - Link in README
   - Mention in issues
   - Share on social media
   - Include in docs

3. **Active moderation**:
   - Regular monitoring
   - Respond quickly
   - Mark answers promptly
   - Handle spam quickly

4. **Encourage participation**:
   - Welcome new users
   - Ask for feedback
   - Celebrate contributions
   - Thank respondents

5. **Convert to issues**:
   - When feature request gets traction
   - When community consensus reached
   - Before starting implementation
   - Link back to discussion

## Best Practices

1. **Clear categories**:
   - Help users find right place
   - Reduce duplicates
   - Organize discussions

2. **Active responses**:
   - Answer within 24-48 hours
   - Welcome new members
   - Thank contributors

3. **Maintain guidelines**:
   - Enforce respectfully
   - Explain reasoning
   - Support community norms

4. **Regular moderation**:
   - Hide spam/off-topic
   - Manage uncivil behavior
   - Archive old discussions

5. **Leverage analytics**:
   - Identify common questions
   - Improve documentation
   - Track community health

## Related Skills

- `/manage-labels` - Organize with labels
- `/setup-github-projects` - Task tracking
- `/manage-ci-workflows` - CI/CD workflows
