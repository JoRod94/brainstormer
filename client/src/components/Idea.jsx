import React, { Component } from 'react';
import styled from "styled-components";

const Idea = styled("div")`
    border: 1px solid black;
    border-radius: 3px;
    margin-bottom: 10vh;
    margin-top: 10px;
`;

const Header = styled("h3")`
`;

const Board = styled("div")`
    display: flex;
    justify-content: space-around;
`;

const PropertyColumn = styled("div")`
    display: flex;
    flex-direction: column;
    padding: 0 10px;
    padding-bottom: 20px;
    border: 1px solid black;
    border-radius: 3px;
    width: 180px;
`;

const Property = styled("span")`
    padding: 10px;
    border: 1px solid black;
    border-radius: 3px;
    margin: 3px 0;
`;

const Footer = styled("div")`
    margin-top: 20px;
    border-top: 1px solid black;
    padding: 20px;
`;

export default ({what, why, whynot, votes}) => (
  <Idea>
    <Header>
        Idea: {what}
    </Header>
    <Board>
        <PropertyColumn>
            <h3>Why</h3>
            {why.map(p => (<Property>{p}</Property>))}
        </PropertyColumn>
        <PropertyColumn>
            <h3>Why Not</h3>
            {whynot.map(p => (<Property>{p}</Property>))}
        </PropertyColumn>
    </Board>
    <Footer>
        Votes: {votes}
    </Footer>
  </Idea>
)