import React, { Component } from 'react';
import styled from "styled-components";
import Idea from "./Idea";

const IdeaPanel = styled("div")`
    display: flex;
    flex-direction: column;
    padding: 20px;
    height: 85vh;
    overflow: auto;
`;

const IdeaPanelWrapper = styled("div")`
  width: 100%;
`;

const Index = styled("div")`
    display: flex;
    flex-direction: column;
`;

export default ({ideas, currentIdea}) => (
    <IdeaPanelWrapper>
        <IdeaPanel>
            <Index>
                <h3>All Ideas:</h3>
                {ideas.map((idea, index) => (<span>{index}: {idea.what}</span>))}
            </Index>
            {currentIdea >= 0 && <Idea {...ideas[currentIdea]}/>}
        </IdeaPanel>
    </IdeaPanelWrapper>
)