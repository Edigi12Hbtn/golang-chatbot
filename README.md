# golang-chatbot

##Golang based chatbot for a restaurant

To run the repository open a terminal, go to $GOPATH/src/github.com, and clone the repository.

Enter to the directory golang-chatbot/api-chatbot and run go run main/*

Now the api is listening at http://localhost:3000/sky-restaurant.

Open anoter terminal window and go to $GOPATH/src/github.com/golang-chatbot/web-chatbot and run the following commands:

    - npm install
    - npm run serve

Now the web page is running. Go to http://localhost:8080/#/ and open and write a greeting, an opinion about the food or order some food.
![Screenshot](img/api-chatbot.png)