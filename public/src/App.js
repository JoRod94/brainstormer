import React, { Component } from 'react';
import './App.css';
import styled from "styled-components";
import Message from "./components/Message";
import IdeaPanel from './components/IdeaPanel';

const ChatBox = styled("div")`
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: auto;
  height: 85vh;
`;

const ChatBoxWrapper = styled("div")`
  border-right: 2px solid black;
  width: 100%;
`;

const InputFooter = styled("div")`
  padding: 0 20px;
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  border-top: 2px solid black;
  height: 10vh;
  color: white;
  text-align: center;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  background-color: white;
`;

const UsernameInput = styled("input")`
  border: none;
  padding: 5px;
  border: 2px solid black;
  box-shadow: 5px 5px 0px 0px rgba(0,0,0,1);
`;

const ChatInput = styled("input")`
  border: none;
  padding: 5px;
  border: 2px solid black;
  width: 60vw;
  box-shadow: 5px 5px 0px 0px rgba(0,0,0,1);
`;

const SubmitButton = styled("input")`
  border: none;
  padding: 5px;
  border: 2px solid black;
  background: white;
  color: black;
  font-weight: 500;
  margin-left: 10px;
  box-shadow: 5px 5px 0px 0px rgba(0,0,0,1);
`;

const ContentBox = styled("div")`
  display: flex;
  flex-direction: row;
  width: 100%;
`;

const Wrapper = styled("div")`
  display: flex;
  flex-direction: column;
`;

const Puller = styled("div")`
  height: 1px;
  margin-top: 10vh;
`;

const Header = styled("div")`
  border-bottom: 2px solid black;
  height: 5vh;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const Greeting = styled("div")`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
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
      ideas: [],
      currentIdea: -1,
      chatSocket,
    }
  }

  onChatSocketConnect(event){
    console.log("Connected...");
  }

  onChatSocketMessage(event){
    const parsedData = JSON.parse(event.data);
    if(parsedData.messages || parsedData.ideas){
      this.setState(prevState => ({
        messages: parsedData.messages || prevState.messages,
        ideas: parsedData.ideas,
        currentIdea: parsedData.currentIdea,
      }));
    }
    else{
      this.setState(prevState => ({messages: [...prevState.messages, parsedData]}));
    }
  }

  onMessageSubmit(event){
    event.preventDefault();
    const newMessage = {
      username: this.state.username,
      text: event.target[0].value,
    }
    if(event.target[0].value.length != ""){
      this.state.chatSocket.send(JSON.stringify(newMessage));
    }
    event.target[0].value = "";
    this.refs["bottom"].scrollIntoView({behavior: "smooth"});
  }

  onUsernameSubmit(event){
    event.preventDefault();
    this.setState({username: event.target[0].value});
  }

  render() {
    const { messages, username, ideas, currentIdea } = this.state;

    if(!username){
      return (
        <Greeting className="App">
          <Header>
            <h1> B R A I N S T O R M E R </h1>
          </Header>
          <form onSubmit={this.onUsernameSubmit}>
            <h1>Pick a username</h1><br/>
            <UsernameInput type="text" name="username"/>
            <SubmitButton type="submit" value="Submit"/>
          </form>
        </Greeting>  
      )
    }

    return (
      <Wrapper className="App">
        <Header>
          <h2> B R A I N S T O R M E R </h2>
        </Header>
        <ContentBox>
          <ChatBoxWrapper>
            <ChatBox >
              {messages.map(function(m){
                return <Message {...m}/>
              })}
              <Puller ref="bottom"/>
            </ChatBox>
          </ChatBoxWrapper>
          {ideas.length > 0 && <IdeaPanel ideas={ideas} currentIdea={currentIdea}/>}
        </ContentBox>
        <InputFooter>
          <form onSubmit={this.onMessageSubmit}>
            <ChatInput type="text" name="message" placeholder="Commands: /newidea  /idea  /why  /whynot  /vote"/>
            <SubmitButton type="submit" value="SUBMIT"/>
          </form>
        </InputFooter>
      </Wrapper>
    );
  }
}

export default App;
