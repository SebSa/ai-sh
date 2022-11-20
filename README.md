# AI Shell

This is a simple golang binary which allow one to use GPT-3 to generate shell commands via natural language prompts.

## Usage

Authenticate with OpenAI inputting your API key:

```bash
ai auth <api-key>
```

To obtain an API key, you need to create an account on [OpenAI](https://beta.openai.com/).

Then you can go to the [OpenAI Api Keys](https://beta.openai.com/account/api-keys) page and create a new API key.

Then, you can use the `ai ask` command to generate shell commands:

```bash
ai ask "ls -l"
```

## Installation

Download latest release from releases page, rename it "ai" or what ever else you want and save it in a directory in your PATH.

On linux you can do this with:

```bash

mv ai-linux-amd64 /usr/local/bin/ai

```

