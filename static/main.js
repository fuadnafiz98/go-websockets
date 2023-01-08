"use strict";

const messages = document.getElementById("messages");
const inputMessage = document.getElementById("input");
const submitButton = document.getElementById("submit");

const addMessage = (data) => {
  const messageBox = document.createElement("div");
  messageBox.className =
    "flex w-full items-start justify-between bg-zinc-800 px-6 py-2";
  const message = document.createElement("div");
  message.className = "flex-1";
  const username = document.createElement("div");
  username.className = "text-lg font-bold text-zinc-300";
  username.innerText = "Username";

  const comment = document.createElement("div");
  comment.className = "text-zinc-100";
  comment.innerText = data;

  message.appendChild(username);
  message.appendChild(comment);

  const timeDiv = document.createElement("div");
  timeDiv.className = "text-md text-zinc-400";
  timeDiv.innerText = "12.20 PM";

  messageBox.appendChild(message);
  messageBox.appendChild(timeDiv);

  messages.append(messageBox);
};

const connect = () => {
  const socket = new WebSocket(`ws://${location.host}/subscribe`);
  socket.addEventListener("open", (event) => {
    console.log("New Connection stablished!");
  });
  socket.addEventListener("close", (event) => {
    console.log("Disconncted", event.code, " ", event.reason);
    if (event.code !== 1001) {
      setTimeout(connect, 1000);
    }
  });
  socket.addEventListener("message", (event) => {
    addMessage(event.data);
  });
};

submitButton.addEventListener("click", async (event) => {
  event.preventDefault();
  const response = await fetch("/publish", {
    method: "POST",
    body: inputMessage.value,
  });
  inputMessage.value = "";
  inputMessage.focus();
  messages.scrollTo({
    left: 0,
    top: messages.scrollHeight,
    behavior: "smooth",
  });
});

const main = () => {
  connect();
};

main();
