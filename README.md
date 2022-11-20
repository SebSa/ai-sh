# AI Shell

This is a simple golang project which allows one to use GPT-3 to generate shell commands via natural language prompts and then execute them while in a shell.

Obviously this is quite unpredictable, but it can be save time at the shell and it's faster than googling for the command you want to run in most cases.

## Usage

Pre-requisites:

-   [GPT-3 API key](https://beta.openai.com/)

To obtain an API key, you need to create an account on [OpenAI](https://beta.openai.com/).

Then you can go to the [OpenAI Api Keys](https://beta.openai.com/account/api-keys) page and create a new API key.

Once you have the API key, you can run the following command and the program will save the key for later use:

```bash
ai auth <api-key>
```

Then, you can use the `ai ask` command to generate shell commands.

For example:

```bash
ai ask "List all files in current directory and subdirectories, print total average of their size in Mb"
```

Returned output:
    
```bash
find . -type f -exec du -sh {} + | awk '{ total += $1; count++ } END { print total/count/1024 }'

Run this command? (y/n): 
y
0.00406703
```

## Installation

Download latest release from releases page, rename it "ai" or what ever else you want and save it in a directory in your PATH.

On linux you can do this with:

```bash
mv ai-linux-amd64 /usr/local/bin/ai
chmod +x /usr/local/bin/ai
```