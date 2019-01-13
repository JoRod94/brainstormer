import React, { Component } from 'react';
import styled from "styled-components";

const Message = styled("div")`
  text-align: left;
  padding: 2px 10px;
  border-bottom: 1px solid black;
`;

export default ({username, text, timestamp}) => (
  <Message>
    <p>{username}: {text} </p>
  </Message>
)