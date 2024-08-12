docker build -t forum-app . 
docker run --name forum-container -p 8080:8080 forum-app