import React, { Component } from 'react';
import './App.css';
import styled from "styled-components";
import Message from "./components/Message";

const ChatBox = styled("div")`
  width: 100%;
  overflow: scroll;
`;

const InputFooter = styled("div")`
  padding: 0 20px;
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  border-top: 1px solid black;
  height: 60px;
  color: white;
  text-align: center;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  background-color: white;
`;

const InputForm = styled("form")`
  height: 100%;
`;

const UsernameInput = styled("input")`
  border: none;
  padding: 5px;
  border: 1px solid black;
  border-radius: 3px;
`;

const ChatInput = styled("input")`
  border: none;
  padding: 5px;
  border: 1px solid black;
  border-radius: 3px;
  width: 700px;
`;

const SubmitButton = styled("input")`
  border: none;
  padding: 5px;
  border: 1px solid black;
  background: black;
  border-radius: 3px;
  color: white;
`;

const Wrapper = styled("div")`
  display: flex;
  flex-direction: column;
`;

const Puller = styled("div")`
  margin-top: 60px;
  height: 60px;
`;

class App extends Component {
  constructor(props){
    super(props);
    const chatSocket = new WebSocket(`ws://${window.location.host}/ws`);
    this.onChatSocketMessage = this.onChatSocketMessage.bind(this);
    this.onMessageSubmit = this.onMessageSubmit.bind(this);
    this.onUsernameSubmit = this.onUsernameSubmit.bind(this);
    chatSocket.addEventListener('open', this.onChatSocketConnect);
    chatSocket.addEventListener('message', this.onChatSocketMessage);

    this.state = {
      messages: [],
      chatSocket,
    }
  }

  onChatSocketConnect(event){
    console.log("Connected...");
  }

  onChatSocketMessage(event){
    const parsedData = JSON.parse(event.data);
    if(parsedData.messages){
      this.setState({messages: parsedData.messages});  
    }
    else{
      this.setState({messages: [...this.state.messages, parsedData]});
    }
  }

  onMessageSubmit(event){
    event.preventDefault();
    const newMessage = {
      username: this.state.username,
      timestamp: "today",
      text: event.target[0].value,
    }
    this.state.chatSocket.send(JSON.stringify(newMessage));
    event.target[0].value = "";
    this.refs["bottom"].scrollIntoView({behavior: "smooth"});
  }

  onUsernameSubmit(event){
    event.preventDefault();
    this.setState({username: event.target[0].value});
  }

  render() {
    const { messages, username } = this.state;

    if(!username){
      return (
        <div className="App">
          <InputForm onSubmit={this.onUsernameSubmit}>
            What is your username?<br/>
            <UsernameInput type="text" name="username"/>
            <SubmitButton type="submit" value="Submit"/>
          </InputForm>
        </div>  
      )
    }

    return (
      <Wrapper>
        <ChatBox className="App">
          {messages.map(function(m){
            return <Message {...m}/>
          })}
        </ChatBox>
        <Puller ref="bottom"/>
        <InputFooter>
          <form onSubmit={this.onMessageSubmit}>
            <ChatInput type="text" name="message"/>
            <SubmitButton type="submit" value="Submit"/>
          </form>
        </InputFooter>
      </Wrapper>
    );
  }
}

export default App;
