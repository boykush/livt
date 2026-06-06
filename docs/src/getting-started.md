# Getting Started

## Installation

```bash
go install github.com/boykush/livt@latest
```

## Quick Start

1. Create the required directories:

```bash
mkdir -p stories discoveries/usm discoveries/example-mappings ubiquitous
```

2. Create your first story in `stories/my-first-story.md`:

```markdown
---
name: My first story
---

As a user
I want to do something
So that I get value
```

3. Build and serve:

```bash
livt serve
```

4. Open <http://localhost:3000> in your browser. Every page has a sidebar to
   switch between Example Mappings, Story Maps, Stories, and Ubiquitous Language.
   Open **Stories** to find the story you just created:

![Stories list and sidebar navigation](images/stories-index.png)
