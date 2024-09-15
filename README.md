# Go Demo App

## Overview

Go Demo App is a sample application written in Go. This application demonstrates basic command-line argument parsing and provides usage information. It is structured to showcase functionality such as managing application configurations and handling various roles and ports.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Command-Line Arguments](#command-line-arguments)
- [Logging](#logging)
- [Changelog](#changelog)
- [Contributing](#contributing)

## Prerequisites

- Go (version 1.14 or higher)
- Make sure you have a working Go environment set up on your machine.

## Installation

Clone the repository:

```bash
git clone https://github.com/den-vasyliev/go-demo-app.git
cd go-demo-app
```

Install the required dependencies:

```bash
go mod tidy
```

## Usage

To run the application, you need to provide the necessary command-line arguments. You can view the usage instructions by executing:

```bash
go run src/main.go -h
```

This will print the following usage information:

```
Usage: app [-name name] [-role role] [-port port]
```

## Command-Line Arguments

The application accepts the following command-line arguments:

- `-name name`: Specify the name of the application or user.
- `-role role`: Define the role assigned to the application or user.
- `-port port`: Indicate the port on which the application should run.

To represent an application components diagram that includes components like NATS, Redis, MySQL, and various application roles (ASCII, Data, IMG, API), you can structure it as follows. Below is a textual description of the diagram, and I'll suggest how you might visualize it.

### Application Components Diagram Description

#### Components:
1. **API Layer**
   - Service that handles incoming HTTP requests.
   - Connects to various application roles and services.

2. **Application Roles**
   - **ASCII Role**
     - Responsible for handling ASCII-related functionality. Communicates with the API and possibly interacts with NATS and Redis.
   - **Data Role**
     - Manages data processing and storage, likely interacting with MySQL for persistent storage.
   - **IMG Role**
     - Handles image processing. May communicate with both the API and Redis for caching images and may also involve NATS for notifications/events.

3. **Message Broker**
   - **NATS**
     - Used for communication between different application roles (ASCII, Data, IMG). Facilitates asynchronous messaging and event handling.

4. **Caching System**
   - **Redis**
     - Provides caching for frequently accessed data. Used by roles for quick data retrieval to improve performance.

5. **Database**
   - **MySQL**
     - Relational database for persistent data storage. Used primarily by the Data role for CRUD operations.


### Example Diagram Representation (Textual)

```
+------------------------+
|        API Layer       |<----------------------------+
|    (Handles requests)  |                             |
+-----------+------------+                             |
            |                                          |
            |                                          |
            v                                          |
+-----------+------------+                             |
|      ASCII Role       |                             |
|  (Handles ASCII data) |                             |
+-----------+------------+                             |
            |                                          |
            |                                          |
            |                                          |
+-----------v------------+  +-----------------+  +-----------------+
|      NATS              |  |   Data Role     |  |    IMG Role     |
| (Message Broker)       |  |  (Data storage) |  | (Image handling)|
+-----------+------------+  +-----------------+  +-----------------+
            |                                          |
            |                                          |
            |                                          |
+-----------v------------+                             |
|       Redis            |                             |
| (Caching System)       |                             |
+-----------+------------+                             |
            |                                          |
            |                                          |
            v                                          |
+-----------v------------+                             |
|       MySQL            |<---------------------------+
|  (Persistent Storage)  |
+------------------------+
```

Based on the provided code snippets, here's a textual diagram illustrating the process of handling a user request for image conversion to ASCII through the NATS messaging system. This diagram focuses on the interaction between the `subscribeAndPublish` function and the `ImgHandler` function.

### Textual Diagram of NATS Message Interconnection for Image Conversion

```
+---------------------------+
|       User Request        |
|      (API Layer)         |
|     (HTTP Post Request)  |
+------------+--------------+
             |
             v
+------------+--------------+
|    subscribeAndPublish    |<------------------------------------------------+
| (src/apiHandler.go)      |                                                 |
+------------+--------------+                                                 |
             |                                                                |
             | Publish Request to NATS                                        |
             |                                                                |
             v                                                                |
+------------+--------------+                                                 |
|          NATS            |                                                 |
| (Message Broker)         |                                                 |
+------------+--------------+                                                 |
             |                                                                |
             | Receive Message (Reply)                                        |
             |                                                                |
             v                                                                |
+------------+--------------+                                                 |
|        ImgHandler        |                                                 |
| (src/imgHandler.go)     |                                                 |
+------------+--------------+                                                 |
             |                                                                |
             | Convert Image to ASCII                                        |
             |                                                                |
             v                                                                |
+------------+--------------+                                                 |
|   Convert Library         |                                                 |
| (image2ascii/convert)    |                                                 |
+------------+--------------+                                                 |
             |                                                                |
             | Return ASCII String                                           |
             |                                                                |
             v                                                                |
+------------+--------------+                                                 |
|       NATS Publish        |                                                 |
|   Send Reply back to API  |                                                 |
+---------------------------+                                                 |
```

### Explanation:

1. **User Request:**
   - The process starts with the user sending an HTTP POST request to the API layer.
2. **subscribeAndPublish:**
   - The `subscribeAndPublish` function subscribes to a unique reply-to channel and publishes the image processing request (along with necessary parameters) to the NATS subject.
3. **NATS:**
   - The NATS message broker facilitates communication between the API layer and the ImgHandler.
4. **ImgHandler:**
   - The `ImgHandler` handles the message received, retrieves any relevant options, and performs the image conversion using the `image2ascii` library.
5. **Convert Library:**
   - The image is converted to an ASCII string.
6. **Return ASCII String:**
   - The final ASCII string is sent back through NATS, which is then relayed to the API layer.
