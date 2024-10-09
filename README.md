## Forum

## Authors
1. Yousif Alsayyad
2. S.Hasan
3. Isa
4. S.Hussain

## Overview
This project is a web-based forum that allows users to communicate with each other by creating posts and comments, associating categories with posts, liking and disliking posts and comments, and filtering content. The project uses SQLite for data storage and Docker for containerization.

## Features
1. User Authentication:

- User Registration: 
- Users can register by providing an email, username, and password. Or using one of the two OAuth methods provided (github and google)
- The system checks if the email is already registered and returns an error if it is.
- Passwords are securely encrypted.

- Login Session: 
- Users can log in using their credentials, and a session is created using cookies.
- Sessions have an expiration date to limit how long a user stays logged in.
- Sessions are stored in the database.
- A UUID is used for session management.

2. Communication:

- Posts and comments:
- Registered users can create posts and comments.
- Posts can be associated with one or more categories chosen by the user.
- All posts and comments are visible to both registered and non-registered users.
- Non-registered users can only view posts and comments.

3. Likes and Dislikes:

- Only registered users can like or dislike posts and comments.
- The total number of likes and dislikes is visible to all users.

4. Filtering:

- Users can filter posts by categories, created posts, and liked posts.
- Liked and created posts are visible for each user in his user-info page.
- Filtering by created and liked posts is only available to registered users and applies to the logged-in user.

# Docker:

This project is containerized using Docker. Docker simplifies the deployment process by packaging the application and its dependencies into a single container.

## Setup Instructions

1. Clone the repository:
```
git clone [repository-url]
cd forum-project
```
2. Build the docker image and run the container using the script in the bash file:
```
bash docker.sh 
OR
./docker.sh
```
3. Access the page on `http://localhost:8080`




