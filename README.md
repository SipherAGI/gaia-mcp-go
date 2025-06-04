# Gaia AI Image Generator for MCP-Compatible AI Systems

**Transform your AI conversations by adding powerful AI image generation capabilities!**

Create stunning AI-generated images directly within your favorite AI applications using Gaia's advanced image generation technology. Whether you're a creative professional, content creator, or simply someone who loves visual storytelling, this MCP server brings professional-quality AI image generation to any AI system that supports the Model Context Protocol (MCP).

## Table of Contents

- [What is Gaia MCP Server?](#what-is-gaia-mcp-server)
- [Supported AI Systems](#supported-ai-systems)
- [What You Can Do](#what-you-can-do)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Step 1: Get Your Gaia API Key](#step-1-get-your-gaia-api-key)
  - [Step 2: Download the Server](#step-2-download-the-server)
  - [Step 3: Configure Your AI System](#step-3-configure-your-ai-system)
    - [Claude Desktop Setup](#claude-desktop-setup)
    - [Other MCP-Compatible Systems](#other-mcp-compatible-systems)
  - [Step 4: Test Your Setup](#step-4-test-your-setup)
- [Available Image Tools](#available-image-tools)
- [Example Usage](#example-usage)
- [Troubleshooting](#troubleshooting)
- [Requirements & Credits](#requirements--credits)
- [Support](#support)
- [License](#license)

## What is Gaia MCP Server?

The Gaia MCP Server is a bridge that connects MCP-compatible AI systems with Gaia's powerful AI image generation platform. Think of it as a translator that allows your AI assistant to understand your image requests and communicate with Gaia's servers to create the images you want.

**No technical knowledge required!** Simply follow our step-by-step setup guide, and you'll be generating images in minutes.

## Supported AI Systems

This MCP server works with any AI system that supports the **Model Context Protocol (MCP)**. Currently supported systems include:

### ✅ Fully Tested & Documented

- **Claude Desktop** - Complete setup guide provided below
- **Claude Web** (with MCP support)

### ✅ MCP-Compatible (Should Work)

- **Cline** (VS Code Extension)
- **Zed Editor** (with MCP support)
- **Continue** (VS Code/JetBrains Extension)
- **Other MCP-compatible AI tools and editors**

### 🔧 Custom Integration

- **Any application** that can communicate via MCP stdio or SSE protocols
- **Custom AI systems** with MCP client implementation

_Don't see your preferred AI system listed? If it supports MCP, this server should work with it! Check our [troubleshooting section](#troubleshooting) or [open an issue](https://github.com/SipherAGI/gaia-mcp-go/issues) for help._

## What You Can Do

With Gaia MCP Server integrated into your AI system, you can:

✨ **Generate Original Images**: Create unique artwork, illustrations, and photos from text descriptions  
🎨 **Remix Existing Images**: Transform and reimagine existing images with new styles or elements  
📸 **Enhance Faces**: Improve facial details and clarity in portraits  
🔍 **Upscale Images**: Increase image resolution and quality  
📤 **Upload Images**: Easily work with images from web URLs

All of this happens seamlessly within your AI conversations!

## Getting Started

### Prerequisites

Before you begin, make sure you have:

- **An MCP-compatible AI system** (see [Supported AI Systems](#supported-ai-systems) above)
- **A Gaia account** with available credits ([Sign up here](https://protogaia.com))
- **10 minutes** to complete the setup

### Step 1: Get Your Gaia API Key

Your API key is like a password that allows the server to access your Gaia account:

1. **Visit Gaia's Website**: Go to [protogaia.com](https://protogaia.com) and log in to your account
2. **Access Your Profile**: Click on your profile picture in the top-right corner
3. **Go to Settings**: Select "Account Settings" or "Settings" from the dropdown menu
4. **Find Security Section**: Look for a "Security" or "API Keys" section
5. **Create New API Key**: Click "Create New API Key" or similar button
6. **Copy Your Key**: **Important!** Copy the generated API key immediately and save it somewhere safe - you won't be able to see it again

### Step 2: Download the Server

1. **Visit the Releases Page**: Go to the [releases section](https://github.com/SipherAGI/gaia-mcp-go/releases) of this project
2. **Download for Your System**:
   - **Windows**: Download `gaia-mcp-go-windows.exe`
   - **Mac**: Download `gaia-mcp-go-darwin`
   - **Linux**: Download `gaia-mcp-go-linux`
3. **Save the File**: Save it to a location you'll remember (like your Desktop or Downloads folder)
4. **Note the Full Path**: You'll need to know exactly where you saved this file for the next step

### Step 3: Configure Your AI System

The configuration process varies depending on your AI system. We provide detailed instructions for Claude Desktop below, and general guidance for other systems.

#### Claude Desktop Setup

**Detailed step-by-step setup for Claude Desktop:**

1. **Open Claude Desktop** on your computer

2. **Access Settings**:

   - **Windows**: Go to File → Settings
   - **Mac**: Go to Claude → Settings (or use Cmd+,)

3. **Open Developer Settings**:

   - Click on the "Developer" tab
   - Click "Edit Config" button

4. **Add Configuration**: Replace everything in the configuration file with this code:

   ```json
   {
     "mcpServers": {
       "gaia-mcp-server": {
         "command": "/full/path/to/your/downloaded/gaia-mcp-go-file",
         "args": ["stdio", "--api-key=YOUR_ACTUAL_API_KEY_HERE"]
       }
     }
   }
   ```

5. **Customize Your Configuration**:

   - Replace `/full/path/to/your/downloaded/gaia-mcp-go-file` with the actual location where you saved the downloaded file
   - Replace `YOUR_ACTUAL_API_KEY_HERE` with the API key you copied from Gaia

   **Example**:

   ```json
   {
     "mcpServers": {
       "gaia-mcp-server": {
         "command": "/Users/johndoe/Desktop/gaia-mcp-go-darwin",
         "args": ["stdio", "--api-key=gaia_1234567890abcdef"]
       }
     }
   }
   ```

6. **Save and Close** the configuration window

#### Other MCP-Compatible Systems

**For other AI systems that support MCP:**

The general configuration pattern is:

```json
{
  "name": "gaia-mcp-server",
  "command": "/path/to/gaia-mcp-go-[your-system]",
  "args": ["stdio", "--api-key=YOUR_GAIA_API_KEY"]
}
```

**System-Specific Notes:**

- **Cline (VS Code)**: Add the server configuration to your Cline MCP settings
- **Zed Editor**: Configure in your Zed MCP settings file
- **Continue**: Add to your Continue extension MCP configuration
- **Custom Systems**: Use the stdio transport with the command and args shown above

_Need help with your specific AI system? [Open an issue](https://github.com/SipherAGI/gaia-mcp-go/issues) and we'll help you get set up!_

### Step 4: Test Your Setup

1. **Restart your AI system** completely (close and reopen the application)

2. **Start a New Conversation**

3. **Test Image Generation**: Try asking your AI assistant something like:

   - "Generate an image of a peaceful mountain lake at sunset"
   - "Create a cartoon drawing of a friendly robot"
   - "Make an image of a cozy coffee shop interior"

4. **Success!** If everything is working, you should see your AI assistant generate and display your requested image directly in the conversation.

## Available Image Tools

Here's what you can do with each tool:

### 🎨 Generate Image

**What it does**: Creates brand new images from your text descriptions
**Example**: "Generate an image of a futuristic city skyline with flying cars"

### 🔄 Remix Image

**What it does**: Takes an existing image and creates new variations or applies different styles
**Example**: "Remix this photo to look like a watercolor painting"

### 👤 Face Enhancer

**What it does**: Improves the quality and detail of faces in portraits
**Example**: "Enhance the facial features in this portrait"

### 🔍 Upscaler

**What it does**: Increases image resolution and overall quality
**Example**: "Upscale this image to make it higher resolution"

### 📤 Upload Image

**What it does**: Allows you to work with images from web URLs
**Example**: Upload an image from a website to use with other tools

## Example Usage

Here are some conversation examples to get you started:

**Simple Image Generation:**

- "Generate an image of a golden retriever playing in a park"
- "Create a logo design for a bakery with warm, friendly colors"

**Creative Projects:**

- "Generate an image of a steampunk-style airship flying over Victorian London"
- "Create a book cover design for a mystery novel set in the 1920s"

**Photo Enhancement:**

- "Upload this image: [URL] and enhance the faces in it"
- "Take this low-resolution image and upscale it to make it clearer"

## Troubleshooting

### Common Issues and Solutions

**Problem**: Claude says it can't generate images

- **Solution**: Make sure you restarted Claude Desktop after configuration
- **Check**: Verify your API key is correct and has sufficient credits

**Problem**: "Command not found" error

- **Solution**: Double-check the file path in your configuration
- **Check**: Make sure you downloaded the correct version for your operating system

**Problem**: API key errors

- **Solution**: Generate a new API key from your Gaia account
- **Check**: Ensure there are no extra spaces in your API key

**Problem**: Images aren't generating

- **Solution**: Check your Gaia account credit balance
- **Check**: Try with a simpler image request first

### Need More Help?

If you're still having trouble:

1. Check the [Issues section](https://github.com/SipherAGI/gaia-mcp-go/issues) for similar problems
2. Create a new issue with details about your problem
3. Contact Gaia support through their website

## Requirements & Credits

**System Requirements:**

- Claude Desktop application
- Internet connection
- Gaia account with available credits

**About Credits:**

- Each image generation uses Gaia credits from your account
- Different tools may use different amounts of credits
- Check your credit balance at [protogaia.com/settings/account](https://protogaia.com/settings/account?tab=Plans)
- Purchase additional credits as needed

## Support

**For Technical Issues**: [Open an issue on GitHub](https://github.com/SipherAGI/gaia-mcp-go/issues)  
**For Gaia Account Issues**: [Contact Gaia Support](https://protogaia.com/support)  
**For Claude Desktop Issues**: [Contact Anthropic Support](https://support.anthropic.com)

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

---

**Ready to start creating amazing AI images?** Follow the setup guide above and begin generating stunning visuals in your Claude Desktop conversations today!
