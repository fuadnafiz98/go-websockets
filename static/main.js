"use strict";

const messages = document.getElementById("messages");
const inputMessage = document.getElementById("input");
const submitButton = document.getElementById("submit");

const scrollToView = () => {
  messages.scrollTo({
    left: 0,
    top: messages.scrollHeight,
    behavior: "smooth",
  });
};

const sendWelcomeMessage = (data) => {
  const messageBox = document.createElement("div");
  messageBox.className =
    "flex w-full items-start justify-between bg-green-700 px-6 py-2";
  const message = document.createElement("div");
  message.className = "text-green-100";
  message.innerText = data.message;

  const timeDiv = document.createElement("div");
  timeDiv.className = "text-md text-green-200";
  timeDiv.innerText = data.created.split("T")[1].slice(0, 8);

  messageBox.appendChild(message);
  messageBox.appendChild(timeDiv);

  messages.append(messageBox);
};

const sendLogoutMessage = (data) => {
  const messageBox = document.createElement("div");
  messageBox.className =
    "flex w-full items-start justify-between bg-red-700 px-6 py-2";
  const message = document.createElement("div");
  message.className = "text-red-100";
  message.innerText = data.message;

  const timeDiv = document.createElement("div");
  timeDiv.className = "text-md text-red-200";
  timeDiv.innerText = data.created.split("T")[1].slice(0, 8);

  messageBox.appendChild(message);
  messageBox.appendChild(timeDiv);

  messages.append(messageBox);
};

const addMessage = (data) => {
  console.log(data);
  if (data.messageType === "WELCOME_MESSAGE") {
    sendWelcomeMessage(data);
    scrollToView();
    return;
  }
  if (data.messageType == "LEAVE_MESSAGE") {
    sendLogoutMessage(data);
    scrollToView();
    return;
  }
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
  comment.innerText = data.message;

  message.appendChild(username);
  message.appendChild(comment);

  const timeDiv = document.createElement("div");
  timeDiv.className = "text-md text-zinc-400";
  timeDiv.innerText = data.created.split("T")[1].slice(0, 8);

  messageBox.appendChild(message);
  messageBox.appendChild(timeDiv);

  messages.append(messageBox);
  scrollToView();
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
    addMessage(JSON.parse(event.data));
  });
};

const sendMessage = async () => {
  const response = await fetch("/publish", {
    method: "POST",
    body: inputMessage.value,
  });
  inputMessage.value = "";
  inputMessage.focus();
};

submitButton.addEventListener("click", async (event) => {
  event.preventDefault();
  await sendMessage();
});

const main = () => {
  connect();
};

main();
